package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rookgm/gophkeeper/internal/models"
)

type APIClientSecreter interface {
	CreateSecret(ctx context.Context, req models.SecretRequest, token string) (*models.SecretResponse, error)
	GetSecret(ctx context.Context, id uuid.UUID, token string) (*models.SecretResponse, error)
	DeleteSecret(ctx context.Context, id uuid.UUID, token string) (*models.SecretResponse, error)
	UpdateSecret(ctx context.Context, id uuid.UUID, req models.SecretRequest, token string) (*models.SecretResponse, error)
}

type AESGSMEncryptor interface {
	EncryptPwd(plaintext []byte, password string) ([]byte, error)
	DecryptPwd(ciphertext []byte, password string) ([]byte, error)
}

// SecretService implements SecretService interface
type SecretService struct {
	apiClient APIClientSecreter
	encryptor AESGSMEncryptor
	tokener   Tokener
}

// NewSecretService creates a new SecretService instance
func NewSecretService(apiClient APIClientSecreter, encryptor AESGSMEncryptor, tokener Tokener) *SecretService {
	return &SecretService{apiClient: apiClient, encryptor: encryptor, tokener: tokener}
}

// AddCredentials adds secret credential data
func (s *SecretService) AddCredentials(ctx context.Context, req models.Credentials, masterPassword string) (*models.Credentials, error) {
	// marshaling credentials
	credData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling credentials: %v\n", err)
	}
	// encrypt data with master password
	credEnc, err := s.encryptor.EncryptPwd(credData, masterPassword)
	if err != nil {
		return nil, fmt.Errorf("Error encrypting credentials: %v\n", err)
	}
	// forming secret request
	secReq := models.SecretRequest{
		Name: req.Name,
		Type: models.Credential,
		Note: req.Note,
		Data: credEnc,
	}
	// load token
	token, err := s.tokener.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading token: %v\n", err)
	}
	// sending a request to create a secret
	respSec, err := s.apiClient.CreateSecret(ctx, secReq, token)
	if err != nil {
		return nil, fmt.Errorf("Error creating secret: %v\n", err)
	}

	return &models.Credentials{
		ID:   respSec.ID,
		Name: respSec.Name,
		Note: respSec.Note,
	}, nil
}

// AddText adds secret text data
func (s *SecretService) AddText(ctx context.Context, req models.TextData, masterPassword string) (*models.TextData, error) {
	// marshaling text data
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling text data: %v\n", err)
	}
	// encrypt text data
	textEnc, err := s.encryptor.EncryptPwd(data, masterPassword)
	if err != nil {
		return nil, fmt.Errorf("Error encrypting text data: %v\n", err)
	}
	// forming secret request
	secReq := models.SecretRequest{
		Name: req.Name,
		Type: models.Text,
		Note: req.Note,
		Data: textEnc,
	}
	// load token
	token, err := s.tokener.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading token: %v\n", err)
	}
	// sending a request to create a secret
	respSec, err := s.apiClient.CreateSecret(ctx, secReq, token)
	if err != nil {
		return nil, fmt.Errorf("Error creating secret: %v\n", err)
	}
	return &models.TextData{
		ID:   respSec.ID,
		Name: respSec.Name,
		Note: respSec.Note,
	}, nil
}

// AddBinary adds secret binary data
func (s *SecretService) AddBinary(ctx context.Context, req models.BinaryData, masterPassword string) (*models.BinaryData, error) {
	// marshaling binary data
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling binary data: %v\n", err)
	}
	// encrypt binary data
	binEnc, err := s.encryptor.EncryptPwd(data, masterPassword)
	if err != nil {
		return nil, fmt.Errorf("Error encrypting text data: %v\n", err)
	}
	// forming secret request
	secReq := models.SecretRequest{
		Name: req.Name,
		Type: models.Credential,
		Note: req.Note,
		Data: binEnc,
	}
	// load token
	token, err := s.tokener.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading token: %v\n", err)
	}
	// sending a request to create a secret
	respSec, err := s.apiClient.CreateSecret(ctx, secReq, token)
	if err != nil {
		return nil, fmt.Errorf("Error creating secret: %v\n", err)
	}
	return &models.BinaryData{
		ID:   respSec.ID,
		Name: respSec.Name,
		Note: respSec.Note,
	}, nil
}

