package models

type BankCard struct {
	CardNumber      string `json:"cardNumber"`      // CardNumber is full bank card number
	ExpirationMonth string `json:"expirationMonth"` // ExpirationMonth is two-digit expiration month of the card.
	ExpirationYear  string `json:"expirationYear"`  // ExpirationYear is four-digit expiration year of the card
	CardHolderName  string `json:"cardHolderName"`  // CardHolderName is name of the cardholder as it appears on the card
	Cvv             string `json:"cvv"`             // Cvv is Card Verification Value (CVV) or Card Security Code (CSC).
	BillingAddress  string `json:"billingAddress"`  // BillingAddress an object containing details of the cardholder's billing address
	CardType        string `json:"cardType"`        // CardType is type of bank card (e.g., Visa, Mastercard, American Express)
	IssuingBank     string `json:"issuingBank"`     // IssuingBank name of the bank that issued the bank card
}
