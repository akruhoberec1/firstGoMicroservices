package repository

import (
	"authService/internal/config"
	"authService/internal/model"
	"fmt"
	"time"
)

func GetRefreshToken(refreshToken string) (*model.Token, error) {
	token := &model.Token{}

	err := config.DB.QueryRow(
		"SELECT id, user_id, refresh_token, expires_at, revoked, created_at FROM tokens WHERE refresh_token = $1",
		refreshToken,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.Revoked,
		&token.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve refresh token: %w", err)
	}

	return token, nil
}

func StoreRefreshToken(token *model.Token) error {
	_, err := config.DB.Exec("INSERT INTO tokens (refresh_token, expires_at, user_id) VALUES ($1, $2, $3)", token.Token, token.ExpiresAt, token.UserID)
	return err
}

func RevokeRefreshToken(refreshToken string) error {
	result, err := config.DB.Exec("UPDATE tokens SET revoked = TRUE WHERE refresh_token = $1", refreshToken)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("refresh token not found or already revoked")
	}

	return nil
}

func SetAllTokensRevokedForUser(userID int) error {
	_, err := config.DB.Exec("UPDATE tokens SET revoked = TRUE WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all refresh tokens for user: %w", err)
	}
	return nil
}

func DeleteExpiredTokens() error {
	_, err := config.DB.Exec("DELETE FROM tokens WHERE expires_at < $1", time.Now())
	return err
}
