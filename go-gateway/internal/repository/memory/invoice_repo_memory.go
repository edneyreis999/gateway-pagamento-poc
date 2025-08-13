package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
)

// InvoiceRepositoryMemory implements domain.InvoiceRepository using in-memory storage.
type InvoiceRepositoryMemory struct {
	invoices map[string]*domain.Invoice
	mu       sync.RWMutex
}

// NewInvoiceRepositoryMemory creates a new in-memory invoice repository.
func NewInvoiceRepositoryMemory() *InvoiceRepositoryMemory {
	return &InvoiceRepositoryMemory{
		invoices: make(map[string]*domain.Invoice),
	}
}

// Create stores a new invoice in memory.
func (r *InvoiceRepositoryMemory) Create(ctx context.Context, i *domain.Invoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.invoices[i.ID]; exists {
		return errors.New("invoice: already exists")
	}

	// Create a copy to avoid external modifications
	invoiceCopy := *i
	r.invoices[i.ID] = &invoiceCopy
	return nil
}

// GetByID retrieves an invoice by its ID.
func (r *InvoiceRepositoryMemory) GetByID(ctx context.Context, id string) (*domain.Invoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	invoice, exists := r.invoices[id]
	if !exists {
		return nil, domain.ErrInvoiceNotFound
	}

	// Return a copy to avoid external modifications
	invoiceCopy := *invoice
	return &invoiceCopy, nil
}

// GetByAccountID retrieves all invoices for a specific account.
func (r *InvoiceRepositoryMemory) GetByAccountID(ctx context.Context, accountID string) ([]*domain.Invoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var invoices []*domain.Invoice
	for _, invoice := range r.invoices {
		if invoice.AccountID == accountID {
			// Create a copy to avoid external modifications
			invoiceCopy := *invoice
			invoices = append(invoices, &invoiceCopy)
		}
	}

	return invoices, nil
}

// UpdateStatus updates the status of an existing invoice.
func (r *InvoiceRepositoryMemory) UpdateStatus(ctx context.Context, id string, status domain.Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	invoice, exists := r.invoices[id]
	if !exists {
		return domain.ErrInvoiceNotFound
	}

	return invoice.UpdateStatus(status)
}
