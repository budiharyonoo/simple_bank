package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword is generate brcpt hased password from string
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed convert bcrypt password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword Check the password
func CheckPassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
