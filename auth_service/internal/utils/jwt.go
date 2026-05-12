package utils

import (
	"time"

	"github.com/adityawaradkar/gratia/auth_service/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(
	userID string,
	role string,
	secret string,
	expirationMinutes int,
) (string, error) {

	claims := model.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(
					time.Duration(expirationMinutes) * time.Minute,
				),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(
	tokenString string,
	secret string,
) (*model.JWTClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&model.JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.JWTClaims)

	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}