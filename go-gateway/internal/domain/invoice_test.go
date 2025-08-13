package domain

import (
	"testing"
	"time"
)

func TestNewInvoice(t *testing.T) {
	tests := []struct {
		name           string
		accountID      string
		description    string
		paymentType    string
		amount         float64
		cardLastDigits string
		expectedError  bool
	}{
		{
			name:           "valid invoice",
			accountID:      "test-account-id",
			description:    "Test invoice description",
			paymentType:    "credit_card",
			amount:         100.50,
			cardLastDigits: "1234",
			expectedError:  false,
		},
		{
			name:           "empty account ID",
			accountID:      "",
			description:    "Test invoice description",
			paymentType:    "credit_card",
			amount:         100.50,
			cardLastDigits: "1234",
			expectedError:  true,
		},
		{
			name:           "description too short",
			accountID:      "test-account-id",
			description:    "Te",
			paymentType:    "credit_card",
			amount:         100.50,
			cardLastDigits: "1234",
			expectedError:  true,
		},
		{
			name:           "empty payment type",
			accountID:      "test-account-id",
			description:    "Test invoice description",
			paymentType:    "",
			amount:         100.50,
			cardLastDigits: "1234",
			expectedError:  true,
		},
		{
			name:           "negative amount",
			accountID:      "test-account-id",
			description:    "Test invoice description",
			paymentType:    "credit_card",
			amount:         -50.00,
			cardLastDigits: "1234",
			expectedError:  true,
		},
		{
			name:           "zero amount",
			accountID:      "test-account-id",
			description:    "Test invoice description",
			paymentType:    "credit_card",
			amount:         0.00,
			cardLastDigits: "1234",
			expectedError:  true,
		},
		{
			name:           "valid invoice without card digits",
			accountID:      "test-account-id",
			description:    "Test invoice description",
			paymentType:    "pix",
			amount:         100.50,
			cardLastDigits: "",
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invoice, err := NewInvoice(tt.accountID, tt.description, tt.paymentType, tt.amount, tt.cardLastDigits)

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

			if invoice.ID == "" {
				t.Error("expected invoice ID to be set")
			}
			if invoice.AccountID != tt.accountID {
				t.Errorf("expected account ID %s, got %s", tt.accountID, invoice.AccountID)
			}
			if invoice.Amount != tt.amount {
				t.Errorf("expected amount %f, got %f", tt.amount, invoice.Amount)
			}
			if invoice.Status != StatusPending {
				t.Errorf("expected status %s, got %s", StatusPending, invoice.Status)
			}
			if invoice.Description != tt.description {
				t.Errorf("expected description %s, got %s", tt.description, invoice.Description)
			}
			if invoice.PaymentType != tt.paymentType {
				t.Errorf("expected payment type %s, got %s", tt.paymentType, invoice.PaymentType)
			}
			if invoice.CardLastDigits != tt.cardLastDigits {
				t.Errorf("expected card last digits %s, got %s", tt.cardLastDigits, invoice.CardLastDigits)
			}
			if invoice.CreatedAt.IsZero() {
				t.Error("expected created at to be set")
			}
			if invoice.UpdatedAt.IsZero() {
				t.Error("expected updated at to be set")
			}
		})
	}
}

func TestInvoice_Process(t *testing.T) {
	invoice, err := NewInvoice("test-account-id", "Test invoice", "credit_card", 100.50, "1234")
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Test processing pending invoice
	err = invoice.Process()
	if err != nil {
		t.Errorf("failed to process pending invoice: %v", err)
	}

	// Change status to approved to test processing non-pending invoice
	err = invoice.UpdateStatus(StatusApproved)
	if err != nil {
		t.Fatalf("failed to update status to approved: %v", err)
	}

	// Test processing non-pending invoice
	err = invoice.Process()
	if err == nil {
		t.Error("expected error when processing non-pending invoice")
	}
}

