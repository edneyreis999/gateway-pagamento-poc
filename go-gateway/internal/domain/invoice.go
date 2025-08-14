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

// TestInvoiceProcessor implements a processor for testing that allows full control
type TestInvoiceProcessor struct {
	nextStatus Status
	shouldErr  bool
	err        error
}

// NewTestInvoiceProcessor creates a new test processor
func NewTestInvoiceProcessor() *TestInvoiceProcessor {
	return &TestInvoiceProcessor{
		nextStatus: StatusApproved,
		shouldErr:  false,
		err:        nil,
	}
}

// SetNextStatus sets the next status that will be returned
func (p *TestInvoiceProcessor) SetNextStatus(status Status) {
	p.nextStatus = status
}

// SetError sets an error to be returned
func (p *TestInvoiceProcessor) SetError(err error) {
	p.shouldErr = true
	p.err = err
}

// ClearError clears any error condition
func (p *TestInvoiceProcessor) ClearError() {
	p.shouldErr = false
	p.err = nil
}

// ProcessInvoice processes an invoice with the configured test behavior
func (p *TestInvoiceProcessor) ProcessInvoice(invoice *Invoice) error {
	if p.shouldErr {
		return p.err
	}

	if invoice.Status != StatusPending {
		return errors.New("invoice: can only process pending invoices")
	}

	// Apply the same business rule: invoices with amount > 10000 stay pending
	if invoice.Amount > 10000 {
		return nil
	}

	invoice.Status = p.nextStatus
	invoice.UpdatedAt = time.Now().UTC()

	return nil
}

// InvoiceProcessor defines the interface for processing invoices
type InvoiceProcessor interface {
	ProcessInvoice(invoice *Invoice) error
}

// DefaultInvoiceProcessor implements the default random processing logic
type DefaultInvoiceProcessor struct {
	randomSource *rand.Rand
}

// NewDefaultInvoiceProcessor creates a new default processor with current time seed
func NewDefaultInvoiceProcessor() *DefaultInvoiceProcessor {
	return &DefaultInvoiceProcessor{
		randomSource: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

// NewDefaultInvoiceProcessorWithSeed creates a new default processor with a specific seed
func NewDefaultInvoiceProcessorWithSeed(seed int64) *DefaultInvoiceProcessor {
	return &DefaultInvoiceProcessor{
		randomSource: rand.New(rand.NewSource(seed)),
	}
}

// ProcessInvoice processes an invoice using random logic (70% approved, 30% rejected)
func (p *DefaultInvoiceProcessor) ProcessInvoice(invoice *Invoice) error {
	if invoice.Amount > 10000 {
		return nil
	}

	if invoice.Status != StatusPending {
		return errors.New("invoice: can only process pending invoices")
	}

	var newStatus Status
	if p.randomSource.Float64() <= 0.7 {
		newStatus = StatusApproved
	} else {
		newStatus = StatusRejected
	}

	invoice.Status = newStatus
	invoice.UpdatedAt = time.Now().UTC()

	return nil
}

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
	processor      InvoiceProcessor // Processor for this invoice
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
		processor:      NewDefaultInvoiceProcessor(),
	}, nil
}

// NewInvoiceWithProcessor creates a new Invoice with a custom processor
func NewInvoiceWithProcessor(accountID, description, paymentType string, amount float64, cardLastDigits string, processor InvoiceProcessor) (*Invoice, error) {
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
		processor:      processor,
	}, nil
}

// SetProcessor allows changing the processor for an invoice
func (i *Invoice) SetProcessor(processor InvoiceProcessor) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.processor = processor
}

// Process updates the invoice status using the configured processor
func (i *Invoice) Process() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.processor == nil {
		i.processor = NewDefaultInvoiceProcessor()
	}

	return i.processor.ProcessInvoice(i)
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
