package api

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rookgm/gophkeeper/internal/models"
)

// CreateSecret performs create secret request
func (c *Client) CreateSecret(ctx context.Context, req models.SecretRequest, token string) (*models.SecretResponse, error) {
	// POST /api/user/secrets
	resp := &models.SecretResponse{}
	err := c.doRequest(ctx, "POST", "/api/user/secrets", &token, req, resp)
	if err != nil {
		return nil, fmt.Errorf("Create secret failed: %w\n", err)
	}
	return resp, nil
}

// GetSecret performs get secret request
func (c *Client) GetSecret(ctx context.Context, id uuid.UUID, token string) (*models.SecretResponse, error) {
	// GET /api/user/secrets/{id}
	resp := &models.SecretResponse{}
	err := c.doRequest(ctx, "GET", "/api/user/secrets/"+id.String(), &token, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("Get secret failed: %w\n", err)
	}
	return resp, nil
}

// DeleteSecret performs delete secret request
func (c *Client) DeleteSecret(ctx context.Context, id uuid.UUID, token string) error {
	// DELETE /api/user/secrets/{id}
	err := c.doRequest(ctx, "DELETE", "/api/user/secrets/"+id.String(), &token, nil, nil)
	if err != nil {
		fmt.Errorf("Delete secret failed: %w\n", err)
	}
	return nil
}

// UpdateSecret performs update secret request
func (c *Client) UpdateSecret(ctx context.Context, id uuid.UUID, req models.SecretRequest, token string) error {
	// PUT /api/user/secrets/{id}
	resp := &models.SecretResponse{}
	err := c.doRequest(ctx, "PUT", "/api/user/secrets/"+id.String(), &token, req, resp)
	if err != nil {
		return fmt.Errorf("Update secret failed: %w\n", err)
	}
	return nil
}

// Sync performs syncing data
func (c *Client) Sync(ctx context.Context) {
	// POST /api/user/secrets/sync
}
