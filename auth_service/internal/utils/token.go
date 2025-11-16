package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSecureToken(size int) (string, error) {
	if size <= 0 {
		size = 32
	}

	bytes := make([]byte, size)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes)
	return token, nil
}
