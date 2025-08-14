package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	"github.com/devfullcycle/imersao22/go-gateway/internal/service"
)

// MockInvoiceService is a mock implementation of InvoiceServicePort for testing
type MockInvoiceService struct {
	invoices            map[string]*service.InvoiceOutput
	accounts            map[string]*service.AccountOutput
	createError         error
	getByIDError        error
	getByAccountIDError error
}

func NewMockInvoiceService() *MockInvoiceService {
	return &MockInvoiceService{
		invoices: make(map[string]*service.InvoiceOutput),
		accounts: make(map[string]*service.AccountOutput),
	}
}

func (m *MockInvoiceService) Create(ctx context.Context, in service.InvoiceCreateInput) (*service.InvoiceOutput, error) {
	// Check for mock errors first
	if m.createError != nil {
		return nil, m.createError
	}

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
	if m.getByIDError != nil {
		return nil, m.getByIDError
	}
	invoice, exists := m.invoices[id]
	if !exists {
		return nil, domain.ErrInvoiceNotFound
	}
	return invoice, nil
}

func (m *MockInvoiceService) GetByAccountID(ctx context.Context, accountID string) ([]*service.InvoiceOutput, error) {
	if m.getByAccountIDError != nil {
		return nil, m.getByAccountIDError
	}
	var invoices []*service.InvoiceOutput
	for _, invoice := range m.invoices {
		if invoice.AccountID == accountID {
			invoices = append(invoices, invoice)
		}
	}
	return invoices, nil
}

