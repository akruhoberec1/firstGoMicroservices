package service

import (
	"authService/internal/model"
	"authService/internal/repository"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

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

func verifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func GetUserByUsername(username string) (*model.User, error) {
	return repository.GetUserByUsername(username)
}
