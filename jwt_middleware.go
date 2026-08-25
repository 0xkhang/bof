package main

import (
	"context"
	"net/http"

	"github.com/0xkhangle/bof/internal/auth"
	"github.com/0xkhangle/bof/internal/utils"
)

func (cfg *apiConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := utils.GetBearerToken(r)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, err, "Missing token")
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, err, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
