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

func TestPostgresInvoiceRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresInvoiceRepository(db)
	ctx := context.Background()

	invoice := &domain.Invoice{
		ID:             "inv-1",
		AccountID:      "acc-1",
		Amount:         100.50,
		Status:         domain.StatusPending,
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO invoices (id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)")).
		WithArgs(invoice.ID, invoice.AccountID, invoice.Amount, invoice.Status, invoice.Description, invoice.PaymentType, invoice.CardLastDigits, invoice.CreatedAt, invoice.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.Create(ctx, invoice); err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet: %v", err)
	}
}

func TestPostgresInvoiceRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresInvoiceRepository(db)
	ctx := context.Background()

	invoice := &domain.Invoice{
		ID:             "inv-1",
		AccountID:      "acc-1",
		Amount:         100.50,
		Status:         domain.StatusPending,
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	rows := sqlmock.NewRows([]string{"id", "account_id", "amount", "status", "description", "payment_type", "card_last_digits", "created_at", "updated_at"}).
		AddRow(invoice.ID, invoice.AccountID, invoice.Amount, invoice.Status, invoice.Description, invoice.PaymentType, invoice.CardLastDigits, invoice.CreatedAt, invoice.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at FROM invoices WHERE id = $1")).
		WithArgs(invoice.ID).WillReturnRows(rows)

	got, err := repo.GetByID(ctx, invoice.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.ID != invoice.ID {
		t.Fatalf("expected same id")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet: %v", err)
	}
}

func TestPostgresInvoiceRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresInvoiceRepository(db)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at FROM invoices WHERE id = $1")).
		WithArgs("nope").WillReturnError(sql.ErrNoRows)

	_, err = repo.GetByID(ctx, "nope")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err != domain.ErrInvoiceNotFound {
		t.Fatalf("expected not found, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet: %v", err)
	}
}

func TestPostgresInvoiceRepository_GetByAccountID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresInvoiceRepository(db)
	ctx := context.Background()

	accountID := "acc-1"
	invoices := []*domain.Invoice{
		{
			ID:             "inv-1",
			AccountID:      accountID,
			Amount:         100.00,
			Status:         domain.StatusPending,
			Description:    "Invoice 1",
			PaymentType:    "credit_card",
			CardLastDigits: "1234",
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		},
		{
			ID:             "inv-2",
			AccountID:      accountID,
			Amount:         200.00,
			Status:         domain.StatusApproved,
			Description:    "Invoice 2",
			PaymentType:    "debit_card",
			CardLastDigits: "5678",
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "account_id", "amount", "status", "description", "payment_type", "card_last_digits", "created_at", "updated_at"})
	for _, invoice := range invoices {
		rows.AddRow(invoice.ID, invoice.AccountID, invoice.Amount, invoice.Status, invoice.Description, invoice.PaymentType, invoice.CardLastDigits, invoice.CreatedAt, invoice.UpdatedAt)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at FROM invoices WHERE account_id = $1 ORDER BY created_at DESC")).
		WithArgs(accountID).WillReturnRows(rows)

	got, err := repo.GetByAccountID(ctx, accountID)
	if err != nil {
		t.Fatalf("get by account id: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 invoices, got %d", len(got))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet: %v", err)
	}
}

func TestPostgresInvoiceRepository_GetByAccountID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresInvoiceRepository(db)
	ctx := context.Background()

	accountID := "acc-2"

	rows := sqlmock.NewRows([]string{"id", "account_id", "amount", "status", "description", "payment_type", "card_last_digits", "created_at", "updated_at"})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at FROM invoices WHERE account_id = $1 ORDER BY created_at DESC")).
		WithArgs(accountID).WillReturnRows(rows)

	got, err := repo.GetByAccountID(ctx, accountID)
	if err != nil {
		t.Fatalf("get by account id: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 invoices, got %d", len(got))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet: %v", err)
	}
}

func TestPostgresInvoiceRepository_UpdateStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresInvoiceRepository(db)
	ctx := context.Background()

	invoiceID := "inv-1"
	newStatus := domain.StatusApproved

	mock.ExpectExec(regexp.QuoteMeta("UPDATE invoices SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2")).
		WithArgs(newStatus, invoiceID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.UpdateStatus(ctx, invoiceID, newStatus); err != nil {
		t.Fatalf("update status: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet: %v", err)
	}
}

func TestPostgresInvoiceRepository_UpdateStatus_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresInvoiceRepository(db)
	ctx := context.Background()

	invoiceID := "missing"
	newStatus := domain.StatusApproved

	mock.ExpectExec(regexp.QuoteMeta("UPDATE invoices SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2")).
		WithArgs(newStatus, invoiceID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	updateErr := repo.UpdateStatus(ctx, invoiceID, newStatus)
	if updateErr == nil {
		t.Fatalf("expected error")
	}
	if updateErr != domain.ErrInvoiceNotFound {
		t.Fatalf("expected not found, got %v", updateErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet: %v", err)
	}
}
