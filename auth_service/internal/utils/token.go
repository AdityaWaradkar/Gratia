package utils

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// GenerateSecureToken returns a cryptographically random, URL-safe string for refresh tokens
func GenerateSecureToken(length int) (string, error) {
    b := make([]byte, length)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return base64.RawURLEncoding.EncodeToString(b), nil
}

// GenerateAccessToken creates a signed JWT containing user identity and role claims
func GenerateAccessToken(userID string, role string, secret string, ttl time.Duration) (string, error) {
    claims := jwt.MapClaims{
        "sub":  userID,
        "role": role,
        "exp":  time.Now().Add(ttl).Unix(),
        "iat":  time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

// ValidateAccessToken parses and validates the signature and expiration of a JWT
func ValidateAccessToken(tokenString string, secret string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return []byte(secret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, fmt.Errorf("invalid token claims")
}