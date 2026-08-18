package auth

import (
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
