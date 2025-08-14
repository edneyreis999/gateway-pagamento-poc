package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	"github.com/devfullcycle/imersao22/go-gateway/internal/service"
)

// MockInvoiceService is a mock implementation of InvoiceServicePort for testing
type MockInvoiceService struct {
	invoices map[string]*service.InvoiceOutput
	accounts map[string]*service.AccountOutput
}

func NewMockInvoiceService() *MockInvoiceService {
	return &MockInvoiceService{
		invoices: make(map[string]*service.InvoiceOutput),
		accounts: make(map[string]*service.AccountOutput),
	}
}

func (m *MockInvoiceService) Create(ctx context.Context, in service.InvoiceCreateInput) (*service.InvoiceOutput, error) {
	// Simulate domain validation
	if in.Amount <= 0 {
		return nil, domain.ErrInvoiceNegativeValue
	}
	if len(in.Description) < 3 {
		return nil, domain.ErrInvalidDescription
	}
	if len(in.PaymentType) == 0 {
		return nil, domain.ErrInvalidPaymentType
	}

	// Simulate creation
	invoice := &service.InvoiceOutput{
		ID:             "test-invoice-id",
		AccountID:      "test-account-id", // Mock account ID
		Amount:         in.Amount,
		Status:         "pending",
		Description:    in.Description,
		PaymentType:    in.PaymentType,
		CardLastDigits: in.CardLastDigits,
	}

	m.invoices[invoice.ID] = invoice
	return invoice, nil
}

func (m *MockInvoiceService) GetByID(ctx context.Context, id string) (*service.InvoiceOutput, error) {
	invoice, exists := m.invoices[id]
	if !exists {
		return nil, domain.ErrInvoiceNotFound
	}
	return invoice, nil
}

func (m *MockInvoiceService) GetByAccountID(ctx context.Context, accountID string) ([]*service.InvoiceOutput, error) {
	var invoices []*service.InvoiceOutput
	for _, invoice := range m.invoices {
		if invoice.AccountID == accountID {
			invoices = append(invoices, invoice)
		}
	}
	return invoices, nil
}

func (m *MockInvoiceService) GetAccountByAPIKey(ctx context.Context, apiKey string) (*service.AccountOutput, error) {
	// Return a mock account
	return &service.AccountOutput{
		ID:        "test-account-id",
		Name:      "Test Account",
		Email:     "test@example.com",
		APIKey:    apiKey,
		Balance:   1000.0,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func TestInvoiceHandler_CreateInvoice(t *testing.T) {
	tests := []struct {
		name           string
		input          service.InvoiceCreateInput
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "valid invoice creation",
			input: service.InvoiceCreateInput{
				APIKey:         "test-api-key-123",
				Amount:         100.50,
				Description:    "Test invoice",
				PaymentType:    "credit_card",
				CardLastDigits: "1234",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "invalid amount",
			input: service.InvoiceCreateInput{
				APIKey:      "test-api-key-123",
				Amount:      -50.00,
				Description: "Test invoice",
				PaymentType: "credit_card",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "invalid description",
			input: service.InvoiceCreateInput{
				APIKey:      "test-api-key-123",
				Amount:      100.00,
				Description: "Te",
				PaymentType: "credit_card",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := NewMockInvoiceService()
			handler := NewInvoiceHandler(mockSvc)

			inputJSON, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/invoices", bytes.NewBuffer(inputJSON))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-API-KEY", "test-api-key-123")

			w := httptest.NewRecorder()
			handler.PostInvoices()(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectedError && w.Code == http.StatusCreated {
				var response service.InvoiceOutput
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}

				if response.ID == "" {
					t.Error("expected invoice ID to be set")
				}
				if response.Status != "pending" {
					t.Errorf("expected status 'pending', got '%s'", response.Status)
				}
			}
		})
	}
}

func TestInvoiceHandler_GetInvoicesByAccountID(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	handler := NewInvoiceHandler(mockSvc)

	// Create a test invoice first
	testInvoice := &service.InvoiceOutput{
		ID:             "test-invoice-id",
		AccountID:      "test-account-id",
		Amount:         100.00,
		Status:         "pending",
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
	}
	mockSvc.invoices[testInvoice.ID] = testInvoice

	req := httptest.NewRequest(http.MethodGet, "/invoices?account_id=test-account-id", nil)
	req.Header.Set("X-API-KEY", "test-api-key")

	w := httptest.NewRecorder()
	handler.GetInvoices()(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []*service.InvoiceOutput
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("expected 1 invoice, got %d", len(response))
	}

	if response[0].ID != testInvoice.ID {
		t.Errorf("expected invoice ID %s, got %s", testInvoice.ID, response[0].ID)
	}
}

func TestInvoiceHandler_GetInvoiceByID(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	handler := NewInvoiceHandler(mockSvc)

	// Create a test invoice first
	testInvoice := &service.InvoiceOutput{
		ID:             "test-invoice-id",
		AccountID:      "test-account-id",
		Amount:         100.00,
		Status:         "pending",
		Description:    "Test invoice",
		PaymentType:    "credit_card",
		CardLastDigits: "1234",
	}
	mockSvc.invoices[testInvoice.ID] = testInvoice

	req := httptest.NewRequest(http.MethodGet, "/invoices/test-invoice-id", nil)
	req.Header.Set("X-API-KEY", "test-api-key")

	w := httptest.NewRecorder()
	handler.GetInvoiceByID()(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response service.InvoiceOutput
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response.ID != testInvoice.ID {
		t.Errorf("expected invoice ID %s, got %s", testInvoice.ID, response.ID)
	}
}

func TestInvoiceHandler_MissingAPIKey(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	handler := NewInvoiceHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/invoices?account_id=test-account-id", nil)
	// No X-API-KEY header

	w := httptest.NewRecorder()
	handler.GetInvoices()(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestInvoiceHandler_MethodNotAllowed(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	handler := NewInvoiceHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/invoices", nil)
	req.Header.Set("X-API-KEY", "test-api-key")

	w := httptest.NewRecorder()
	handler.GetInvoices()(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}
