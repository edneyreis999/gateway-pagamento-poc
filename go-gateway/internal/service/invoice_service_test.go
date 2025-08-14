package service

import (
	"context"
	"testing"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	"github.com/devfullcycle/imersao22/go-gateway/internal/repository/memory"
)

// Mock AccountService for testing
type mockAccountService struct {
	accounts map[string]*AccountOutput
}

func newMockAccountService() *mockAccountService {
	return &mockAccountService{
		accounts: make(map[string]*AccountOutput),
	}
}

func (m *mockAccountService) Create(ctx context.Context, in AccountCreateInput) (*AccountOutput, error) {
	// Not needed for invoice tests
	return nil, nil
}

func (m *mockAccountService) GetByID(ctx context.Context, id string) (*AccountOutput, error) {
	// Not needed for invoice tests
	return nil, nil
}

func (m *mockAccountService) GetByAPIKey(ctx context.Context, apiKey string) (*AccountOutput, error) {
	if account, exists := m.accounts[apiKey]; exists {
		return account, nil
	}
	return nil, domain.ErrAccountNotFound
}

func (m *mockAccountService) UpdateBalance(ctx context.Context, apiKey string, amount float64) error {
	// Not needed for invoice tests
	return nil
}

func (m *mockAccountService) addTestAccount(apiKey, accountID string) {
	m.accounts[apiKey] = &AccountOutput{
		ID:      accountID,
		Name:    "Test Account",
		Email:   "test@example.com",
		APIKey:  apiKey,
		Balance: 1000.0,
	}
}

