package domain

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidAmount        = errors.New("invoice: invalid amount")
	ErrInvalidDescription   = errors.New("invoice: invalid description")
	ErrInvalidPaymentType   = errors.New("invoice: invalid payment type")
	ErrInvalidStatus        = errors.New("invoice: invalid status")
	ErrInvoiceNegativeValue = errors.New("invoice: amount must be positive")
)

// Status represents the possible states of an invoice
type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

// Invoice represents a payment invoice that belongs to an account
type Invoice struct {
	ID             string
	AccountID      string
	Amount         float64
	Status         Status
	Description    string
	PaymentType    string
	CardLastDigits string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	mu             sync.RWMutex
}

// NewInvoice creates a new Invoice with generated ID and timestamps.
func NewInvoice(accountID, description, paymentType string, amount float64, cardLastDigits string) (*Invoice, error) {
	if len(accountID) == 0 {
		return nil, errors.New("invoice: account ID is required")
	}
	if len(description) < 3 {
		return nil, ErrInvalidDescription
	}
	if len(paymentType) == 0 {
		return nil, ErrInvalidPaymentType
	}
	if amount <= 0 {
		return nil, ErrInvoiceNegativeValue
	}

	now := time.Now().UTC()
	return &Invoice{
		ID:             uuid.New().String(),
		AccountID:      accountID,
		Amount:         amount,
		Status:         StatusPending,
		Description:    description,
		PaymentType:    paymentType,
		CardLastDigits: cardLastDigits,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// Process updates the invoice status and emits events for processing
func (i *Invoice) Process() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.Status != StatusPending {
		return errors.New("invoice: can only process pending invoices")
	}

	randomSource := rand.New(rand.NewSource(time.Now().Unix()))
	var newStatus Status

	if randomSource.Float64() <= 0.7 {
		newStatus = StatusApproved
	} else {
		newStatus = StatusRejected
	}

	i.Status = newStatus

	// In a real implementation, this would emit events
	// For now, we just update the timestamp
	i.UpdatedAt = time.Now().UTC()

	return nil
}

// UpdateStatus updates the invoice status
func (i *Invoice) UpdateStatus(newStatus Status) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	switch newStatus {
	case StatusPending, StatusApproved, StatusRejected:
		i.Status = newStatus
		i.UpdatedAt = time.Now().UTC()
		return nil
	default:
		return ErrInvalidStatus
	}
}

// IsPending checks if the invoice is in pending status
func (i *Invoice) IsPending() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.Status == StatusPending
}

// IsApproved checks if the invoice is approved
func (i *Invoice) IsApproved() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.Status == StatusApproved
}

// IsRejected checks if the invoice is rejected
func (i *Invoice) IsRejected() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.Status == StatusRejected
}
