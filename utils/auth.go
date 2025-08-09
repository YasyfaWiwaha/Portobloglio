package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password []byte) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "ERROR!", fmt.Errorf("error hashing password. %w", err)
	}
	return string(bytes), err
}

func CheckPassword(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}
