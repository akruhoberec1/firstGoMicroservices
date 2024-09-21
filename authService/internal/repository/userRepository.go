package repository

import (
	"authService/internal/config"
	"authService/internal/model"
	"database/sql"
	"errors"
)

func GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	err := config.DB.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
