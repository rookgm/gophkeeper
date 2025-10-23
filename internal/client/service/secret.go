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
	UpdateSecret(ctx context.Context, id uuid.UUID, token string) (*models.SecretResponse, error)
	DeleteSecret(ctx context.Context, id uuid.UUID, token string) error
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
func (c *SecretService) AddCredentials(ctx context.Context, req models.Credentials, masterPassword string) (*models.Credentials, error) {
	// marshaling credentials
	credData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling credentials: %v\n", err)
	}
	// encrypt data with master password
	credEnc, err := c.encryptor.EncryptPwd(credData, masterPassword)
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
	token, err := c.tokener.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading token: %v\n", err)
	}
	// sending a request to create a secret
	respSec, err := c.apiClient.CreateSecret(ctx, secReq, token)
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
func (c *SecretService) AddText(ctx context.Context, req models.TextData, masterPassword string) (*models.TextData, error) {
	// marshaling text data
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling text data: %v\n", err)
	}
	// encrypt text data
	textEnc, err := c.encryptor.EncryptPwd(data, masterPassword)
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
	token, err := c.tokener.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading token: %v\n", err)
	}
	// sending a request to create a secret
	respSec, err := c.apiClient.CreateSecret(ctx, secReq, token)
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
func (c *SecretService) AddBinary(ctx context.Context, req models.BinaryData, masterPassword string) (*models.BinaryData, error) {
	// marshaling binary data
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling binary data: %v\n", err)
	}
	// encrypt binary data
	binEnc, err := c.encryptor.EncryptPwd(data, masterPassword)
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
	token, err := c.tokener.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading token: %v\n", err)
	}
	// sending a request to create a secret
	respSec, err := c.apiClient.CreateSecret(ctx, secReq, token)
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
func (c *SecretService) AddBankCard(ctx context.Context, req models.BankCard, masterPassword string) (*models.BankCard, error) {
	// marshaling bank card
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling bank card: %v\n", err)
	}
	// encrypt bank card
	cardEnc, err := c.encryptor.EncryptPwd(data, masterPassword)
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
	token, err := c.tokener.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading token: %v\n", err)
	}
	// sending a request to create a secret
	respSec, err := c.apiClient.CreateSecret(ctx, secReq, token)
	if err != nil {
		return nil, fmt.Errorf("Error creating secret: %v\n", err)
	}
	return &models.BankCard{
		ID:   respSec.ID,
		Name: respSec.Name,
		Note: respSec.Note,
	}, nil
}

func (c *SecretService) DeleteSecret(ctx context.Context, id uuid.UUID, token string) error {
	return nil
}

func (c *SecretService) GetSecret(ctx context.Context, id uuid.UUID, token string) (*models.Secret, error) {
	return nil, nil
}

func (c *SecretService) ListSecrets(ctx context.Context, token string) ([]models.Secret, error) {
	return nil, nil
}