func TestInvoice_UpdateStatus(t *testing.T) {
	invoice, err := NewInvoice("test-account-id", "Test invoice", "credit_card", 100.50, "1234")
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Test updating to approved status
	err = invoice.UpdateStatus(StatusApproved)
	if err != nil {
		t.Errorf("failed to update status to approved: %v", err)
	}
	if invoice.Status != StatusApproved {
		t.Errorf("expected status %s, got %s", StatusApproved, invoice.Status)
	}

	// Test updating to rejected status
	err = invoice.UpdateStatus(StatusRejected)
	if err != nil {
		t.Errorf("failed to update status to rejected: %v", err)
	}
	if invoice.Status != StatusRejected {
		t.Errorf("expected status %s, got %s", StatusRejected, invoice.Status)
	}

	// Test updating to pending status
	err = invoice.UpdateStatus(StatusPending)
	if err != nil {
		t.Errorf("failed to update status to pending: %v", err)
	}
	if invoice.Status != StatusPending {
		t.Errorf("expected status %s, got %s", StatusPending, invoice.Status)
	}

	// Test updating to invalid status
	err = invoice.UpdateStatus("invalid_status")
	if err != ErrInvalidStatus {
		t.Errorf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestInvoice_StatusChecks(t *testing.T) {
	invoice, err := NewInvoice("test-account-id", "Test invoice", "credit_card", 100.50, "1234")
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Test initial status
	if !invoice.IsPending() {
		t.Error("expected invoice to be pending initially")
	}
	if invoice.IsApproved() {
		t.Error("expected invoice not to be approved initially")
	}
	if invoice.IsRejected() {
		t.Error("expected invoice not to be rejected initially")
	}

	// Test approved status
	err = invoice.UpdateStatus(StatusApproved)
	if err != nil {
		t.Fatalf("failed to update status to approved: %v", err)
	}

	if invoice.IsPending() {
		t.Error("expected invoice not to be pending after approval")
	}
	if !invoice.IsApproved() {
		t.Error("expected invoice to be approved")
	}
	if invoice.IsRejected() {
		t.Error("expected invoice not to be rejected")
	}

	// Test rejected status
	err = invoice.UpdateStatus(StatusRejected)
	if err != nil {
		t.Fatalf("failed to update status to rejected: %v", err)
	}

	if invoice.IsPending() {
		t.Error("expected invoice not to be pending after rejection")
	}
	if invoice.IsApproved() {
		t.Error("expected invoice not to be approved after rejection")
	}
	if !invoice.IsRejected() {
		t.Error("expected invoice to be rejected")
	}
}

func TestInvoice_Timestamps(t *testing.T) {
	beforeCreation := time.Now().UTC()

	invoice, err := NewInvoice("test-account-id", "Test invoice", "credit_card", 100.50, "1234")
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	afterCreation := time.Now().UTC()

	// Check that timestamps are set correctly
	if invoice.CreatedAt.Before(beforeCreation) || invoice.CreatedAt.After(afterCreation) {
		t.Errorf("created at timestamp %v is not within expected range [%v, %v]",
			invoice.CreatedAt, beforeCreation, afterCreation)
	}

	if invoice.UpdatedAt.Before(beforeCreation) || invoice.UpdatedAt.After(afterCreation) {
		t.Errorf("updated at timestamp %v is not within expected range [%v, %v]",
			invoice.UpdatedAt, beforeCreation, afterCreation)
	}

	// Check that timestamps are updated when status changes
	originalUpdatedAt := invoice.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	err = invoice.UpdateStatus(StatusApproved)
	if err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	if !invoice.UpdatedAt.After(originalUpdatedAt) {
		t.Error("expected updated at timestamp to be updated after status change")
	}
}

func TestInvoice_Concurrency(t *testing.T) {
	invoice, err := NewInvoice("test-account-id", "Test invoice", "credit_card", 100.50, "1234")
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Test concurrent status updates
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			err := invoice.UpdateStatus(StatusApproved)
			if err != nil {
				t.Errorf("failed to update status concurrently: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify final status
	if invoice.Status != StatusApproved {
		t.Errorf("expected final status %s, got %s", StatusApproved, invoice.Status)
	}
}
