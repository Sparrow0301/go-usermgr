package security

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given password string using bcrypt.
func HashPassword(password string, cost int) (string, error) {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword compares the hashed password with the provided plain password.
func VerifyPassword(hashed, password string) error {
	if hashed == "" || password == "" {
		return errors.New("password cannot be empty")
	}
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
