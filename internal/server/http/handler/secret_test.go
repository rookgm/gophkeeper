package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/rookgm/gophkeeper/internal/models"
	"github.com/rookgm/gophkeeper/internal/server/http/handler/mocks"
	"github.com/rookgm/gophkeeper/internal/server/http/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSecretHandler_CreateUserSecret(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		body           any
		token          *models.TokenPayload
		setup          func(t *testing.T) *mocks.MockSecretService
		wantStatusCode int
		wantBody       *models.SecretResponse
	}{
		// 201 - secret is created successfully
		{
			name: "valid_request_return_201",
			body: &models.SecretRequest{
				Name: "github",
				Type: models.Credential,
				Note: "dev",
				Data: []byte("user:password"),
			},
			token: &models.TokenPayload{
				ID:     uuid.New(),
				UserID: id,
			},
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().CreateSecret(gomock.Any(), gomock.Any()).Return(&models.Secret{
					ID:        id,
					UserID:    id,
					Name:      "github",
					Type:      models.Credential,
					Note:      "dev",
					Data:      []byte("user:password"),
					CreatedAt: now,
					UpdatedAt: now,
				}, nil).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusCreated,
			wantBody: &models.SecretResponse{
				ID:        id,
				Name:      "github",
				Type:      models.Credential,
				Note:      "dev",
				Data:      []byte("user:password"),
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		// 400 - bad request
		{
			name: "bad_request_return_400",
			body: "{user:password}",
			token: &models.TokenPayload{
				ID:     id,
				UserID: id,
			},
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().CreateSecret(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       nil,
		},
		// 500 - internal server error;
		{
			name:  "internal_error_return_500(get userid)",
			body:  "{user:password}",
			token: nil,
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().CreateSecret(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       nil,
		},
		{
			name: "internal_error_return_500",
			body: &models.SecretRequest{
				Name: "github",
				Type: models.Credential,
				Note: "dev",
				Data: []byte("user:password"),
			},
			token: &models.TokenPayload{
				ID:     uuid.New(),
				UserID: id,
			},
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().CreateSecret(gomock.Any(), gomock.Any()).Return(nil, models.ErrDataNotFound).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data []byte
			var err error

			if str, ok := tt.body.(string); ok {
				data = []byte(str)
			} else {
				data, err = json.Marshal(tt.body)
				require.NoError(t, err)
			}

			bodyReader := bytes.NewReader(data)

			req, err := http.NewRequest(http.MethodPost, "/api/user/secrets", bodyReader)
			if err != nil {
				t.Fatal("failed to create http request ", err)
			}
			w := httptest.NewRecorder()
			svc := tt.setup(t)
			ctx := context.WithValue(req.Context(), middleware.AuthPayloadKey, tt.token)

			handler := NewSecretHandler(svc)
			handler.CreateUserSecret(w, req.WithContext(ctx))

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)
			if tt.wantBody != nil {
				body, err := io.ReadAll(res.Body)
				require.NoError(t, err)

				if tt.wantBody != nil {
					err = json.Unmarshal(body, tt.wantBody)
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestSecretHandler_GetUserSecret(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		target         string
		body           any
		token          *models.TokenPayload
		setup          func(t *testing.T) *mocks.MockSecretService
		wantStatusCode int
		wantBody       *models.SecretResponse
	}{
		// 200 - secret is received successfully
		{
			name:   "valid_request_return_200",
			target: "/api/user/secrets/" + id.String(),
			body: &models.SecretRequest{
				Name: "github",
				Type: models.Credential,
				Note: "dev",
				Data: []byte("user:password"),
			},
			token: &models.TokenPayload{
				ID:     uuid.New(),
				UserID: id,
			},
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().GetSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.Secret{
					ID:        id,
					UserID:    id,
					Name:      "github",
					Type:      models.Credential,
					Note:      "dev",
					Data:      []byte("user:password"),
					CreatedAt: now,
					UpdatedAt: now,
				}, nil).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusOK,
			wantBody: &models.SecretResponse{
				ID:        id,
				Name:      "github",
				Type:      models.Credential,
				Note:      "dev",
				Data:      []byte("user:password"),
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		// 400 - bad request
		{
			name:   "bad_request_return_400",
			target: "/api/user/secrets/1",
			body:   "{user:password}",
			token: &models.TokenPayload{
				ID:     id,
				UserID: id,
			},
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().GetSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       nil,
		},
		// 500 - internal server error
		{
			name:   "internal_error_return_500(get userid)",
			target: "/api/user/secrets/" + id.String(),
			body:   "{user:password}",
			token:  nil,
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().GetSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       nil,
		},
		{
			name:   "internal_error_return_500",
			target: "/api/user/secrets/" + id.String(),
			body: &models.SecretRequest{
				Name: "github",
				Type: models.Credential,
				Note: "dev",
				Data: []byte("user:password"),
			},
			token: &models.TokenPayload{
				ID:     uuid.New(),
				UserID: id,
			},
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().GetSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       nil,
		},
		{
			name:   "not_found_return_404",
			target: "/api/user/secrets/" + id.String(),
			body: &models.SecretRequest{
				Name: "github",
				Type: models.Credential,
				Note: "dev",
				Data: []byte("user:password"),
			},
			token: &models.TokenPayload{
				ID:     uuid.New(),
				UserID: id,
			},
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().GetSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, models.ErrDataNotFound).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusNotFound,
			wantBody:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data []byte
			var err error

			if str, ok := tt.body.(string); ok {
				data = []byte(str)
			} else {
				data, err = json.Marshal(tt.body)
				require.NoError(t, err)
			}

			bodyReader := bytes.NewReader(data)

			req, err := http.NewRequest(http.MethodGet, tt.target, bodyReader)
			if err != nil {
				t.Fatal("failed to create http request ", err)
			}
			req.SetPathValue("id", id.String())

			w := httptest.NewRecorder()
			svc := tt.setup(t)
			ctx := context.WithValue(req.Context(), middleware.AuthPayloadKey, tt.token)

			handler := NewSecretHandler(svc)

			router := chi.NewRouter()
			router.Get("/api/user/secrets/{id}", handler.GetUserSecret)

			router.ServeHTTP(w, req.WithContext(ctx))

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)
			if tt.wantBody != nil {
				body, err := io.ReadAll(res.Body)
				require.NoError(t, err)

				if tt.wantBody != nil {
					err = json.Unmarshal(body, tt.wantBody)
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestSecretHandler_DeleteUserSecret(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name           string
		target         string
		token          *models.TokenPayload
		setup          func(t *testing.T) *mocks.MockSecretService
		wantStatusCode int
	}{
		// 200 - secret is deleted successfully
		{
			name:   "valid_request_return_200",
			target: "/api/user/secrets/" + id.String(),
			token: &models.TokenPayload{
				ID:     uuid.New(),
				UserID: id,
			},
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().DeleteSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusOK,
		},
		// 400 - bad request
		{
			name:   "bad_request_return_400",
			target: "/api/user/secrets/1",
			token: &models.TokenPayload{
				ID:     id,
				UserID: id,
			},
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().DeleteSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusBadRequest,
		},
		// 500 - internal server error
		{
			name:   "internal_error_return_500(get userid)",
			target: "/api/user/secrets/" + id.String(),
			token:  nil,
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().DeleteSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:   "not_found_return_404",
			target: "/api/user/secrets/" + id.String(),
			token: &models.TokenPayload{
				ID:     uuid.New(),
				UserID: id,
			},
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().DeleteSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.ErrDataNotFound).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, tt.target, nil)
			if err != nil {
				t.Fatal("failed to create http request ", err)
			}
			req.SetPathValue("id", id.String())

			w := httptest.NewRecorder()
			svc := tt.setup(t)
			ctx := context.WithValue(req.Context(), middleware.AuthPayloadKey, tt.token)

			handler := NewSecretHandler(svc)

			router := chi.NewRouter()
			router.Delete("/api/user/secrets/{id}", handler.DeleteUserSecret)

			router.ServeHTTP(w, req.WithContext(ctx))

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)
		})
	}
}

func TestSecretHandler_UpdateUserSecret(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name           string
		target         string
		body           any
		token          *models.TokenPayload
		setup          func(t *testing.T) *mocks.MockSecretService
		wantStatusCode int
	}{
		// 200 - secret is received successfully
		{
			name:   "valid_request_return_200",
			target: "/api/user/secrets/" + id.String(),
			body: &models.SecretRequest{
				Name: "github",
				Type: models.Credential,
				Note: "dev",
				Data: []byte("user:password"),
			},
			token: &models.TokenPayload{
				ID:     uuid.New(),
				UserID: id,
			},
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().UpdateSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

				return svc
			},
			wantStatusCode: http.StatusOK,
		},
		// 400 - bad request
		{
			name:   "bad_request_return_400",
			target: "/api/user/secrets/1",
			body:   "{user:password}",
			token: &models.TokenPayload{
				ID:     id,
				UserID: id,
			},
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().UpdateSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusBadRequest,
		},
		// 500 - internal server error
		{
			name:   "internal_error_return_500(get userid)",
			target: "/api/user/secrets/" + id.String(),
			body:   "{user:password}",
			token:  nil,
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().UpdateSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:   "internal_error_return_500",
			target: "/api/user/secrets/" + id.String(),
			body: &models.SecretRequest{
				Name: "github",
				Type: models.Credential,
				Note: "dev",
				Data: []byte("user:password"),
			},
			token: &models.TokenPayload{
				ID:     uuid.New(),
				UserID: id,
			},
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().UpdateSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error")).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:   "not_found_return_404",
			target: "/api/user/secrets/" + id.String(),
			body: &models.SecretRequest{
				Name: "github",
				Type: models.Credential,
				Note: "dev",
				Data: []byte("user:password"),
			},
			token: &models.TokenPayload{
				ID:     uuid.New(),
				UserID: id,
			},
			setup: func(t *testing.T) *mocks.MockSecretService {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				svc := mocks.NewMockSecretService(ctrl)
				svc.EXPECT().UpdateSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.ErrDataNotFound).AnyTimes()
				return svc
			},
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data []byte
			var err error

			if str, ok := tt.body.(string); ok {
				data = []byte(str)
			} else {
				data, err = json.Marshal(tt.body)
				require.NoError(t, err)
			}

			bodyReader := bytes.NewReader(data)

			req, err := http.NewRequest(http.MethodPut, tt.target, bodyReader)
			if err != nil {
				t.Fatal("failed to create http request ", err)
			}
			req.SetPathValue("id", id.String())

			w := httptest.NewRecorder()
			svc := tt.setup(t)
			ctx := context.WithValue(req.Context(), middleware.AuthPayloadKey, tt.token)

			handler := NewSecretHandler(svc)

			router := chi.NewRouter()
			router.Put("/api/user/secrets/{id}", handler.UpdateUserSecret)

			router.ServeHTTP(w, req.WithContext(ctx))

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)
		})
	}
}