func (m *MockInvoiceService) GetAccountByAPIKey(ctx context.Context, apiKey string) (*service.AccountOutput, error) {
	account, exists := m.accounts[apiKey]
	if !exists {
		return nil, domain.ErrAccountNotFound
	}
	return account, nil
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

	// Add test account first
	testAccount := &service.AccountOutput{
		ID:      "test-account-id",
		Name:    "Test Account",
		Email:   "test@example.com",
		APIKey:  "test-api-key",
		Balance: 1000.0,
	}
	mockSvc.accounts["test-api-key"] = testAccount

	// Create a test invoice
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

	if len(response) > 0 && response[0].ID != testInvoice.ID {
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

// ============================================================================
// TESTES ADICIONAIS PARA COBERTURA COMPLETA
// ============================================================================

func TestInvoiceHandler_GetInvoiceByID_NotFound(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	handler := NewInvoiceHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/invoices/non-existent-id", nil)
	req.Header.Set("X-API-KEY", "test-api-key")

	w := httptest.NewRecorder()
	handler.GetInvoiceByID()(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestInvoiceHandler_GetInvoiceByID_InvalidID(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	handler := NewInvoiceHandler(mockSvc)

	// Test with invalid URL path
	req := httptest.NewRequest(http.MethodGet, "/invoices/", nil)
	req.Header.Set("X-API-KEY", "test-api-key")

	w := httptest.NewRecorder()
	handler.GetInvoiceByID()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestInvoiceHandler_GetInvoiceByID_EmptyID(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	handler := NewInvoiceHandler(mockSvc)

	// Test with empty ID in path
	req := httptest.NewRequest(http.MethodGet, "/invoices//", nil)
	req.Header.Set("X-API-KEY", "test-api-key")

	w := httptest.NewRecorder()
	handler.GetInvoiceByID()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestInvoiceHandler_GetInvoices_AccountNotFound(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	// Don't add any test account, so GetAccountByAPIKey will fail
	handler := NewInvoiceHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/invoices?account_id=test-account-id", nil)
	req.Header.Set("X-API-KEY", "test-api-key")

	w := httptest.NewRecorder()
	handler.GetInvoices()(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestInvoiceHandler_GetInvoices_EmptyResult(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	// Add test account but no invoices
	testAccount := &service.AccountOutput{
		ID:      "test-account-id",
		Name:    "Test Account",
		Email:   "test@example.com",
		APIKey:  "test-api-key",
		Balance: 1000.0,
	}
	mockSvc.accounts["test-api-key"] = testAccount
	handler := NewInvoiceHandler(mockSvc)

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

	if len(response) != 0 {
		t.Errorf("expected 0 invoices, got %d", len(response))
	}
}

func TestInvoiceHandler_CreateInvoice_ServiceError(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	// Add test account
	testAccount := &service.AccountOutput{
		ID:      "test-account-id",
		Name:    "Test Account",
		Email:   "test@example.com",
		APIKey:  "test-api-key",
		Balance: 1000.0,
	}
	mockSvc.accounts["test-api-key"] = testAccount
	// Set service to return error for Create
	mockSvc.createError = errors.New("service error")
	handler := NewInvoiceHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/invoices", bytes.NewBufferString(`{"amount":100.00,"description":"Test invoice","payment_type":"credit_card"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "test-api-key")

	w := httptest.NewRecorder()
	handler.PostInvoices()(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestInvoiceHandler_CreateInvoice_DomainValidationError(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	// Add test account
	testAccount := &service.AccountOutput{
		ID:      "test-account-id",
		Name:    "Test Account",
		Email:   "test@example.com",
		APIKey:  "test-api-key",
		Balance: 1000.0,
	}
	mockSvc.accounts["test-api-key"] = testAccount
	// Set service to return domain validation error
	mockSvc.createError = domain.ErrInvalidDescription
	handler := NewInvoiceHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/invoices", bytes.NewBufferString(`{"amount":100.00,"description":"ab","payment_type":"credit_card"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "test-api-key")

	w := httptest.NewRecorder()
	handler.PostInvoices()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestInvoiceHandler_CreateInvoice_AccountNotFoundError(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	// Don't add test account, so GetAccountByAPIKey will fail
	handler := NewInvoiceHandler(mockSvc)

	// Debug: verify mock is working
	_, err := mockSvc.GetAccountByAPIKey(context.Background(), "test-api-key")
	if err != domain.ErrAccountNotFound {
		t.Fatalf("mock should return ErrAccountNotFound, got: %v", err)
	}

	// Debug: verify handler is using mock
	if handler.svc != mockSvc {
		t.Fatalf("handler is not using the mock service")
	}

	// Debug: verify mock implements interface
	var _ InvoiceServicePort = mockSvc

	// Debug: verify mock is being called
	req := httptest.NewRequest(http.MethodPost, "/invoices", bytes.NewBufferString(`{"amount":100.00,"description":"Test invoice","payment_type":"credit_card"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "test-api-key")

	w := httptest.NewRecorder()
	handler.PostInvoices()(w, req)

	// Debug: print response body
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())

	// Debug: print handler service
	t.Logf("Handler service: %T", handler.svc)
	t.Logf("Mock service: %T", mockSvc)

	// Debug: print mock createError
	t.Logf("Mock createError: %v", mockSvc.createError)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestInvoiceHandler_GetInvoiceByID_ServiceError(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	// Add test account
	testAccount := &service.AccountOutput{
		ID:      "test-account-id",
		Name:    "Test Account",
		Email:   "test@example.com",
		APIKey:  "test-api-key",
		Balance: 1000.0,
	}
	mockSvc.accounts["test-api-key"] = testAccount
	// Set service to return error for GetByID
	mockSvc.getByIDError = errors.New("service error")
	handler := NewInvoiceHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/invoices/test-invoice-id", nil)
	req.Header.Set("X-API-KEY", "test-api-key")

	w := httptest.NewRecorder()
	handler.GetInvoiceByID()(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestInvoiceHandler_GetInvoices_ServiceError(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	// Add test account
	testAccount := &service.AccountOutput{
		ID:      "test-account-id",
		Name:    "Test Account",
		Email:   "test@example.com",
		APIKey:  "test-api-key",
		Balance: 1000.0,
	}
	mockSvc.accounts["test-api-key"] = testAccount
	// Set service to return error for GetByAccountID
	mockSvc.getByAccountIDError = errors.New("service error")
	handler := NewInvoiceHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/invoices?account_id=test-account-id", nil)
	req.Header.Set("X-API-KEY", "test-api-key")

	w := httptest.NewRecorder()
	handler.GetInvoices()(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestInvoiceHandler_GetInvoiceByID_MethodNotAllowed(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	handler := NewInvoiceHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/invoices/test-invoice-id", nil)
	req.Header.Set("X-API-KEY", "test-api-key")

	w := httptest.NewRecorder()
	handler.GetInvoiceByID()(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestInvoiceHandler_GetInvoices_MethodNotAllowed(t *testing.T) {
	mockSvc := NewMockInvoiceService()
	handler := NewInvoiceHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/invoices", bytes.NewBufferString(`{"amount":100.00,"description":"Test invoice","payment_type":"credit_card"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "test-api-key")

	w := httptest.NewRecorder()
	handler.GetInvoices()(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}
