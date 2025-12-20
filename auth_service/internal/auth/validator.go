package auth

import (
	"errors"
	"regexp"
)

// ValidateEmail checks email format
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email required")
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

// ValidatePassword checks password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}
