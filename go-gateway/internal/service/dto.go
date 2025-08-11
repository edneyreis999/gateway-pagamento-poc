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
