package middleware

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/rookgm/gophkeeper/internal/models"
	"github.com/rookgm/gophkeeper/internal/server/http/middleware/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockHandler is a simple handler to be wrapped by the middleware
func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(("auth")))
}

func TestAuth(t *testing.T) {
	handler := http.HandlerFunc(mockHandler)

	tests := []struct {
		name        string
		tokenString string
		setup       func(t *testing.T) *mocks.MockTokenService
		statusCode  int
	}{
		// token exist
		{
			name:        "authorized",
			tokenString: "token",
			setup: func(t *testing.T) *mocks.MockTokenService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				ts := mocks.NewMockTokenService(ctrl)
				ts.EXPECT().VerifyToken(gomock.Any()).Return(&models.TokenPayload{
					ID:     uuid.New(),
					UserID: uuid.New(),
				}, nil).AnyTimes()
				return ts
			},
			statusCode: http.StatusOK,
		},
		// empty token
		{
			name:        "empty_token",
			tokenString: "",
			setup: func(t *testing.T) *mocks.MockTokenService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				ts := mocks.NewMockTokenService(ctrl)
				ts.EXPECT().VerifyToken(gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
				return ts
			},
			statusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			w := httptest.NewRecorder()
			req.Header.Set("Authorization", "Bearer "+tt.tokenString)
			ts := tt.setup(t)
			Auth(ts)(handler).ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.statusCode, res.StatusCode, "status code is not equal")
		})
	}
}
