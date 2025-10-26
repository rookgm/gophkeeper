package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/rookgm/gophkeeper/internal/models"
	"time"
)

type SecretRepository interface {
	CreateSecret(ctx context.Context, sec *models.Secret) (*models.Secret, error)
	GetSecretByID(ctx context.Context, secretID, userID uuid.UUID) (*models.Secret, error)
	DeleteSecretByID(ctx context.Context, secretID, userID uuid.UUID) error
	UpdateSecretByID(ctx context.Context, secretID uuid.UUID, sec *models.Secret) error
}

// SecretService implements SecretService interface
type SecretService struct {
	repo SecretRepository
}

// NewSecretService create new secret service instance
func NewSecretService(repo SecretRepository) *SecretService {
	return &SecretService{repo: repo}
}

// CreateSecret creates a new secret
func (s *SecretService) CreateSecret(ctx context.Context, sec *models.Secret) (*models.Secret, error) {
	return s.repo.CreateSecret(ctx, sec)
}

// GetSecret gets secret
func (s *SecretService) GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*models.Secret, error) {
	return s.repo.GetSecretByID(ctx, secretID, userID)
}

// DeleteSecret deletes secret
func (s *SecretService) DeleteSecret(ctx context.Context, secretID, userID uuid.UUID) error {
	return s.repo.DeleteSecretByID(ctx, secretID, userID)
}

// UpdateSecret updates secret
func (s *SecretService) UpdateSecret(ctx context.Context, secretID uuid.UUID, sec *models.Secret) error {
	// check existing secret
	secCur, err := s.repo.GetSecretByID(ctx, secretID, sec.UserID)
	if err != nil {
		return err
	}
	// update secret
	secCur.Name = sec.Name
	secCur.Type = sec.Type
	secCur.Note = sec.Note
	secCur.Data = sec.Data
	secCur.UpdatedAt = time.Now()

	return s.repo.UpdateSecretByID(ctx, secretID, secCur)
}
