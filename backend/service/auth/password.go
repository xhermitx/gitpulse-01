package auth

import "golang.org/x/crypto/bcrypt"

func ComparePassword(p1, p2 []byte) bool {
	if err := bcrypt.CompareHashAndPassword(p1, p2); err != nil {
		return false
	}
	return true
}

func HashedPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
