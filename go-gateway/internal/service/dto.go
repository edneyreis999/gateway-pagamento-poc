package service

import (
	"time"
)

// AccountCreateInput is the input DTO to create an account.
type AccountCreateInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// AccountOutput is the output DTO for account responses.
type AccountOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	APIKey    string    `json:"api_key"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddBalanceInput is the input DTO for adding balance.
type AddBalanceInput struct {
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

// InvoiceCreateInput is the input DTO to create an invoice.
type InvoiceCreateInput struct {
	APIKey         string  `json:"api_key"`
	Amount         float64 `json:"amount"`
	Description    string  `json:"description"`
	PaymentType    string  `json:"payment_type"`
	CardLastDigits string  `json:"card_last_digits,omitempty"`
}

// InvoiceOutput is the output DTO for invoice responses.
type InvoiceOutput struct {
	ID             string    `json:"id"`
	AccountID      string    `json:"account_id"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"`
	Description    string    `json:"description"`
	PaymentType    string    `json:"payment_type"`
	CardLastDigits string    `json:"card_last_digits,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
