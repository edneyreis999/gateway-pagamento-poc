package memory

import (
	"context"
	"testing"
	"time"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
)

func TestInvoiceRepositoryMemory_Create(t *testing.T) {
	repo := NewInvoiceRepositoryMemory()
	ctx := context.Background()

	// Create a test invoice
	invoice := &domain.Invoice{
		ID:             "test-id",
		AccountID:      "test-account-id",
		Amount:         100.50,
		Status:         domain.StatusPending,
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	// Test creating invoice
	err := repo.Create(ctx, invoice)
	if err != nil {
		t.Errorf("failed to create invoice: %v", err)
	}

	// Test creating duplicate invoice
	err = repo.Create(ctx, invoice)
	if err == nil {
		t.Error("expected error when creating duplicate invoice")
	}
}

func TestInvoiceRepositoryMemory_GetByID(t *testing.T) {
	repo := NewInvoiceRepositoryMemory()
	ctx := context.Background()

	// Create a test invoice
	invoice := &domain.Invoice{
		ID:             "test-id",
		AccountID:      "test-account-id",
		Amount:         100.50,
		Status:         domain.StatusPending,
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	err := repo.Create(ctx, invoice)
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Test getting existing invoice
	retrieved, err := repo.GetByID(ctx, "test-id")
	if err != nil {
		t.Errorf("failed to get invoice by ID: %v", err)
	}

	if retrieved.ID != invoice.ID {
		t.Errorf("expected invoice ID %s, got %s", invoice.ID, retrieved.ID)
	}
	if retrieved.AccountID != invoice.AccountID {
		t.Errorf("expected account ID %s, got %s", invoice.AccountID, retrieved.AccountID)
	}
	if retrieved.Amount != invoice.Amount {
		t.Errorf("expected amount %f, got %f", invoice.Amount, retrieved.Amount)
	}
	if retrieved.Status != invoice.Status {
		t.Errorf("expected status %s, got %s", invoice.Status, retrieved.Status)
	}

	// Test getting non-existent invoice
	_, err = repo.GetByID(ctx, "non-existent-id")
	if err != domain.ErrInvoiceNotFound {
		t.Errorf("expected ErrInvoiceNotFound, got %v", err)
	}
}

func TestInvoiceRepositoryMemory_GetByAccountID(t *testing.T) {
	repo := NewInvoiceRepositoryMemory()
	ctx := context.Background()

	// Create test invoices for different accounts
	invoices := []*domain.Invoice{
		{
			ID:          "invoice-1",
			AccountID:   "account-1",
			Amount:      100.00,
			Status:      domain.StatusPending,
			Description: "Invoice 1",
			PaymentType: "credit_card",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          "invoice-2",
			AccountID:   "account-1",
			Amount:      200.00,
			Status:      domain.StatusApproved,
			Description: "Invoice 2",
			PaymentType: "debit_card",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          "invoice-3",
			AccountID:   "account-2",
			Amount:      150.00,
			Status:      domain.StatusPending,
			Description: "Invoice 3",
			PaymentType: "credit_card",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	}

	for _, invoice := range invoices {
		err := repo.Create(ctx, invoice)
		if err != nil {
			t.Fatalf("failed to create test invoice: %v", err)
		}
	}

	// Test getting invoices for account-1
	account1Invoices, err := repo.GetByAccountID(ctx, "account-1")
	if err != nil {
		t.Errorf("failed to get invoices by account ID: %v", err)
	}

	if len(account1Invoices) != 2 {
		t.Errorf("expected 2 invoices for account-1, got %d", len(account1Invoices))
	}

	// Test getting invoices for account-2
	account2Invoices, err := repo.GetByAccountID(ctx, "account-2")
	if err != nil {
		t.Errorf("failed to get invoices by account ID: %v", err)
	}

	if len(account2Invoices) != 1 {
		t.Errorf("expected 1 invoice for account-2, got %d", len(account2Invoices))
	}

	// Test getting invoices for non-existent account
	nonExistentInvoices, err := repo.GetByAccountID(ctx, "non-existent-account")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(nonExistentInvoices) != 0 {
		t.Errorf("expected 0 invoices for non-existent account, got %d", len(nonExistentInvoices))
	}
}

func TestInvoiceRepositoryMemory_UpdateStatus(t *testing.T) {
	repo := NewInvoiceRepositoryMemory()
	ctx := context.Background()

	// Create a test invoice
	invoice := &domain.Invoice{
		ID:             "test-id",
		AccountID:      "test-account-id",
		Amount:         100.50,
		Status:         domain.StatusPending,
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	err := repo.Create(ctx, invoice)
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Test updating status to approved
	err = repo.UpdateStatus(ctx, "test-id", domain.StatusApproved)
	if err != nil {
		t.Errorf("failed to update status to approved: %v", err)
	}

	// Verify the status was updated
	retrieved, err := repo.GetByID(ctx, "test-id")
	if err != nil {
		t.Errorf("failed to get updated invoice: %v", err)
	}

	if retrieved.Status != domain.StatusApproved {
		t.Errorf("expected status %s, got %s", domain.StatusApproved, retrieved.Status)
	}

	// Test updating status to rejected
	err = repo.UpdateStatus(ctx, "test-id", domain.StatusRejected)
	if err != nil {
		t.Errorf("failed to update status to rejected: %v", err)
	}

	// Verify the status was updated again
	retrieved, err = repo.GetByID(ctx, "test-id")
	if err != nil {
		t.Errorf("failed to get updated invoice: %v", err)
	}

	if retrieved.Status != domain.StatusRejected {
		t.Errorf("expected status %s, got %s", domain.StatusRejected, retrieved.Status)
	}

	// Test updating non-existent invoice
	err = repo.UpdateStatus(ctx, "non-existent-id", domain.StatusApproved)
	if err != domain.ErrInvoiceNotFound {
		t.Errorf("expected ErrInvoiceNotFound, got %v", err)
	}
}

func TestInvoiceRepositoryMemory_Concurrency(t *testing.T) {
	repo := NewInvoiceRepositoryMemory()
	ctx := context.Background()

	// Test concurrent access
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			invoice := &domain.Invoice{
				ID:          "concurrent-id-" + string(rune(id)),
				AccountID:   "account-id",
				Amount:      float64(id * 10),
				Status:      domain.StatusPending,
				Description: "Concurrent invoice",
				PaymentType: "credit_card",
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
			}

			err := repo.Create(ctx, invoice)
			if err != nil {
				t.Errorf("failed to create concurrent invoice: %v", err)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all invoices were created
	invoices, err := repo.GetByAccountID(ctx, "account-id")
	if err != nil {
		t.Errorf("failed to get concurrent invoices: %v", err)
	}

	if len(invoices) != 10 {
		t.Errorf("expected 10 concurrent invoices, got %d", len(invoices))
	}
}
