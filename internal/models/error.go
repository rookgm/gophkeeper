package models

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrConflictData       = errors.New("data conflicts with existing data")
	ErrDataNotFound       = errors.New("data not found")
	ErrInvalidCredentials = errors.New("invalid login or password")
)

type TooManyRequestsError struct {
	RetryAfter time.Duration
}

func (e TooManyRequestsError) Error() string {
	return fmt.Sprintf("Error: too many requests, retry after %v", e.RetryAfter)
}

func NewTooManyRequestsError(retryAfter time.Duration) *TooManyRequestsError {
	return &TooManyRequestsError{RetryAfter: retryAfter}
}
