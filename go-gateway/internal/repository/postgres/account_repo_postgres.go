package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
)

// PostgresAccountRepository implements AccountRepository using database/sql.
type PostgresAccountRepository struct {
	db *sql.DB
}

func NewPostgresAccountRepository(db *sql.DB) *PostgresAccountRepository {
	return &PostgresAccountRepository{db: db}
}

func (r *PostgresAccountRepository) Create(ctx context.Context, a *domain.Account) error {
	const q = `
		INSERT INTO accounts (id, name, email, api_key, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, q, a.ID, a.Name, a.Email, a.APIKey, a.Balance, a.CreatedAt, a.UpdatedAt)
	return err
}

func (r *PostgresAccountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	const q = `
		SELECT id, name, email, api_key, balance, created_at, updated_at
		FROM accounts WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, q, id)
	var a domain.Account
	if err := scanAccount(row, &a); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *PostgresAccountRepository) GetByAPIKey(ctx context.Context, apiKey string) (*domain.Account, error) {
	const q = `
		SELECT id, name, email, api_key, balance, created_at, updated_at
		FROM accounts WHERE api_key = $1
	`
	row := r.db.QueryRowContext(ctx, q, apiKey)
	var a domain.Account
	if err := scanAccount(row, &a); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *PostgresAccountRepository) UpdateBalance(ctx context.Context, id string, amount float64) error {
	// atomically increment balance and update updated_at
	const q = `
		UPDATE accounts
		SET balance = balance + $1, updated_at = $2
		WHERE id = $3
	`
	res, err := r.db.ExecContext(ctx, q, amount, time.Now().UTC(), id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrAccountNotFound
	}
	return nil
}

// scanAccount scans a single row into Account.
func scanAccount(row interface{ Scan(dest ...any) error }, a *domain.Account) error {
	return row.Scan(&a.ID, &a.Name, &a.Email, &a.APIKey, &a.Balance, &a.CreatedAt, &a.UpdatedAt)
}
