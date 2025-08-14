package service

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
)

func TestAccountService_CreateAndGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	svc := NewAccountService(db)
	ctx := context.Background()

	mock.ExpectExec("INSERT INTO accounts").
		WithArgs(sqlmock.AnyArg(), "Acme", "acme@example.com", sqlmock.AnyArg(), 0.0, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	out, err := svc.Create(ctx, AccountCreateInput{Name: "Acme", Email: "acme@example.com"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if out.Name != "Acme" {
		t.Fatalf("expected name Acme")
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "api_key", "balance", "created_at", "updated_at"}).
		AddRow(out.ID, out.Name, out.Email, out.APIKey, out.Balance, time.Now().UTC(), time.Now().UTC())
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, api_key, balance, created_at, updated_at FROM accounts WHERE id = $1")).
		WithArgs(out.ID).WillReturnRows(rows)

	got, err := svc.GetByID(ctx, out.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != out.ID {
		t.Fatalf("expected same id")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet: %v", err)
	}
}

func TestAccountService_UpdateBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	svc := NewAccountService(db)
	ctx := context.Background()

	// Mock GetByAPIKey call
	rows := sqlmock.NewRows([]string{"id", "name", "email", "api_key", "balance", "created_at", "updated_at"}).
		AddRow("acc-1", "Acme", "acme@example.com", "key-1", 100.0, time.Now().UTC(), time.Now().UTC())
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, api_key, balance, created_at, updated_at FROM accounts WHERE api_key = $1")).
		WithArgs("key-1").WillReturnRows(rows)

	// Mock UpdateBalance call (transaction)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM accounts WHERE id = $1 FOR UPDATE")).
		WithArgs("acc-1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("acc-1"))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE accounts SET balance = $1, updated_at = $2 WHERE id = $3")).
		WithArgs(110.0, sqlmock.AnyArg(), "acc-1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if err := svc.UpdateBalance(ctx, "key-1", 10); err != nil {
		t.Fatalf("update: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet: %v", err)
	}
}

func TestAccountService_GetByAPIKey_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	svc := NewAccountService(db)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, api_key, balance, created_at, updated_at FROM accounts WHERE api_key = $1")).
		WithArgs("nope").
		WillReturnError(sql.ErrNoRows)

	_, err = svc.GetByAPIKey(ctx, "nope")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err != domain.ErrAccountNotFound {
		t.Fatalf("expected not found")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet: %v", err)
	}
}
