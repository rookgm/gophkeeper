package crypto

import (
	"bytes"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("Error generating salt: %v", err)
	}

	if len(salt) != saltLength {
		t.Fatalf("salt length does not match")
	}

	salt2, err := GenerateSalt()
	if err != nil {
		t.Fatalf("Error generating salt: %v", err)
	}

	if bytes.Equal(salt, salt2) {
		t.Fatalf("salts should be different")
	}
}

func TestDeriveKey(t *testing.T) {
	password := "secret-password"
	salt := []byte("the salt length must be 32 bytes")

	key := DeriveKey(password, salt)
	if len(key) != keyLength {
		t.Fatalf("key length got = %v, want %v", len(key), keyLength)
	}

	key2 := DeriveKey(password, salt)
	if !bytes.Equal(key, key2) {
		t.Fatalf("key does not match")
	}

	key3 := DeriveKey("not-secret-password", salt)
	if bytes.Equal(key, key3) {
		t.Fatalf("key should be different")
	}

	salt2 := []byte("new salt length must be 32 bytes")
	key4 := DeriveKey(password, salt2)
	if bytes.Equal(key, key4) {
		t.Fatalf("key should be different")
	}
}

func TestAESEncryptor_EncryptDecrypt(t *testing.T) {
	key := []byte("the key length must be 32 bytes!")
	text := []byte("hello world!!!")

	encryptor := AESEncryptor{}

	enc, err := encryptor.Encrypt(text, key)
	if err != nil {
		t.Fatalf("Error encrypting data: %v", err)
	}

	dec, err := encryptor.Decrypt(enc, key)
	if err != nil {
		t.Errorf("Error decrypting data: %v", err)
	}
	if !bytes.Equal(text, dec) {
		t.Errorf("Decrypted text does not match")
	}
}

func TestAESEncryptor_WrongPassword(t *testing.T) {
	key := []byte("the key length must be 32 bytes!")
	key2 := []byte("new key length must be 32 bytes!")
	text := []byte("hello world!!!")

	encryptor := AESEncryptor{}

	enc, err := encryptor.Encrypt(text, key)
	if err != nil {
		t.Fatalf("Error encrypting data: %v", err)
	}

	_, err = encryptor.Decrypt(enc, key2)
	if err == nil {
		t.Fatalf("Cannot decrypt data(wrong key): %v", err)
	}
}

func TestAESEncryptor_EncryptDecryptPwd(t *testing.T) {
	password := "secret-password-as-string"
	plaintext := []byte("hello world!!!")

	encryptor := AESEncryptor{}

	enc, err := encryptor.EncryptPwd(plaintext, password)
	if err != nil {
		t.Fatalf("Error encrypting data: %v", err)
	}

	dec, err := encryptor.DecryptPwd(enc, password)
	if err != nil {
		t.Fatalf("Error decrypting data: %v", err)
	}

	if !bytes.Equal(plaintext, dec) {
		t.Fatalf("Decrypted text does not match")
	}
}