func TestInvoiceService_Create(t *testing.T) {
	repo := memory.NewInvoiceRepositoryMemory()
	mockAccountSvc := newMockAccountService()

	// Add a test account
	testAPIKey := "test-api-key-123"
	testAccountID := "test-account-id"
	mockAccountSvc.addTestAccount(testAPIKey, testAccountID)

	svc := NewInvoiceServiceWithAccountService(nil, mockAccountSvc)
	svc.repo = repo // Override the repo to use memory instead of postgres

	// Test with controlled processor for approval
	t.Run("valid invoice with approval", func(t *testing.T) {
		// Create a test processor that always approves
		testProcessor := domain.NewTestInvoiceProcessor()
		testProcessor.SetNextStatus(domain.StatusApproved)
		svc.SetProcessor(testProcessor)

		input := InvoiceCreateInput{
			APIKey:         testAPIKey,
			Amount:         100.50,
			Description:    "Test invoice",
			PaymentType:    "credit_card",
			CardLastDigits: "1234",
		}

		output, err := svc.Create(context.Background(), input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if output.ID == "" {
			t.Error("expected invoice ID to be set")
		}
		if output.AccountID != testAccountID {
			t.Errorf("expected account ID %s, got %s", testAccountID, output.AccountID)
		}
		if output.Amount != input.Amount {
			t.Errorf("expected amount %f, got %f", input.Amount, output.Amount)
		}
		// After processing, status should be approved, not pending
		if output.Status != "approved" {
			t.Errorf("expected status 'approved', got '%s'", output.Status)
		}
		if output.Description != input.Description {
			t.Errorf("expected description %s, got %s", input.Description, output.Description)
		}
		if output.PaymentType != input.PaymentType {
			t.Errorf("expected payment type %s, got %s", input.PaymentType, output.PaymentType)
		}
		if output.CardLastDigits != input.CardLastDigits {
			t.Errorf("expected card last digits %s, got %s", input.CardLastDigits, output.CardLastDigits)
		}
	})

	// Test with controlled processor for rejection
	t.Run("valid invoice with rejection", func(t *testing.T) {
		// Create a test processor that always rejects
		testProcessor := domain.NewTestInvoiceProcessor()
		testProcessor.SetNextStatus(domain.StatusRejected)
		svc.SetProcessor(testProcessor)

		input := InvoiceCreateInput{
			APIKey:         testAPIKey,
			Amount:         200.00,
			Description:    "Test invoice 2",
			PaymentType:    "debit_card",
			CardLastDigits: "5678",
		}

		output, err := svc.Create(context.Background(), input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if output.ID == "" {
			t.Error("expected invoice ID to be set")
		}
		// After processing, status should be rejected, not pending
		if output.Status != "rejected" {
			t.Errorf("expected status 'rejected', got '%s'", output.Status)
		}
	})

	// Test with high value invoice that should stay pending
	t.Run("high value invoice stays pending", func(t *testing.T) {
		// Create a test processor that always approves
		testProcessor := domain.NewTestInvoiceProcessor()
		testProcessor.SetNextStatus(domain.StatusApproved)
		svc.SetProcessor(testProcessor)

		input := InvoiceCreateInput{
			APIKey:         testAPIKey,
			Amount:         15000.00, // Amount > 10000
			Description:    "High value invoice",
			PaymentType:    "credit_card",
			CardLastDigits: "9999",
		}

		output, err := svc.Create(context.Background(), input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if output.ID == "" {
			t.Error("expected invoice ID to be set")
		}
		// High value invoices should stay pending regardless of processor setting
		if output.Status != "pending" {
			t.Errorf("expected status 'pending' for high value invoice, got '%s'", output.Status)
		}
	})

	// Test with exact value invoice (10000) that should be processed
	t.Run("exact value invoice gets processed", func(t *testing.T) {
		// Create a test processor that always rejects
		testProcessor := domain.NewTestInvoiceProcessor()
		testProcessor.SetNextStatus(domain.StatusRejected)
		svc.SetProcessor(testProcessor)

		input := InvoiceCreateInput{
			APIKey: testAPIKey,

			Amount:         10000.00, // Amount = 10000
			Description:    "Exact value invoice",
			PaymentType:    "credit_card",
			CardLastDigits: "8888",
		}

		output, err := svc.Create(context.Background(), input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if output.ID == "" {
			t.Error("expected invoice ID to be set")
		}
		// Exact value invoices should be processed normally
		if output.Status != "rejected" {
			t.Errorf("expected status 'rejected' for exact value invoice, got '%s'", output.Status)
		}
	})

	// Test validation errors
	t.Run("validation errors", func(t *testing.T) {
		// Reset processor to default for validation tests
		svc.SetProcessor(nil)

		tests := []struct {
			name          string
			input         InvoiceCreateInput
			expectedError bool
		}{
			{
				name: "invalid description too short",
				input: InvoiceCreateInput{
					APIKey: testAPIKey,

					Amount:      100.50,
					Description: "Te",
					PaymentType: "credit_card",
				},
				expectedError: true,
			},
			{
				name: "invalid payment type empty",
				input: InvoiceCreateInput{
					APIKey: testAPIKey,

					Amount:      100.50,
					Description: "Test invoice",
					PaymentType: "",
				},
				expectedError: true,
			},
			{
				name: "invalid amount negative",
				input: InvoiceCreateInput{
					APIKey: testAPIKey,

					Amount:      -50.00,
					Description: "Test invoice",
					PaymentType: "credit_card",
				},
				expectedError: true,
			},
			{
				name: "invalid amount zero",
				input: InvoiceCreateInput{
					APIKey: testAPIKey,

					Amount:      0.00,
					Description: "Test invoice",
					PaymentType: "credit_card",
				},
				expectedError: true,
			},
			{
				name: "invalid API key",
				input: InvoiceCreateInput{
					APIKey: "invalid-api-key",

					Amount:      100.50,
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
				if output.AccountID != testAccountID {
					t.Errorf("expected account ID %s, got %s", testAccountID, output.AccountID)
				}
				if output.Amount != tt.input.Amount {
					t.Errorf("expected amount %f, got %f", tt.input.Amount, output.Amount)
				}
				// After processing, status should not be pending
				if output.Status == "pending" {
					t.Error("expected status to not be pending after processing")
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
	})
}

func TestInvoiceService_GetByID(t *testing.T) {
	repo := memory.NewInvoiceRepositoryMemory()
	mockAccountSvc := newMockAccountService()

	// Add a test account
	testAPIKey := "test-api-key-123"
	testAccountID := "test-account-id"
	mockAccountSvc.addTestAccount(testAPIKey, testAccountID)

	svc := NewInvoiceServiceWithAccountService(nil, mockAccountSvc)
	svc.repo = repo // Override the repo to use memory instead of postgres

	// Create a test invoice first with controlled processor
	testProcessor := domain.NewTestInvoiceProcessor()
	testProcessor.SetNextStatus(domain.StatusApproved)
	svc.SetProcessor(testProcessor)

	input := InvoiceCreateInput{
		APIKey: testAPIKey,

		Amount:         100.50,
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
	}

	created, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Verify that the invoice was processed and is not pending
	if created.Status == "pending" {
		t.Error("expected invoice to not be pending after processing")
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
	mockAccountSvc := newMockAccountService()

	// Add test accounts
	testAPIKey1 := "test-api-key-1"
	testAPIKey2 := "test-api-key-2"
	testAccountID1 := "account-1"
	testAccountID2 := "account-2"
	mockAccountSvc.addTestAccount(testAPIKey1, testAccountID1)
	mockAccountSvc.addTestAccount(testAPIKey2, testAccountID2)

	svc := NewInvoiceServiceWithAccountService(nil, mockAccountSvc)
	svc.repo = repo // Override the repo to use memory instead of postgres

	// Create test invoices for different accounts with controlled processors
	inputs := []InvoiceCreateInput{
		{
			APIKey: testAPIKey1,

			Amount:      100.00,
			Description: "Invoice 1",
			PaymentType: "credit_card",
		},
		{
			APIKey: testAPIKey1,

			Amount:      200.00,
			Description: "Invoice 2",
			PaymentType: "debit_card",
		},
		{
			APIKey: testAPIKey2,

			Amount:      150.00,
			Description: "Invoice 3",
			PaymentType: "credit_card",
		},
	}

	// Create invoices with different statuses for testing
	statuses := []domain.Status{domain.StatusApproved, domain.StatusRejected, domain.StatusApproved}

	for i, input := range inputs {
		// Set processor for this specific invoice
		testProcessor := domain.NewTestInvoiceProcessor()
		testProcessor.SetNextStatus(statuses[i])
		svc.SetProcessor(testProcessor)

		_, err := svc.Create(context.Background(), input)
		if err != nil {
			t.Fatalf("failed to create test invoice %d: %v", i, err)
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

	// Verify that invoices are not pending
	for i, invoice := range invoices {
		if invoice.Status == "pending" {
			t.Errorf("invoice %d should not be pending after processing, got status: %s", i, invoice.Status)
		}
	}

	// Test getting invoices for account-2
	invoices, err = svc.GetByAccountID(context.Background(), "account-2")
	if err != nil {
		t.Errorf("failed to get invoices by account ID: %v", err)
	}

	if len(invoices) != 1 {
		t.Errorf("expected 1 invoice for account-2, got %d", len(invoices))
	}

	// Verify that invoice is not pending
	if invoices[0].Status == "pending" {
		t.Errorf("invoice should not be pending after processing, got status: %s", invoices[0].Status)
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
	mockAccountSvc := newMockAccountService()

	// Add a test account
	testAPIKey := "test-api-key-123"
	testAccountID := "test-account-id"
	mockAccountSvc.addTestAccount(testAPIKey, testAccountID)

	svc := NewInvoiceServiceWithAccountService(nil, mockAccountSvc)
	svc.repo = repo // Override the repo to use memory instead of postgres

	// Create a test invoice first with controlled processor
	testProcessor := domain.NewTestInvoiceProcessor()
	testProcessor.SetNextStatus(domain.StatusRejected) // Start with rejected
	svc.SetProcessor(testProcessor)

	input := InvoiceCreateInput{
		APIKey: testAPIKey,

		Amount:         100.50,
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
	}

	created, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Verify that the invoice was processed and is not pending
	if created.Status == "pending" {
		t.Error("expected invoice to not be pending after processing")
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

func TestInvoiceService_GetAccountByAPIKey(t *testing.T) {
	repo := memory.NewInvoiceRepositoryMemory()
	mockAccountSvc := newMockAccountService()

	// Add a test account
	testAPIKey := "test-api-key-123"
	testAccountID := "test-account-id"
	mockAccountSvc.addTestAccount(testAPIKey, testAccountID)

	svc := NewInvoiceServiceWithAccountService(nil, mockAccountSvc)
	svc.repo = repo

	// Test successful retrieval
	account, err := svc.GetAccountByAPIKey(context.Background(), testAPIKey)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if account.ID != testAccountID {
		t.Errorf("expected account ID %s, got %s", testAccountID, account.ID)
	}

	if account.APIKey != testAPIKey {
		t.Errorf("expected API key %s, got %s", testAPIKey, account.APIKey)
	}

	// Test non-existent API key
	_, err = svc.GetAccountByAPIKey(context.Background(), "non-existent-key")
	if err != domain.ErrAccountNotFound {
		t.Errorf("expected ErrAccountNotFound, got %v", err)
	}
}

func TestInvoiceService_GetByID_NotFound(t *testing.T) {
	repo := memory.NewInvoiceRepositoryMemory()
	mockAccountSvc := newMockAccountService()

	svc := NewInvoiceServiceWithAccountService(nil, mockAccountSvc)
	svc.repo = repo

	// Test getting non-existent invoice
	_, err := svc.GetByID(context.Background(), "non-existent-id")
	if err != domain.ErrInvoiceNotFound {
		t.Errorf("expected ErrInvoiceNotFound, got %v", err)
	}
}

func TestInvoiceService_GetByID_Success(t *testing.T) {
	repo := memory.NewInvoiceRepositoryMemory()
	mockAccountSvc := newMockAccountService()

	// Add a test account
	testAPIKey := "test-api-key-123"
	testAccountID := "test-account-id"
	mockAccountSvc.addTestAccount(testAPIKey, testAccountID)

	svc := NewInvoiceServiceWithAccountService(nil, mockAccountSvc)
	svc.repo = repo

	// Create a test invoice first
	testProcessor := domain.NewTestInvoiceProcessor()
	testProcessor.SetNextStatus(domain.StatusApproved)
	svc.SetProcessor(testProcessor)

	input := InvoiceCreateInput{
		APIKey:         testAPIKey,
		Amount:         100.50,
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
	}

	created, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("failed to create test invoice: %v", err)
	}

	// Test getting the created invoice
	retrieved, err := svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("expected invoice ID %s, got %s", created.ID, retrieved.ID)
	}

	if retrieved.Amount != created.Amount {
		t.Errorf("expected amount %f, got %f", created.Amount, retrieved.Amount)
	}

	if retrieved.Status != created.Status {
		t.Errorf("expected status %s, got %s", created.Status, retrieved.Status)
	}
}

func TestInvoiceService_Create_AccountNotFound(t *testing.T) {
	repo := memory.NewInvoiceRepositoryMemory()
	mockAccountSvc := newMockAccountService()

	svc := NewInvoiceServiceWithAccountService(nil, mockAccountSvc)
	svc.repo = repo

	input := InvoiceCreateInput{
		APIKey:         "non-existent-api-key",
		Amount:         100.50,
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
	}

	_, err := svc.Create(context.Background(), input)
	if err != domain.ErrAccountNotFound {
		t.Errorf("expected ErrAccountNotFound, got %v", err)
	}
}

func TestInvoiceService_Create_InvalidInput(t *testing.T) {
	repo := memory.NewInvoiceRepositoryMemory()
	mockAccountSvc := newMockAccountService()

	// Add a test account
	testAPIKey := "test-api-key-123"
	testAccountID := "test-account-id"
	mockAccountSvc.addTestAccount(testAPIKey, testAccountID)

	svc := NewInvoiceServiceWithAccountService(nil, mockAccountSvc)
	svc.repo = repo

	// Test with negative amount
	input1 := InvoiceCreateInput{
		APIKey:         testAPIKey,
		Amount:         -100.50,
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
	}

	_, err := svc.Create(context.Background(), input1)
	if err != domain.ErrInvoiceNegativeValue {
		t.Errorf("expected ErrInvoiceNegativeValue, got %v", err)
	}

	// Test with zero amount
	input2 := InvoiceCreateInput{
		APIKey:         testAPIKey,
		Amount:         0.0,
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
	}

	_, err = svc.Create(context.Background(), input2)
	if err != domain.ErrInvoiceNegativeValue {
		t.Errorf("expected ErrInvoiceNegativeValue, got %v", err)
	}

	// Test with short description
	input3 := InvoiceCreateInput{
		APIKey:         testAPIKey,
		Amount:         100.50,
		Description:    "ab",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
	}

	_, err = svc.Create(context.Background(), input3)
	if err != domain.ErrInvalidDescription {
		t.Errorf("expected ErrInvalidDescription, got %v", err)
	}

	// Test with empty payment type
	input4 := InvoiceCreateInput{
		APIKey:         testAPIKey,
		Amount:         100.50,
		Description:    "Test invoice",
		PaymentType:    "",
		CardLastDigits: "1234",
	}

	_, err = svc.Create(context.Background(), input4)
	if err != domain.ErrInvalidPaymentType {
		t.Errorf("expected ErrInvalidPaymentType, got %v", err)
	}
}

func TestInvoiceService_Create_RepositoryError(t *testing.T) {
	// Create a mock repository that always returns an error
	mockRepo := &mockInvoiceRepository{
		createError: domain.ErrInvoiceNotFound, // Use any domain error
	}
	mockAccountSvc := newMockAccountService()

	// Add a test account
	testAPIKey := "test-api-key-123"
	testAccountID := "test-account-id"
	mockAccountSvc.addTestAccount(testAPIKey, testAccountID)

	svc := NewInvoiceServiceWithAccountService(nil, mockAccountSvc)
	svc.repo = mockRepo // Override the repo to use our mock

	// Create a test invoice with controlled processor
	testProcessor := domain.NewTestInvoiceProcessor()
	testProcessor.SetNextStatus(domain.StatusApproved)
	svc.SetProcessor(testProcessor)

	input := InvoiceCreateInput{
		APIKey:         testAPIKey,
		Amount:         100.50,
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
	}

	_, err := svc.Create(context.Background(), input)
	if err != domain.ErrInvoiceNotFound {
		t.Errorf("expected ErrInvoiceNotFound, got %v", err)
	}
}

// Mock repository for testing repository errors
type mockInvoiceRepository struct {
	createError error
}

func (m *mockInvoiceRepository) Create(ctx context.Context, invoice *domain.Invoice) error {
	return m.createError
}

func (m *mockInvoiceRepository) GetByID(ctx context.Context, id string) (*domain.Invoice, error) {
	return nil, domain.ErrInvoiceNotFound
}

func (m *mockInvoiceRepository) GetByAccountID(ctx context.Context, accountID string) ([]*domain.Invoice, error) {
	return nil, nil
}

func (m *mockInvoiceRepository) UpdateStatus(ctx context.Context, id string, status domain.Status) error {
	return domain.ErrInvoiceNotFound
}
