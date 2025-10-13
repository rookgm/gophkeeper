package middleware

import (
	"context"
	"github.com/rookgm/gophkeeper/internal/models"
	"github.com/rookgm/gophkeeper/internal/service"
	"net/http"
)

type contextKey string

const (
	authPayloadKey contextKey = "auth_payload"
)

// Auth gets the token from the cookie and passes it to the context
func Auth(ts service.TokenService) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				http.Error(w, "can not get cookie", http.StatusUnauthorized)
				return
			}

			payload, err := ts.VerifyToken(cookie.Value)
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
