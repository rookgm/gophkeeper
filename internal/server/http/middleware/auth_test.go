package middleware

import (
	"github.com/golang/mock/gomock"
	"github.com/rookgm/gophkeeper/internal/server/http/handler/mocks"
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
	//token := auth.NewAuthToken([]byte("secretkey"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ts := mocks.NewMockTokenService(ctrl)

	handler := http.HandlerFunc(mockHandler)
	wrapperHandler := Auth(ts)(handler)

	tests := []struct {
		name        string
		tokenString string
		statusCode  int
	}{
		// token exist
		{
			name:        "authorized",
			tokenString: "token",
			statusCode:  http.StatusOK,
		},
		// empty token
		{
			name:        "empty_token",
			tokenString: "",
			statusCode:  http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			w := httptest.NewRecorder()
			req.AddCookie(&http.Cookie{Name: "auth_token", Value: tt.tokenString})
			wrapperHandler.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.statusCode, res.StatusCode, "status code is not equal")
		})
	}
}
