package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

var HashCost = bcrypt.DefaultCost

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), HashCost)

	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}
