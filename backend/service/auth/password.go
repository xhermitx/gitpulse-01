package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func ComparePassword(hashed []byte, plaintext []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashed, plaintext)

	return err == nil
}

func HashedPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
