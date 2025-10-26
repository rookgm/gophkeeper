package models

import (
	"github.com/google/uuid"
	"time"
)

// User is user entity
type User struct {
	ID        uuid.UUID
	Login     string
	Password  string
	CreatedAt time.Time
}

// RegisterRequest is user registration data
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginRequest is user login data
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginResponse is login response data
type LoginResponse struct {
	Token string `json:"token"`
}

// TokenPayload is payload contains user id's
type TokenPayload struct {
	ID     uuid.UUID
	UserID uuid.UUID
}
