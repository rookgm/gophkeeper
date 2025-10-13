package models

type CreditCard struct {
	// CardNumber is full credit card number
	CardNumber string `json:"cardNumber"`
	// ExpirationMonth is two-digit expiration month of the card.
	ExpirationMonth string `json:"expirationMonth"`
	// ExpirationYear is four-digit expiration year of the card
	ExpirationYear string `json:"expirationYear"`
	// CardHolderName is name of the cardholder as it appears on the card
	CardHolderName string `json:"cardHolderName"`
	// Cvv is Card Verification Value (CVV) or Card Security Code (CSC).
	// It is typically only used during transaction processing.
	Cvv string `json:"cvv"`
	// BillingAddress an object containing details of the cardholder's billing address
	BillingAddress string `json:"billingAddress"`
	// CardType is type of credit card (e.g., Visa, Mastercard, American Express)
	CardType string `json:"cardType"`
	// IssuingBank name of the bank that issued the credit card
	IssuingBank string `json:"issuingBank"`
}
