package models

import (
	"github.com/google/uuid"
	"time"
)

type SecretType int

const (
	Credential SecretType = iota
	Text
	Binary
	Card
)

// return string of SecretType
func (st SecretType) String() string {
	switch st {
	case Credential:
		return "credentials"
	case Text:
		return "text"
	case Binary:
		return "binary"
	case Card:
		return "card"
	default:
		return "unknown"
	}
}

// Secret is secret entity
type Secret struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Type      SecretType
	Note      string
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Credentials struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Note     string    `json:"note,omitempty"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
}

type TextData struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Note    string    `json:"note,omitempty"`
	Content string    `json:"content"`
}

type BinaryData struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	FileName string    `json:"file_name"`
	Note     string    `json:"note,omitempty"`
	Data     []byte    `json:"data"`
}

type BankCard struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Note            string    `json:"note,omitempty"`
	CardNumber      string    `json:"cardNumber"`      // CardNumber is full bank card number
	ExpirationMonth string    `json:"expirationMonth"` // ExpirationMonth is two-digit expiration month of the card.
	ExpirationYear  string    `json:"expirationYear"`  // ExpirationYear is four-digit expiration year of the card
	CardHolderName  string    `json:"cardHolderName"`  // CardHolderName is name of the cardholder as it appears on the card
	Cvv             string    `json:"cvv"`             // Cvv is Card Verification Value (CVV) or Card Security Code (CSC).
	BillingAddress  string    `json:"billingAddress"`  // BillingAddress an object containing details of the cardholder's billing address
	CardType        string    `json:"cardType"`        // CardType is type of bank card (e.g., Visa, Mastercard, American Express)
	IssuingBank     string    `json:"issuingBank"`     // IssuingBank name of the bank that issued the bank card
}

type SecretRequest struct {
	Name string     `json:"name"`
	Type SecretType `json:"type"`
	Note string     `json:"note,omitempty"`
	Data []byte     `json:"data"`
}

type SecretResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Type      SecretType `json:"type"`
	Note      string     `json:"note,omitempty"`
	Data      []byte     `json:"data"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
