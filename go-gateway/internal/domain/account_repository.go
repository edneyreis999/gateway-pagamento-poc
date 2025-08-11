package domain

import "context"

// AccountRepository defines persistence operations for Account.
type AccountRepository interface {
	Create(ctx context.Context, a *Account) error
	GetByID(ctx context.Context, id string) (*Account, error)
	GetByAPIKey(ctx context.Context, apiKey string) (*Account, error)
	UpdateBalance(ctx context.Context, id string, amount float64) error
}

// Domain-level errors for repository implementations.
var (
	ErrAccountNotFound = Err("account: not found")
)

type Err string

func (e Err) Error() string { return string(e) }
