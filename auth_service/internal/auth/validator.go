package auth

import (
    "errors"
    "regexp"
)

// emailRegex is compiled once at startup to prevent severe performance penalties during runtime
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail verifies the structural integrity of an email address string
func ValidateEmail(email string) error {
    if email == "" {
        return errors.New("email required")
    }
    if !emailRegex.MatchString(email) {
        return errors.New("invalid email format")
    }
    return nil
}

// ValidatePassword ensures the cryptographic material meets minimum entropy requirements
func ValidatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    return nil
}