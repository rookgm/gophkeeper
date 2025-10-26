package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rookgm/gophkeeper/internal/models"
	"github.com/rookgm/gophkeeper/internal/server/repository/postgres"
)

const (
	insertSecretQuery = `
							INSERT INTO secrets (user_id, name, type, note, data) 
							VALUES ($1, $2, $3, $4, $5)
							RETURNING id, user_id, name, type, note, data, created_at, updated_at`

	selectSecretByIDQuery = `
							SELECT id, user_id, name, type, note, data, created_at, updated_at FROM secrets
							WHERE id = $1 and user_id = $2 and deleted=false
`
	deleteSecretQuery = `
							UPDATE secrets SET deleted=true WHERE id = $1 AND user_id=$2 
`
	updateSecretQuery = `
							UPDATE secrets
							SET name = $3 type=$4 note=$5 data=$6 updated_at=$7
							WHERE id = $1 and user_id = $2
	`
)

type SecretRepository struct {
	db *postgres.DB
}

// NewSecretRepository creates new secret repository
func NewSecretRepository(db *postgres.DB) *SecretRepository {
	return &SecretRepository{db: db}
}

// CreateSecret adds the secret to the repository
func (r *SecretRepository) CreateSecret(ctx context.Context, sec *models.Secret) (*models.Secret, error) {
	var res models.Secret
	err := r.db.QueryRow(ctx, insertSecretQuery, sec.UserID, sec.Name, sec.Type, sec.Note, sec.Data).Scan(&res.ID, &res.UserID, &res.Name, &res.Type, &res.Note, &res.Data, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrDataNotFound
		}
		return nil, err
	}
	return &res, nil
}

// GetSecretByID get secret from repository
func (r *SecretRepository) GetSecretByID(ctx context.Context, secretID, userID uuid.UUID) (*models.Secret, error) {
	var res models.Secret
	err := r.db.QueryRow(ctx, selectSecretByIDQuery, secretID, userID).Scan(&res.ID, &res.UserID, &res.Name, &res.Type, &res.Note, &res.Data, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrDataNotFound
		}
		return nil, err
	}
	return &res, nil
}

// DeleteSecretByID marks secret as deleted and returns secret information
func (r *SecretRepository) DeleteSecretByID(ctx context.Context, secretID, userID uuid.UUID) error {
	cmd, err := r.db.Exec(ctx, deleteSecretQuery, secretID, userID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return models.ErrDataNotFound
	}
	return nil
}

// UpdateSecretByID updates secret information
func (r *SecretRepository) UpdateSecretByID(ctx context.Context, secretID uuid.UUID, sec *models.Secret) error {
	cmd, err := r.db.Exec(ctx, updateSecretQuery, secretID, sec.UserID, sec.Name, sec.Type, sec.Note, sec.Data, sec.UpdatedAt)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return models.ErrDataNotFound
	}
	return nil
}
