package domain

import "context"

// InvoiceRepository defines persistence operations for Invoice.
type InvoiceRepository interface {
	Create(ctx context.Context, i *Invoice) error
	GetByID(ctx context.Context, id string) (*Invoice, error)
	GetByAccountID(ctx context.Context, accountID string) ([]*Invoice, error)
	UpdateStatus(ctx context.Context, id string, status Status) error
}

// Domain-level errors for repository implementations.
var (
	ErrInvoiceNotFound = Err("invoice: not found")
)
