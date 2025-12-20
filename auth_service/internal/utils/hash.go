package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword compares password with hash
func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}
