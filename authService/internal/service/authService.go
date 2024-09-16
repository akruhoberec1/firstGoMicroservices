package service

import (
	"authService/internal/config"
	"authService/internal/model"
	"authService/internal/repository"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var jwtSecretKey = config.JwtSecretKey

func Authenticate(username, password string) (map[string]string, error) {
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username")
	}

	if !verifyPassword(user.Password, password) {
		return nil, errors.New("invalid password")
	}

	accessToken, err := generateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	tokens := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken.Token,
	}

	return tokens, nil
}

func verifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func generateAccessToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecretKey)
}

func generateRefreshToken(userID int) (*model.Token, error) {
	refreshToken := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	token := &model.Token{
		Token:     refreshToken,
		ExpiresAt: expiresAt,
		UserID:    userID,
	}

	err := repository.StoreRefreshToken(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func RegisterUser(username, password string) error {
	exists, err := repository.UserExists(username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username already taken")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	return repository.CreateUser(username, hashedPassword)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func RevokeRefreshToken(refreshToken string) error {
	return repository.RevokeRefreshToken(refreshToken)
}
