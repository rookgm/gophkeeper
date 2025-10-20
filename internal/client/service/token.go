package service

import (
	"errors"
	"os"
)

type TokenService struct {
	name string
}

// NewToken creates a new token
func NewToken(name string) *TokenService {
	return &TokenService{name: name}
}

// Save saves token to file
func (t *TokenService) Save(token string) error {
	return os.WriteFile(t.name, []byte(token), 0600)
}

// Load reads token from file
func (t *TokenService) Load() (string, error) {
	b, err := os.ReadFile(t.name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Remove removes token file
func (t *TokenService) Remove() error {
	if _, err := os.Stat(t.name); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return os.Remove(t.name)
}
