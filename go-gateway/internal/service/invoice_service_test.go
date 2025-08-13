package service

import (
	"context"
	"testing"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	"github.com/devfullcycle/imersao22/go-gateway/internal/repository/memory"
)

func TestInvoiceService_Create(t *testing.T) {
	repo := memory.NewInvoiceRepositoryMemory()
	svc := &InvoiceService{repo: repo}

	tests := []struct {
		name          string
		input         InvoiceCreateInput
		expectedError bool
	}{
		{
			name: "valid invoice",
			input: InvoiceCreateInput{
				AccountID:      "test-account-id",
				Amount:         100.50,
				Description:    "Test invoice",
				PaymentType:    "credit_card",
				CardLastDigits: "1234",
			},
			expectedError: false,
		},
		{
			name: "invalid description too short",
			input: InvoiceCreateInput{
				AccountID:   "test-account-id",
				Amount:      100.50,
				Description: "Te",
				PaymentType: "credit_card",
			},
			expectedError: true,
		},
		{
			name: "invalid payment type empty",
			input: InvoiceCreateInput{
				AccountID:   "test-account-id",
				Amount:      100.50,
				Description: "Test invoice",
				PaymentType: "",
			},
			expectedError: true,
		},
		{
			name: "invalid amount negative",
			input: InvoiceCreateInput{
				AccountID:   "test-account-id",
				Amount:      -50.00,
				Description: "Test invoice",
				PaymentType: "credit_card",
			},
			expectedError: true,
		},
		{
			name: "invalid amount zero",
			input: InvoiceCreateInput{
				AccountID:   "test-account-id",
				Amount:      0.00,
				Description: "Test invoice",
				PaymentType: "credit_card",
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := svc.Create(context.Background(), tt.input)

			if tt.expectedError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if output.ID == "" {
				t.Error("expected invoice ID to be set")
			}
			if output.AccountID != tt.input.AccountID {
				t.Errorf("expected account ID %s, got %s", tt.input.AccountID, output.AccountID)
			}
			if output.Amount != tt.input.Amount {
				t.Errorf("expected amount %f, got %f", tt.input.Amount, output.Amount)
			}
			if output.Status != "pending" {
				t.Errorf("expected status 'pending', got '%s'", output.Status)
			}
			if output.Description != tt.input.Description {
				t.Errorf("expected description %s, got %s", tt.input.Description, output.Description)
			}
			if output.PaymentType != tt.input.PaymentType {
				t.Errorf("expected payment type %s, got %s", tt.input.PaymentType, output.PaymentType)
			}
			if output.CardLastDigits != tt.input.CardLastDigits {
				t.Errorf("expected card last digits %s, got %s", tt.input.CardLastDigits, output.CardLastDigits)
			}
		})
	}
}

func TestInvoiceService_GetByID(t *testing.T) {
	repo := memory.NewInvoiceRepositoryMemory()
	svc := &InvoiceService{repo: repo}

	// Create a test invoice first
	input := InvoiceCreateInput{
		AccountID:      "test-account-id",
		Amount:         100.50,
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
	}

	created, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Test getting by ID
	retrieved, err := svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Errorf("failed to get invoice by ID: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("expected invoice ID %s, got %s", created.ID, retrieved.ID)
	}

	// Test getting non-existent invoice
	_, err = svc.GetByID(context.Background(), "non-existent-id")
	if err != domain.ErrInvoiceNotFound {
		t.Errorf("expected ErrInvoiceNotFound, got %v", err)
	}
}

func TestInvoiceService_GetByAccountID(t *testing.T) {
	repo := memory.NewInvoiceRepositoryMemory()
	svc := &InvoiceService{repo: repo}

	// Create test invoices for different accounts
	inputs := []InvoiceCreateInput{
		{
			AccountID:   "account-1",
			Amount:      100.00,
			Description: "Invoice 1",
			PaymentType: "credit_card",
		},
		{
			AccountID:   "account-1",
			Amount:      200.00,
			Description: "Invoice 2",
			PaymentType: "debit_card",
		},
		{
			AccountID:   "account-2",
			Amount:      150.00,
			Description: "Invoice 3",
			PaymentType: "credit_card",
		},
	}

	for _, input := range inputs {
		_, err := svc.Create(context.Background(), input)
		if err != nil {
			t.Fatalf("failed to create test invoice: %v", err)
		}
	}

	// Test getting invoices for account-1
	invoices, err := svc.GetByAccountID(context.Background(), "account-1")
	if err != nil {
		t.Errorf("failed to get invoices by account ID: %v", err)
	}

	if len(invoices) != 2 {
		t.Errorf("expected 2 invoices for account-1, got %d", len(invoices))
	}

	// Test getting invoices for account-2
	invoices, err = svc.GetByAccountID(context.Background(), "account-2")
	if err != nil {
		t.Errorf("failed to get invoices by account ID: %v", err)
	}

	if len(invoices) != 1 {
		t.Errorf("expected 1 invoice for account-2, got %d", len(invoices))
	}

	// Test getting invoices for non-existent account
	invoices, err = svc.GetByAccountID(context.Background(), "non-existent-account")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(invoices) != 0 {
		t.Errorf("expected 0 invoices for non-existent account, got %d", len(invoices))
	}
}

func TestInvoiceService_UpdateStatus(t *testing.T) {
	repo := memory.NewInvoiceRepositoryMemory()
	svc := &InvoiceService{repo: repo}

	// Create a test invoice first
	input := InvoiceCreateInput{
		AccountID:      "test-account-id",
		Amount:         100.50,
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
	}

	created, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Test updating status to approved
	err = svc.UpdateStatus(context.Background(), created.ID, domain.StatusApproved)
	if err != nil {
		t.Errorf("failed to update status to approved: %v", err)
	}

	// Verify the status was updated
	retrieved, err := svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Errorf("failed to get updated invoice: %v", err)
	}

	if retrieved.Status != "approved" {
		t.Errorf("expected status 'approved', got '%s'", retrieved.Status)
	}

	// Test updating status to rejected
	err = svc.UpdateStatus(context.Background(), created.ID, domain.StatusRejected)
	if err != nil {
		t.Errorf("failed to update status to rejected: %v", err)
	}

	// Verify the status was updated again
	retrieved, err = svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Errorf("failed to get updated invoice: %v", err)
	}

	if retrieved.Status != "rejected" {
		t.Errorf("expected status 'rejected', got '%s'", retrieved.Status)
	}

	// Test updating non-existent invoice
	err = svc.UpdateStatus(context.Background(), "non-existent-id", domain.StatusApproved)
	if err != domain.ErrInvoiceNotFound {
		t.Errorf("expected ErrInvoiceNotFound, got %v", err)
	}
}
