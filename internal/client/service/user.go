package service

import (
	"context"
	"fmt"
	"github.com/rookgm/gophkeeper/internal/models"
)

type APIClient interface {
	Register(ctx context.Context, req models.RegisterRequest) error
	Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
}

type Tokener interface {
	Save(token string) error
	Load() (string, error)
	Remove() error
}

// UserService implements UserService interface
type UserService struct {
	apiClient APIClient
	tokener   Tokener
}

func NewUserService(apiClient APIClient, tokener Tokener) *UserService {
	return &UserService{apiClient: apiClient, tokener: tokener}
}

// RegisterUser is registers new user
func (us *UserService) RegisterUser(ctx context.Context, user, password string) error {
	req := models.RegisterRequest{
		Login:    user,
		Password: password,
	}
	// do register user
	err := us.apiClient.Register(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

// LoginUser performs user logging
func (us *UserService) LoginUser(ctx context.Context, user, password string) error {
	req := models.LoginRequest{
		Login:    user,
		Password: password,
	}
	// do login user
	resp, err := us.apiClient.Login(ctx, req)
	if err != nil {
		return err
	}
	// save token to file
	if err := us.tokener.Save(resp.Token); err != nil {
		return fmt.Errorf("error saving token: %v", err)
	}

	return nil
}
