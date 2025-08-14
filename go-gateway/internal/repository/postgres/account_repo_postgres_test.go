package postgres

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
)

func TestPostgresAccountRepository_CreateAndGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresAccountRepository(db)
	ctx := context.Background()

	a := &domain.Account{
		ID:        "acc-1",
		Name:      "Acme",
		Email:     "acme@example.com",
		APIKey:    "key-1",
		Balance:   0,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO accounts (id, name, email, api_key, balance, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
		WithArgs(a.ID, a.Name, a.Email, a.APIKey, a.Balance, a.CreatedAt, a.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.Create(ctx, a); err != nil {
		t.Fatalf("create: %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "api_key", "balance", "created_at", "updated_at"}).
		AddRow(a.ID, a.Name, a.Email, a.APIKey, a.Balance, a.CreatedAt, a.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, api_key, balance, created_at, updated_at FROM accounts WHERE id = $1")).
		WithArgs(a.ID).WillReturnRows(rows)

	got, err := repo.GetByID(ctx, a.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.ID != a.ID {
		t.Fatalf("expected same id")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet: %v", err)
	}
}

func TestPostgresAccountRepository_GetByAPIKey_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresAccountRepository(db)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, api_key, balance, created_at, updated_at FROM accounts WHERE api_key = $1")).
		WithArgs("nope").WillReturnError(sql.ErrNoRows)

	_, err = repo.GetByAPIKey(ctx, "nope")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err != domain.ErrAccountNotFound {
		t.Fatalf("expected not found, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet: %v", err)
	}
}

func TestPostgresAccountRepository_UpdateBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresAccountRepository(db)
	ctx := context.Background()

	account := &domain.Account{
		ID:        "acc-1",
		Name:      "Acme",
		Email:     "acme@example.com",
		APIKey:    "key-1",
		Balance:   150.0, // Updated balance
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM accounts WHERE id = $1 FOR UPDATE")).
		WithArgs("acc-1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("acc-1"))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE accounts SET balance = $1, updated_at = $2 WHERE id = $3")).
		WithArgs(150.0, account.UpdatedAt, "acc-1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if err := repo.UpdateBalance(ctx, account); err != nil {
		t.Fatalf("update balance: %v", err)
	}

	// Not found: lock query returns no rows
	missingAccount := &domain.Account{
		ID:        "missing",
		Name:      "Missing",
		Email:     "missing@example.com",
		APIKey:    "key-missing",
		Balance:   50.0,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM accounts WHERE id = $1 FOR UPDATE")).
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	err = repo.UpdateBalance(ctx, missingAccount)
	if err != domain.ErrAccountNotFound {
		t.Fatalf("expected not found, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet: %v", err)
	}
}