// AddBankCard adds secret bank data
func (s *SecretService) AddBankCard(ctx context.Context, req models.BankCard, masterPassword string) (*models.BankCard, error) {
	// marshaling bank card
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling bank card: %v\n", err)
	}
	// encrypt bank card
	cardEnc, err := s.encryptor.EncryptPwd(data, masterPassword)
	if err != nil {
		return nil, fmt.Errorf("Error encrypting text data: %v\n", err)
	}
	// forming secret request
	secReq := models.SecretRequest{
		Name: req.Name,
		Type: models.Credential,
		Note: req.Note,
		Data: cardEnc,
	}
	// load token
	token, err := s.tokener.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading token: %v\n", err)
	}
	// sending a request to create a secret
	respSec, err := s.apiClient.CreateSecret(ctx, secReq, token)
	if err != nil {
		return nil, fmt.Errorf("Error creating secret: %v\n", err)
	}
	return &models.BankCard{
		ID:   respSec.ID,
		Name: respSec.Name,
		Note: respSec.Note,
	}, nil
}

// GetSecret gets secret data
func (s *SecretService) GetSecret(ctx context.Context, id uuid.UUID, masterPassword string) (*models.SecretResponse, error) {
	token, err := s.tokener.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading token: %v\n", err)
	}
	// get secret with encrypted data
	respEnc, err := s.apiClient.GetSecret(ctx, id, token)
	if err != nil {
		return nil, fmt.Errorf("Error getting secret: %v\n", err)
	}
	// decrypt data
	dec, err := s.encryptor.DecryptPwd(respEnc.Data, masterPassword)
	if err != nil {
		return nil, fmt.Errorf("Error decrypting secret data: %v\n", err)
	}
	return &models.SecretResponse{
		ID:        respEnc.ID,
		Name:      respEnc.Name,
		Type:      respEnc.Type,
		Note:      respEnc.Note,
		Data:      dec,
		CreatedAt: respEnc.CreatedAt,
		UpdatedAt: respEnc.UpdatedAt,
	}, nil
}

// DeleteSecret deletes secret
func (s *SecretService) DeleteSecret(ctx context.Context, id uuid.UUID, masterPassword string) (*models.SecretResponse, error) {
	token, err := s.tokener.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading token: %v\n", err)
	}
	// get secret with encrypted data
	respEnc, err := s.apiClient.DeleteSecret(ctx, id, token)
	if err != nil {
		return nil, fmt.Errorf("Error getting secret: %v\n", err)
	}
	// decrypt data
	dec, err := s.encryptor.DecryptPwd(respEnc.Data, masterPassword)
	if err != nil {
		return nil, fmt.Errorf("Error decrypting secret data: %v\n", err)
	}
	return &models.SecretResponse{
		ID:        respEnc.ID,
		Name:      respEnc.Name,
		Type:      respEnc.Type,
		Note:      respEnc.Note,
		Data:      dec,
		CreatedAt: respEnc.CreatedAt,
		UpdatedAt: respEnc.UpdatedAt,
	}, nil
}

func (s *SecretService) UpdateSecret(ctx context.Context, id uuid.UUID, req models.SecretRequest, masterPassword string) (*models.SecretResponse, error) {
	// encrypt data
	dataEnc, err := s.encryptor.EncryptPwd(req.Data, masterPassword)
	if err != nil {
		return nil, fmt.Errorf("Error encrypting secret data: %v\n", err)
	}
	// load token
	token, err := s.tokener.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading token: %v\n", err)
	}

	reqEnc := models.SecretRequest{
		Name: req.Name,
		Type: req.Type,
		Note: req.Note,
		Data: dataEnc,
	}

	respEnc, err := s.apiClient.UpdateSecret(ctx, id, reqEnc, token)
	if err != nil {
		return nil, fmt.Errorf("Error updating secret: %v\n", err)
	}
	dec, err := s.encryptor.DecryptPwd(respEnc.Data, masterPassword)
	if err != nil {
		return nil, fmt.Errorf("Error decrypting secret data: %v\n", err)
	}

	return &models.SecretResponse{
			ID:        respEnc.ID,
			Name:      respEnc.Name,
			Type:      respEnc.Type,
			Note:      respEnc.Note,
			Data:      dec,
			CreatedAt: respEnc.CreatedAt,
			UpdatedAt: respEnc.UpdatedAt,
		},
		nil
}
