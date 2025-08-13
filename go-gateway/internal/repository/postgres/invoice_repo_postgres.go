package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
)

// PostgresInvoiceRepository implements domain.InvoiceRepository using PostgreSQL.
type PostgresInvoiceRepository struct {
	db *sql.DB
}

// NewPostgresInvoiceRepository creates a new PostgreSQL invoice repository.
func NewPostgresInvoiceRepository(db *sql.DB) *PostgresInvoiceRepository {
	return &PostgresInvoiceRepository{db: db}
}

// Create stores a new invoice in PostgreSQL.
func (r *PostgresInvoiceRepository) Create(ctx context.Context, i *domain.Invoice) error {
	query := `
		INSERT INTO invoices (id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		i.ID, i.AccountID, i.Amount, i.Status, i.Description, i.PaymentType, i.CardLastDigits, i.CreatedAt, i.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves an invoice by its ID from PostgreSQL.
func (r *PostgresInvoiceRepository) GetByID(ctx context.Context, id string) (*domain.Invoice, error) {
	query := `
		SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at
		FROM invoices
		WHERE id = $1
	`

	var invoice domain.Invoice
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&invoice.ID, &invoice.AccountID, &invoice.Amount, &invoice.Status, &invoice.Description,
		&invoice.PaymentType, &invoice.CardLastDigits, &invoice.CreatedAt, &invoice.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrInvoiceNotFound
		}
		return nil, err
	}

	return &invoice, nil
}

// GetByAccountID retrieves all invoices for a specific account from PostgreSQL.
func (r *PostgresInvoiceRepository) GetByAccountID(ctx context.Context, accountID string) ([]*domain.Invoice, error) {
	query := `
		SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at
		FROM invoices
		WHERE account_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []*domain.Invoice
	for rows.Next() {
		var invoice domain.Invoice
		err := rows.Scan(
			&invoice.ID, &invoice.AccountID, &invoice.Amount, &invoice.Status, &invoice.Description,
			&invoice.PaymentType, &invoice.CardLastDigits, &invoice.CreatedAt, &invoice.UpdatedAt)

		if err != nil {
			return nil, err
		}

		invoices = append(invoices, &invoice)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return invoices, nil
}

// UpdateStatus updates the status of an existing invoice in PostgreSQL.
func (r *PostgresInvoiceRepository) UpdateStatus(ctx context.Context, id string, status domain.Status) error {
	query := `
		UPDATE invoices
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrInvoiceNotFound
	}

	return nil
}
