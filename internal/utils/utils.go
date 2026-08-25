package utils

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"github.com/google/uuid"
)

type contextKey string
const UserIDKey contextKey = "userID"

var (
	ErrInvalidHeaderFormat = errors.New("Invalid header format")
	ErrTokenNotProvided    = errors.New("Can't extract token")
	ErrInvalidID           = errors.New("Invalid ID format")
)

func IsEmailValid(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}
	return nil
}

func GetBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 {
		return "", ErrInvalidHeaderFormat
	}

	if parts[0] != "Bearer" {
		return "", ErrInvalidHeaderFormat
	}

	if parts[1] == "" {
		return "", ErrTokenNotProvided
	}

	return parts[1], nil
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID := ctx.Value(UserIDKey)

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrInvalidID
	}

	return id, nil
}
