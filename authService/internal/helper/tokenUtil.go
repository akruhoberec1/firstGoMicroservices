// File: internal/helper/tokenUtil.go

package helper

import (
	"authService/internal/config"
	"authService/internal/repository"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.JwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired access token")
	}

	return token, nil
}

func ExtractUserIDFromToken(token *jwt.Token) (int, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return 0, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64) // JWT stores numbers as float64
	if !ok {
		return 0, errors.New("invalid user ID in token claims")
	}

	return int(userID), nil
}

func GenerateAccessToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JwtSecretKey)
}

func RefreshAccessToken(refreshToken string) (string, error) {
	token, err := repository.GetRefreshToken(refreshToken)
	if err != nil || token.Revoked || token.ExpiresAt.Before(time.Now()) {
		return "", errors.New("invalid or expired refresh token")
	}

	newAccessToken, err := GenerateAccessToken(token.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	return newAccessToken, nil
}
