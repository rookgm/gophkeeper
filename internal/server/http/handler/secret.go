package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rookgm/gophkeeper/internal/models"
	"github.com/rookgm/gophkeeper/internal/server/http/middleware"
	"net/http"
)

// SecretService is interface for interacting with secret service
type SecretService interface {
	CreateSecret(ctx context.Context, sec *models.Secret) (*models.Secret, error)
	GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*models.Secret, error)
	UpdateSecret(ctx context.Context, secretID uuid.UUID, sec *models.Secret) error
	DeleteSecret(ctx context.Context, secretID, userID uuid.UUID) error
}

// SecretHandler represents HTTP handler for secret-related requests
type SecretHandler struct {
	svc SecretService
}

// NewSecretHandler creates new secret handler instance
func NewSecretHandler(svc SecretService) *SecretHandler {
	return &SecretHandler{svc: svc}
}

// CreateUserSecret creates new user secret
//
// POST /api/user/secrets
//
// code status
// 201 - secret is created successfully;
// 400 - bad request;
// 401 - user is not authorized;
// 405 - method not allowed;
// 500 - internal server error.
func (h *SecretHandler) CreateUserSecret(w http.ResponseWriter, r *http.Request) {
	// only POST method
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	// ger use id
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	var secReq models.SecretRequest

	// read secret
	if err := json.NewDecoder(r.Body).Decode(&secReq); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	sec := models.Secret{
		UserID: userID,
		Name:   secReq.Name,
		Type:   secReq.Type,
		Note:   secReq.Note,
		Data:   secReq.Data,
	}
	// create secret
	secNew, err := h.svc.CreateSecret(r.Context(), &sec)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	// forming request
	secResp := models.SecretResponse{
		ID:        secNew.ID,
		Name:      secNew.Name,
		Type:      secNew.Type,
		Note:      secNew.Note,
		Data:      secNew.Data,
		CreatedAt: secNew.CreatedAt,
		UpdatedAt: secNew.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(secResp); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

// GetUserSecret gets user secret
//
// GET /api/user/secrets/{id}
//
// code status
// 200 - secret is received successfully;
// 400 - bad request;
// 401 - user is not authorized;
// 404 - secret not found;
// 405 - method not allowed;
// 500 - internal server error.
func (h *SecretHandler) GetUserSecret(w http.ResponseWriter, r *http.Request) {
	// only GET method
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	// ger use id
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	// get secret id string
	s := chi.URLParam(r, "id")
	// parse secret id
	secretID, err := uuid.Parse(s)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// get secret
	sec, err := h.svc.GetSecret(r.Context(), secretID, userID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDataNotFound):
			http.Error(w, "secret not found", http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	// forming request
	secResp := models.SecretResponse{
		ID:        sec.ID,
		Name:      sec.Name,
		Type:      sec.Type,
		Note:      sec.Note,
		Data:      sec.Data,
		CreatedAt: sec.CreatedAt,
		UpdatedAt: sec.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(secResp); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

// DeleteUserSecret deletes user secret
//
// DELETE /api/user/secrets/{id}
//
// code status
// 201 - secret is deleted successfully;
// 400 - bad request;
// 401 - user is not authorized;
// 404 - secret not found;
// 405 - method not allowed;
// 500 - internal server error.
func (h *SecretHandler) DeleteUserSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	// ger use id
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	// get secret id string
	s := chi.URLParam(r, "id")
	// parse secret id
	secretID, err := uuid.Parse(s)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
	}
	// delete secret
	err = h.svc.DeleteSecret(r.Context(), secretID, userID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDataNotFound):
			http.Error(w, "secret not found", http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UpdateUserSecret updates secret
//
// PUT /api/user/secrets/{id}
//
// code status
// 201 - secret is created successfully;
// 400 - bad request;
// 401 - user is not authorized;
// 404 - secret not found;
// 405 - method not allowed;
// 500 - internal server error.
func (h *SecretHandler) UpdateUserSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	// ger use id
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	// get secret id string
	s := chi.URLParam(r, "id")
	// parse secret id
	secretID, err := uuid.Parse(s)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	var secReq models.SecretRequest

	// read secret
	if err := json.NewDecoder(r.Body).Decode(&secReq); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// forming secret for update
	sec := models.Secret{
		UserID: userID,
		Name:   secReq.Name,
		Type:   secReq.Type,
		Note:   secReq.Note,
		Data:   secReq.Data,
	}
	// update secret
	err = h.svc.UpdateSecret(r.Context(), secretID, &sec)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDataNotFound):
			http.Error(w, "secret not found", http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
