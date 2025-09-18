package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateSecureToken creates a random, URL-safe token string of given length (default 32 bytes)
func GenerateSecureToken(size int) (string, error) {
	if size <= 0 {
		size = 32
	}

	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Encode to URL-safe base64
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes), nil
}
