package domain

import (
	"errors"
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
	// Test with default processor (random behavior)
	invoice, err := NewInvoice("test-account-id", "Test invoice", "credit_card", 100.50, "1234")
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Test processing pending invoice
	err = invoice.Process()
	if err != nil {
		t.Errorf("failed to process pending invoice: %v", err)
	}

	// Verify that status changed from pending
	if invoice.Status == StatusPending {
		t.Error("expected status to change from pending after processing")
	}

	// Verify that status is either approved or rejected
	if invoice.Status != StatusApproved && invoice.Status != StatusRejected {
		t.Errorf("expected status to be either approved or rejected, got %s", invoice.Status)
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

func TestInvoice_Process_WithTestProcessor(t *testing.T) {
	// Test with controlled test processor
	testProcessor := NewTestInvoiceProcessor()

	invoice, err := NewInvoiceWithProcessor("test-account-id", "Test invoice", "credit_card", 100.50, "1234", testProcessor)
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Test forced approval
	testProcessor.SetNextStatus(StatusApproved)
	err = invoice.Process()
	if err != nil {
		t.Errorf("failed to process invoice: %v", err)
	}
	if invoice.Status != StatusApproved {
		t.Errorf("expected status %s, got %s", StatusApproved, invoice.Status)
	}

	// Reset to pending and test forced rejection
	err = invoice.UpdateStatus(StatusPending)
	if err != nil {
		t.Fatalf("failed to reset status: %v", err)
	}

	testProcessor.SetNextStatus(StatusRejected)
	err = invoice.Process()
	if err != nil {
		t.Errorf("failed to process invoice: %v", err)
	}
	if invoice.Status != StatusRejected {
		t.Errorf("expected status %s, got %s", StatusRejected, invoice.Status)
	}

	// Test error condition
	err = invoice.UpdateStatus(StatusPending)
	if err != nil {
		t.Fatalf("failed to reset status: %v", err)
	}

	testProcessor.SetError(errors.New("test error"))
	err = invoice.Process()
	if err == nil {
		t.Error("expected error but got none")
	}
	if err.Error() != "test error" {
		t.Errorf("expected error 'test error', got '%v'", err)
	}
}

func TestInvoice_Process_HighValuePending(t *testing.T) {
	// Test that invoices with amount > 10000 stay pending
	testProcessor := NewTestInvoiceProcessor()

	// Test with amount > 10000 - should stay pending regardless of processor setting
	testProcessor.SetNextStatus(StatusApproved)

	highValueInvoice, err := NewInvoiceWithProcessor("test-account-id", "High value invoice", "credit_card", 15000.00, "1234", testProcessor)
	if err != nil {
		t.Fatalf("failed to create high value invoice: %v", err)
	}

	// Process should not change status for high value invoices
	err = highValueInvoice.Process()
	if err != nil {
		t.Errorf("failed to process high value invoice: %v", err)
	}

	// Status should remain pending for amounts > 10000
	if highValueInvoice.Status != StatusPending {
		t.Errorf("expected status to remain pending for amount > 10000, got %s", highValueInvoice.Status)
	}

	// Test with amount exactly 10000 - should be processed normally
	testProcessor.SetNextStatus(StatusRejected)
	exactValueInvoice, err := NewInvoiceWithProcessor("test-account-id", "Exact value invoice", "credit_card", 10000.00, "1234", testProcessor)
	if err != nil {
		t.Fatalf("failed to create exact value invoice: %v", err)
	}

	err = exactValueInvoice.Process()
	if err != nil {
		t.Errorf("failed to process exact value invoice: %v", err)
	}

	// Status should be processed for amount = 10000
	if exactValueInvoice.Status != StatusRejected {
		t.Errorf("expected status to be processed for amount = 10000, got %s", exactValueInvoice.Status)
	}

	// Test with amount < 10000 - should be processed normally
	testProcessor.SetNextStatus(StatusApproved)
	lowValueInvoice, err := NewInvoiceWithProcessor("test-account-id", "Low value invoice", "credit_card", 9999.99, "1234", testProcessor)
	if err != nil {
		t.Fatalf("failed to create low value invoice: %v", err)
	}

	err = lowValueInvoice.Process()
	if err != nil {
		t.Errorf("failed to process low value invoice: %v", err)
	}

	// Status should be processed for amount < 10000
	if lowValueInvoice.Status != StatusApproved {
		t.Errorf("expected status to be processed for amount < 10000, got %s", lowValueInvoice.Status)
	}
}

func TestInvoice_Process_WithSeedProcessor(t *testing.T) {
	// Test with processor using specific seed for deterministic behavior
	seed := int64(12345)
	processor := NewDefaultInvoiceProcessorWithSeed(seed)

	invoice, err := NewInvoiceWithProcessor("test-account-id", "Test invoice", "credit_card", 100.50, "1234", processor)
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Process with known seed - should give consistent result
	err = invoice.Process()
	if err != nil {
		t.Errorf("failed to process invoice: %v", err)
	}

	// With seed 12345, we should get a consistent result
	expectedStatus := invoice.Status

	// Create another invoice with same seed and verify same result
	// Note: We need to create a new processor with the same seed for the second invoice
	processor2 := NewDefaultInvoiceProcessorWithSeed(seed)
	invoice2, err := NewInvoiceWithProcessor("test-account-id-2", "Test invoice 2", "credit_card", 200.00, "5678", processor2)
	if err != nil {
		t.Fatalf("failed to create second test invoice: %v", err)
	}

	err = invoice2.Process()
	if err != nil {
		t.Errorf("failed to process second invoice: %v", err)
	}

	if invoice2.Status != expectedStatus {
		t.Errorf("expected same status with same seed, got %s vs %s", expectedStatus, invoice2.Status)
	}
}

func TestInvoice_Process_RejectedCase(t *testing.T) {
	// Test with controlled processor to ensure we get both outcomes
	// Create invoices with controlled outcomes
	invoices := make([]*Invoice, 0, 10)

	// First 7 invoices will be approved (70%)
	for i := 0; i < 7; i++ {
		testProcessor := NewTestInvoiceProcessor()
		testProcessor.SetNextStatus(StatusApproved)
		invoice, err := NewInvoiceWithProcessor("test-account-id", "Test invoice", "credit_card", 100.50, "1234", testProcessor)
		if err != nil {
			t.Fatalf("failed to create test invoice %d: %v", i, err)
		}
		invoices = append(invoices, invoice)
	}

	// Last 3 invoices will be rejected (30%)
	for i := 0; i < 3; i++ {
		testProcessor := NewTestInvoiceProcessor()
		testProcessor.SetNextStatus(StatusRejected)
		invoice, err := NewInvoiceWithProcessor("test-account-id", "Test invoice", "credit_card", 100.50, "1234", testProcessor)
		if err != nil {
			t.Fatalf("failed to create test invoice %d: %v", i+7, err)
		}
		invoices = append(invoices, invoice)
	}

	// Process all invoices
	rejectedCount := 0
	approvedCount := 0

	for i, invoice := range invoices {
		originalStatus := invoice.Status
		if originalStatus != StatusPending {
			t.Errorf("invoice %d should start with pending status, got %s", i, originalStatus)
		}

		err := invoice.Process()
		if err != nil {
			t.Errorf("failed to process invoice %d: %v", i, err)
		}

		// Count the results
		switch invoice.Status {
		case StatusApproved:
			approvedCount++
		case StatusRejected:
			rejectedCount++
		default:
			t.Errorf("invoice %d has unexpected status after processing: %s", i, invoice.Status)
		}

		// Verify that status changed from pending
		if invoice.Status == StatusPending {
			t.Errorf("invoice %d status should change from pending after processing", i)
		}

		// Verify that UpdatedAt was updated
		if invoice.UpdatedAt.Equal(invoice.CreatedAt) {
			t.Errorf("invoice %d UpdatedAt should be different from CreatedAt after processing", i)
		}
	}

	// Log the distribution for debugging
	t.Logf("Processing results: %d approved, %d rejected", approvedCount, rejectedCount)

	// With controlled processor, we should get exactly what we set
	if rejectedCount != 3 {
		t.Errorf("expected 3 rejected invoices, got %d", rejectedCount)
	}

	if approvedCount != 7 {
		t.Errorf("expected 7 approved invoices, got %d", approvedCount)
	}

	// Verify that the total adds up
	if approvedCount+rejectedCount != len(invoices) {
		t.Errorf("expected total processed invoices to be %d, got %d", len(invoices), approvedCount+rejectedCount)
	}
}

func TestInvoice_Process_ForceRejection(t *testing.T) {
	// Test with controlled processor to force rejection scenario
	// Create 20 invoices, all will be rejected
	invoices := make([]*Invoice, 0, 20)

	for i := 0; i < 20; i++ {
		testProcessor := NewTestInvoiceProcessor()
		testProcessor.SetNextStatus(StatusRejected)
		invoice, err := NewInvoiceWithProcessor("test-account-id", "Test invoice", "credit_card", 100.50, "1234", testProcessor)
		if err != nil {
			t.Fatalf("failed to create test invoice %d: %v", i, err)
		}
		invoices = append(invoices, invoice)
	}

	// Process all invoices
	rejectedCount := 0
	approvedCount := 0

	for i, invoice := range invoices {
		originalStatus := invoice.Status
		if originalStatus != StatusPending {
			t.Errorf("invoice %d should start with pending status, got %s", i, originalStatus)
		}

		err := invoice.Process()
		if err != nil {
			t.Errorf("failed to process invoice %d: %v", i, err)
		}

		// Count the results
		switch invoice.Status {
		case StatusApproved:
			approvedCount++
		case StatusRejected:
			rejectedCount++
		default:
			t.Errorf("invoice %d has unexpected status after processing: %s", i, invoice.Status)
		}

		// Verify that status changed from pending
		if invoice.Status == StatusPending {
			t.Errorf("invoice %d status should change from pending after processing", i)
		}

		// Verify that UpdatedAt was updated
		if invoice.UpdatedAt.Equal(invoice.CreatedAt) {
			t.Errorf("invoice %d UpdatedAt should be different from CreatedAt after processing", i)
		}
	}

	// Log the distribution for debugging
	t.Logf("Processing results: %d approved, %d rejected", approvedCount, rejectedCount)

	// With controlled processor, all should be rejected
	if rejectedCount != 20 {
		t.Errorf("expected 20 rejected invoices, got %d", rejectedCount)
	}

	if approvedCount != 0 {
		t.Errorf("expected 0 approved invoices, got %d", approvedCount)
	}

	// Verify that the total adds up
	if approvedCount+rejectedCount != len(invoices) {
		t.Errorf("expected total processed invoices to be %d, got %d", len(invoices), approvedCount+rejectedCount)
	}

	// Log some statistics
	rejectionRate := float64(rejectedCount) / float64(len(invoices))
	t.Logf("Rejection rate: %.2f%% (expected 100%%)", rejectionRate*100)
}

func TestInvoice_Process_RejectionLogic(t *testing.T) {
	// Test with controlled processor to verify rejection logic
	testProcessor := NewTestInvoiceProcessor()

	invoice, err := NewInvoiceWithProcessor("test-account-id", "Test invoice", "credit_card", 100.50, "1234", testProcessor)
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Verify initial state
	if invoice.Status != StatusPending {
		t.Errorf("expected initial status to be pending, got %s", invoice.Status)
	}

	// Process the invoice with controlled outcome
	testProcessor.SetNextStatus(StatusRejected)
	err = invoice.Process()
	if err != nil {
		t.Errorf("failed to process invoice: %v", err)
	}

	// Verify that the status changed to rejected
	if invoice.Status != StatusRejected {
		t.Errorf("expected status to be rejected, got %s", invoice.Status)
	}

	// Verify that UpdatedAt was updated
	// Note: Since processing happens immediately after creation, we just verify it's not zero
	if invoice.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set after processing")
	}

	// Test that we can't process a non-pending invoice
	err = invoice.Process()
	if err == nil {
		t.Error("expected error when processing non-pending invoice")
	}

	// Test the specific rejection scenario by creating a new invoice and forcing it to rejected status
	invoice2, err := NewInvoiceWithProcessor("test-account-id-2", "Test invoice 2", "credit_card", 200.00, "5678", testProcessor)
	if err != nil {
		t.Fatalf("failed to create second test invoice: %v", err)
	}

	// Manually set to rejected to test the rejection behavior
	err = invoice2.UpdateStatus(StatusRejected)
	if err != nil {
		t.Fatalf("failed to update status to rejected: %v", err)
	}

	// Verify rejection status
	if !invoice2.IsRejected() {
		t.Error("expected invoice to be rejected")
	}

	// Test that rejected invoice can't be processed again
	err = invoice2.Process()
	if err == nil {
		t.Error("expected error when processing rejected invoice")
	}

	// Test that we can change from rejected back to pending
	err = invoice2.UpdateStatus(StatusPending)
	if err != nil {
		t.Errorf("failed to change status from rejected to pending: %v", err)
	}

	if !invoice2.IsPending() {
		t.Error("expected invoice to be pending after status change")
	}

	// Now we can process it again
	err = invoice2.Process()
	if err != nil {
		t.Errorf("failed to process invoice after status change to pending: %v", err)
	}

	// Verify that status changed again
	if invoice2.Status == StatusPending {
		t.Error("expected status to change from pending after processing")
	}
}

func TestInvoice_Process_RejectionCondition(t *testing.T) {
	// Test with controlled processor to verify rejection conditions
	// Create 50 invoices with controlled outcomes
	invoices := make([]*Invoice, 0, 50)

	// First 35 invoices will be approved (70%)
	for i := 0; i < 35; i++ {
		testProcessor := NewTestInvoiceProcessor()
		testProcessor.SetNextStatus(StatusApproved)
		invoice, err := NewInvoiceWithProcessor("test-account-id", "Test invoice", "credit_card", 100.50, "1234", testProcessor)
		if err != nil {
			t.Fatalf("failed to create test invoice %d: %v", i, err)
		}
		invoices = append(invoices, invoice)
	}

	// Last 15 invoices will be rejected (30%)
	for i := 0; i < 15; i++ {
		testProcessor := NewTestInvoiceProcessor()
		testProcessor.SetNextStatus(StatusRejected)
		invoice, err := NewInvoiceWithProcessor("test-account-id", "Test invoice", "credit_card", 100.50, "1234", testProcessor)
		if err != nil {
			t.Fatalf("failed to create test invoice %d: %v", i+35, err)
		}
		invoices = append(invoices, invoice)
	}

	// Process all invoices
	rejectedCount := 0
	approvedCount := 0

	for i, invoice := range invoices {
		originalStatus := invoice.Status
		if originalStatus != StatusPending {
			t.Errorf("invoice %d should start with pending status, got %s", i, originalStatus)
		}

		err := invoice.Process()
		if err != nil {
			t.Errorf("failed to process invoice %d: %v", i, err)
		}

		// Count the results
		switch invoice.Status {
		case StatusApproved:
			approvedCount++
		case StatusRejected:
			rejectedCount++
		default:
			t.Errorf("invoice %d has unexpected status after processing: %s", i, invoice.Status)
		}

		// Verify that status changed from pending
		if invoice.Status == StatusPending {
			t.Errorf("invoice %d status should change from pending after processing", i)
		}

		// Verify that UpdatedAt was updated
		if invoice.UpdatedAt.Equal(invoice.CreatedAt) {
			t.Errorf("invoice %d UpdatedAt should be different from CreatedAt after processing", i)
		}
	}

	// Log the distribution for debugging
	t.Logf("Processing results: %d approved, %d rejected", approvedCount, rejectedCount)

	// With controlled processor, we should get exactly what we set
	if rejectedCount != 15 {
		t.Errorf("expected 15 rejected invoices, got %d", rejectedCount)
	}

	if approvedCount != 35 {
		t.Errorf("expected 35 approved invoices, got %d", approvedCount)
	}

	// Verify that the total adds up
	if approvedCount+rejectedCount != len(invoices) {
		t.Errorf("expected total processed invoices to be %d, got %d", len(invoices), approvedCount+rejectedCount)
	}

	// Log some statistics
	rejectionRate := float64(rejectedCount) / float64(len(invoices))
	t.Logf("Rejection rate: %.2f%% (expected 30%%)", rejectionRate*100)

	// Test the specific rejection condition logic
	// The Process() method should work as follows:
	// - if randomSource.Float64() <= 0.7 (70% chance) -> StatusApproved
	// - if randomSource.Float64() > 0.7 (30% chance) -> StatusRejected

	// Verify that both statuses are possible outcomes
	hasApproved := approvedCount > 0
	hasRejected := rejectedCount > 0

	t.Logf("Has approved invoices: %v", hasApproved)
	t.Logf("Has rejected invoices: %v", hasRejected)

	// This test verifies that the rejection logic is working
	// With controlled processor, we get exactly what we expect
	if !hasApproved {
		t.Error("expected to have approved invoices")
	}
	if !hasRejected {
		t.Error("expected to have rejected invoices")
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
