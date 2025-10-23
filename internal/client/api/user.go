package api

import (
	"context"
	"fmt"
	"github.com/rookgm/gophkeeper/internal/models"
)

// Register performs a user creation request
func (c *Client) Register(ctx context.Context, req models.RegisterRequest) error {
	// POST /api/user/register
	err := c.doRequest(ctx, "POST", "/api/user/register", nil, req, nil)
	if err != nil {
		return fmt.Errorf("register request failed: %w", err)
	}
	return nil
}

// Login performs user logging
func (c *Client) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	resp := &models.LoginResponse{}
	// POST /api/user/login
	err := c.doRequest(ctx, "POST", "/api/user/login", nil, req, resp)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}

	return resp, nil
}
