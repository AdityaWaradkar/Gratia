package auth

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
)

// ValidateEmail checks if the given email has a valid format.
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email cannot be empty")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}

	return nil
}

// ValidatePassword checks if the password meets minimum security criteria.
// Criteria: at least 8 characters, 1 uppercase, 1 lowercase, 1 number, 1 special character.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString
		hasSpecial = regexp.MustCompile(`[!@#~$%^&*()+|_.,<>?/{}-]`).MatchString
	)

	if !hasUpper(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber(password) {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial(password) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// ValidateUsername checks if the username is non-empty and alphanumeric
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("username cannot be empty")
	}

	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validUsername.MatchString(username) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}

	return nil
}

// ValidateConfirmPassword ensures password and confirmPassword match
func ValidateConfirmPassword(password, confirmPassword string) error {
	if password != confirmPassword {
		return errors.New("password and confirm password do not match")
	}
	return nil
}
