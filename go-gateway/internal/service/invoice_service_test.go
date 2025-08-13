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

func (m *mockAccountService) UpdateBalance(ctx context.Context, id string, amount float64) error {
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

	tests := []struct {
		name          string
		input         InvoiceCreateInput
		expectedError bool
	}{
		{
			name: "valid invoice",
			input: InvoiceCreateInput{
				APIKey:         testAPIKey,
				AccountID:      testAccountID,
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
				APIKey:      testAPIKey,
				AccountID:   testAccountID,
				Amount:      100.50,
				Description: "Te",
				PaymentType: "credit_card",
			},
			expectedError: true,
		},
		{
			name: "invalid payment type empty",
			input: InvoiceCreateInput{
				APIKey:      testAPIKey,
				AccountID:   testAccountID,
				Amount:      100.50,
				Description: "Test invoice",
				PaymentType: "",
			},
			expectedError: true,
		},
		{
			name: "invalid amount negative",
			input: InvoiceCreateInput{
				APIKey:      testAPIKey,
				AccountID:   testAccountID,
				Amount:      -50.00,
				Description: "Test invoice",
				PaymentType: "credit_card",
			},
			expectedError: true,
		},
		{
			name: "invalid amount zero",
			input: InvoiceCreateInput{
				APIKey:      testAPIKey,
				AccountID:   testAccountID,
				Amount:      0.00,
				Description: "Test invoice",
				PaymentType: "credit_card",
			},
			expectedError: true,
		},
		{
			name: "invalid API key",
			input: InvoiceCreateInput{
				APIKey:      "invalid-api-key",
				AccountID:   testAccountID,
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
	mockAccountSvc := newMockAccountService()

	// Add a test account
	testAPIKey := "test-api-key-123"
	testAccountID := "test-account-id"
	mockAccountSvc.addTestAccount(testAPIKey, testAccountID)

	svc := NewInvoiceServiceWithAccountService(nil, mockAccountSvc)
	svc.repo = repo // Override the repo to use memory instead of postgres

	// Create a test invoice first
	input := InvoiceCreateInput{
		APIKey:         testAPIKey,
		AccountID:      testAccountID,
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

	// Create test invoices for different accounts
	inputs := []InvoiceCreateInput{
		{
			APIKey:      testAPIKey1,
			AccountID:   testAccountID1,
			Amount:      100.00,
			Description: "Invoice 1",
			PaymentType: "credit_card",
		},
		{
			APIKey:      testAPIKey1,
			AccountID:   testAccountID1,
			Amount:      200.00,
			Description: "Invoice 2",
			PaymentType: "debit_card",
		},
		{
			APIKey:      testAPIKey2,
			AccountID:   testAccountID2,
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
	mockAccountSvc := newMockAccountService()

	// Add a test account
	testAPIKey := "test-api-key-123"
	testAccountID := "test-account-id"
	mockAccountSvc.addTestAccount(testAPIKey, testAccountID)

	svc := NewInvoiceServiceWithAccountService(nil, mockAccountSvc)
	svc.repo = repo // Override the repo to use memory instead of postgres

	// Create a test invoice first
	input := InvoiceCreateInput{
		APIKey:         testAPIKey,
		AccountID:      testAccountID,
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
