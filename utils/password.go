package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the becrypt hash of the input password
func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashPassword), nil
}

func CheckPasswordCorrect(password, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
