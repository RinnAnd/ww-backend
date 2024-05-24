package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), 14)
	if err != nil {
		return err
	}
	*password = string(hash)

	return nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
