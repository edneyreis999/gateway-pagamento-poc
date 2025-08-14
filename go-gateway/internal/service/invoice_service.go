package service

import (
	"context"
	"database/sql"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	pg "github.com/devfullcycle/imersao22/go-gateway/internal/repository/postgres"
)

// AccountServicePort defines the interface for AccountService methods needed by InvoiceService
type AccountServicePort interface {
	GetByAPIKey(ctx context.Context, apiKey string) (*AccountOutput, error)
}

// InvoiceService implements domain.InvoiceRepository by delegating to a Postgres repository
// and also provides DTO-based methods for the API/handlers layer.
type InvoiceService struct {
	repo           domain.InvoiceRepository
	accountService AccountServicePort
	processor      domain.InvoiceProcessor // Custom processor for testing
}

func NewInvoiceService(db *sql.DB) *InvoiceService {
	return &InvoiceService{
		repo:           pg.NewPostgresInvoiceRepository(db),
		accountService: NewAccountService(db),
		processor:      nil, // Use default processor
	}
}

// NewInvoiceServiceWithAccountService creates a new InvoiceService with a custom AccountService
func NewInvoiceServiceWithAccountService(db *sql.DB, accountService AccountServicePort) *InvoiceService {
	return &InvoiceService{
		repo:           pg.NewPostgresInvoiceRepository(db),
		accountService: accountService,
		processor:      nil, // Use default processor
	}
}

// SetProcessor allows setting a custom processor for testing
func (s *InvoiceService) SetProcessor(processor domain.InvoiceProcessor) {
	s.processor = processor
}

// Create creates a new invoice from input DTO and returns an output DTO.
func (s *InvoiceService) Create(ctx context.Context, in InvoiceCreateInput) (*InvoiceOutput, error) {
	accountOutput, err := s.accountService.GetByAPIKey(ctx, in.APIKey)
	if err != nil {
		return nil, err
	}

	var invoice *domain.Invoice
	var err2 error

	if s.processor != nil {
		// Use custom processor for testing
		invoice, err2 = domain.NewInvoiceWithProcessor(accountOutput.ID, in.Description, in.PaymentType, in.Amount, in.CardLastDigits, s.processor)
	} else {
		// Use default processor
		invoice, err2 = domain.NewInvoice(accountOutput.ID, in.Description, in.PaymentType, in.Amount, in.CardLastDigits)
	}

	if err2 != nil {
		return nil, err2
	}

	if err := invoice.Process(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, invoice); err != nil {
		return nil, err
	}

	return toInvoiceOutput(invoice), nil
}

// GetByID retrieves an invoice by ID and returns an output DTO.
func (s *InvoiceService) GetByID(ctx context.Context, id string) (*InvoiceOutput, error) {
	invoice, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toInvoiceOutput(invoice), nil
}

// GetByAccountID retrieves all invoices for an account and returns output DTOs.
func (s *InvoiceService) GetByAccountID(ctx context.Context, accountID string) ([]*InvoiceOutput, error) {
	invoices, err := s.repo.GetByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	var outputs []*InvoiceOutput
	for _, invoice := range invoices {
		outputs = append(outputs, toInvoiceOutput(invoice))
	}

	return outputs, nil
}

// UpdateStatus updates the status of an invoice.
func (s *InvoiceService) UpdateStatus(ctx context.Context, id string, status domain.Status) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

// toInvoiceOutput maps domain.Invoice to output DTO.
func toInvoiceOutput(i *domain.Invoice) *InvoiceOutput {
	return &InvoiceOutput{
		ID:             i.ID,
		AccountID:      i.AccountID,
		Amount:         i.Amount,
		Status:         string(i.Status),
		Description:    i.Description,
		PaymentType:    i.PaymentType,
		CardLastDigits: i.CardLastDigits,
		CreatedAt:      i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
	}
}
