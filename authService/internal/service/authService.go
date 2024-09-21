package service

import (
	"authService/internal/helper"
	"authService/internal/model"
	"authService/internal/repository"
	"errors"
	"github.com/google/uuid"
	"time"
)

func Authenticate(username, password string) (map[string]string, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username")
	}

	if !verifyPassword(user.Password, password) {
		return nil, errors.New("invalid password")
	}

	accessToken, err := helper.GenerateAccessToken(user.ID)
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

/*func RevokeRefreshToken(refreshToken string) error {
	return repository.RevokeRefreshToken(refreshToken)
}*/

func RevokeAllRefreshTokens(userID int) error {
	return repository.SetAllTokensRevokedForUser(userID)
}
