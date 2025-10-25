package middleware

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rookgm/gophkeeper/internal/models"
	"github.com/rookgm/gophkeeper/internal/server/service"
	"net/http"
	"strings"
)

type contextKey string

const (
	authPayloadKey contextKey = "auth_payload"
)

// Auth gets the token from the cookie and passes it to the context
func Auth(ts service.TokenService) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "can not get Authorization header", http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid token format", http.StatusUnauthorized)
			}

			token := parts[1]

			payload, err := ts.VerifyToken(token)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), authPayloadKey, payload)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// getAuthPayload extracts authorization token payload from context
func getAuthPayload(ctx context.Context, key contextKey) (*models.TokenPayload, bool) {
	payload, ok := ctx.Value(key).(*models.TokenPayload)
	return payload, ok
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	authPayload, ok := getAuthPayload(ctx, authPayloadKey)
	if !ok || authPayload == nil {
		return uuid.Nil, errors.New("invalid auth payload")
	}
	return authPayload.UserID, nil
}
