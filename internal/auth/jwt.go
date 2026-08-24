package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, tokenExpires time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    "bof-access",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenExpires)),
		Subject:   userID.String(),
	}
	method := jwt.SigningMethodHS256
	t := jwt.NewWithClaims(method, claims)
	s, err := t.SignedString(tokenSecret)
	if err != nil {
		return "", err
	}
	return s, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	keyFunc := func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}

	t, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return uuid.Nil, err
	}

	if !t.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
