package repository

import (
	"authService/internal/config"
	"authService/internal/model"
	"database/sql"
	"time"
)

func GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	err := config.DB.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func UserExists(username string) (bool, error) {
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE username=$1)", username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func CreateUser(username, hashedPassword string) error {
	_, err := config.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, hashedPassword)
	return err
}

func StoreRefreshToken(token *model.Token) error {
	_, err := config.DB.Exec("INSERT INTO tokens (refresh_token, expires_at, user_id) VALUES ($1, $2, $3)", token.Token, token.ExpiresAt, token.UserID)
	return err
}

// DeleteExpiredTokens for future reference when some scheduling happens
func DeleteExpiredTokens() error {
	_, err := config.DB.Exec("DELETE FROM tokens WHERE expires_at < $1", time.Now())
	return err
}

func RevokeRefreshToken(token string) error {
	_, err := config.DB.Exec("DELETE FROM tokens WHERE refresh_token = $1", token)
	return err
}
