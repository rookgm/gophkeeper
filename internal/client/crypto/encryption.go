package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
)

const (
	keyLength  = 32
	iterations = 1000
	saltLength = 32
)

var ErrCiphertext = errors.New("ciphertext too short")

type AESEncryptor struct{}

// NewAESEncryptor creates new AES encryptor
func NewAESEncryptor() *AESEncryptor {
	return &AESEncryptor{}
}

// GenerateSalt creates a new random salt of predefined length using a cryptographically secure random number generator.
// Returns the generated salt or an error if the random generation fails.
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// DeriveKey generates a cryptographic key from a password and salt using PBKDF2 with SHA-256.
// Returns a key of length defined by keyLength constant.
func DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, iterations, keyLength, sha256.New)
}

// Encrypt encrypts data using AES-GSM with provided key
func (e *AESEncryptor) Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-GSM with provided key
func (e *AESEncryptor) Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrCiphertext
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptPwd encrypts data using AES-GSM with provided password string
func (e *AESEncryptor) EncryptPwd(plaintext []byte, password string) ([]byte, error) {
	salt, err := GenerateSalt()
	if err != nil {
		return nil, err
	}

	key := DeriveKey(password, salt)

	enc, err := e.Encrypt(plaintext, key)
	if err != nil {
		return nil, err
	}
	// adding salt to the beginning of the ciphertext
	final := append(salt, enc...)

	return final, nil
}

// DecryptPwd decrypts ciphertext using AES-GSM with provided password string
func (e *AESEncryptor) DecryptPwd(ciphertext []byte, password string) ([]byte, error) {
	if len(ciphertext) < saltLength {
		return nil, ErrCiphertext
	}
	// extract salt
	salt := ciphertext[:saltLength]
	// extract ciphertext
	ctext := ciphertext[saltLength:]

	key := DeriveKey(password, salt)

	dec, err := e.Decrypt(ctext, key)
	if err != nil {
		return nil, err
	}

	return dec, nil
}
