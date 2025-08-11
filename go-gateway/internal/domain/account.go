package domain

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidName   = errors.New("account: invalid name")
	ErrInvalidEmail  = errors.New("account: invalid email")
	ErrNegativeValue = errors.New("account: amount must be positive")
)

// Account represents a client account that owns invoices and holds a balance
// increased when invoices are approved.
type Account struct {
	ID        string
	Name      string
	Email     string
	APIKey    string
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
	mu        sync.RWMutex
}

// NewAccount creates a new Account with generated IDs and timestamps.
func NewAccount(name, email string) (*Account, error) {
	if len(name) < 2 {
		return nil, ErrInvalidName
	}
	// NOTE: keeping email validation minimal for now; can be enhanced.
	if len(email) < 5 || !containsAt(email) {
		return nil, ErrInvalidEmail
	}
	now := time.Now().UTC()
	return &Account{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		APIKey:    uuid.New().String(),
		Balance:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// AddBalance increments the balance by a positive amount and updates UpdatedAt.
func (a *Account) AddBalance(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if amount <= 0 {
		return ErrNegativeValue
	}
	a.Balance += amount
	a.UpdatedAt = time.Now().UTC()
	return nil
}

// containsAt is a small helper to avoid importing regexp now.
func containsAt(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '@' {
			return true
		}
	}
	return false
}
