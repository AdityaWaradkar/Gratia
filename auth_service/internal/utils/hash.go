package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword generates a secure cryptographic hash with a fixed cost factor of 10
func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

// ComparePassword safely verifies a plain text input against stored cryptographic material
func ComparePassword(hash, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}