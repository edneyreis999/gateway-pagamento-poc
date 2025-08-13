package web

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	"github.com/devfullcycle/imersao22/go-gateway/internal/service"
	"github.com/devfullcycle/imersao22/go-gateway/internal/web/handlers"
	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// Setup test database connection
	var err error
	testDB, err = sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/gateway_test?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer testDB.Close()

	// Run tests
	m.Run()
}

func setupTestTables(t *testing.T) {
	// Create test tables
	_, err := testDB.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			api_key VARCHAR(255) NOT NULL UNIQUE,
			balance DECIMAL(10,2) NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("failed to create accounts test table: %v", err)
	}

	_, err = testDB.Exec(`
		CREATE TABLE IF NOT EXISTS invoices (
			id UUID PRIMARY KEY,
			account_id UUID NOT NULL REFERENCES accounts(id),
			amount DECIMAL(10,2) NOT NULL,
			status VARCHAR(50) NOT NULL,
			description TEXT NOT NULL,
			payment_type VARCHAR(50) NOT NULL,
			card_last_digits VARCHAR(4),
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("failed to create invoices test table: %v", err)
	}
}

func cleanupTestTables(t *testing.T) {
	// Clean up test tables
	_, err := testDB.Exec("DELETE FROM invoices")
	if err != nil {
		t.Fatalf("failed to clean up invoices test table: %v", err)
	}

	_, err = testDB.Exec("DELETE FROM accounts")
	if err != nil {
		t.Fatalf("failed to clean up accounts test table: %v", err)
	}
}

func createTestAccount(t *testing.T) *service.AccountOutput {
	accountSvc := service.NewAccountService(testDB)

	account, err := accountSvc.Create(context.Background(), service.AccountCreateInput{
		Name:  "Test Account",
		Email: "test@example.com",
	})
	if err != nil {
		t.Fatalf("failed to create test account: %v", err)
	}

	return account
}

func TestInvoiceEndpoints_Integration(t *testing.T) {
	setupTestTables(t)
	defer cleanupTestTables(t)

	// Create test account
	account := createTestAccount(t)

	// Create invoice service and handler
	invoiceSvc := service.NewInvoiceService(testDB)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceSvc)

	t.Run("Create Invoice", func(t *testing.T) {
		input := service.InvoiceCreateInput{
			AccountID:      account.ID,
			Amount:         100.50,
			Description:    "Test invoice",
			PaymentType:    "credit_card",
			CardLastDigits: "1234",
		}

		inputJSON, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/invoices", bytes.NewBuffer(inputJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		invoiceHandler.PostInvoices()(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var response service.InvoiceOutput
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		if response.ID == "" {
			t.Error("expected invoice ID to be set")
		}
		if response.AccountID != account.ID {
			t.Errorf("expected account ID %s, got %s", account.ID, response.AccountID)
		}
		if response.Amount != input.Amount {
			t.Errorf("expected amount %f, got %f", input.Amount, response.Amount)
		}
		if response.Status != "pending" {
			t.Errorf("expected status 'pending', got '%s'", response.Status)
		}
	})

	t.Run("Get Invoices by Account ID", func(t *testing.T) {
		// Create another invoice for the same account
		input := service.InvoiceCreateInput{
			AccountID:   account.ID,
			Amount:      200.00,
			Description: "Another test invoice",
			PaymentType: "debit_card",
		}

		_, err := invoiceSvc.Create(context.Background(), input)
		if err != nil {
			t.Fatalf("failed to create second test invoice: %v", err)
		}

		// Test getting invoices by account ID
		req := httptest.NewRequest(http.MethodGet, "/invoices?account_id="+account.ID, nil)
		req.Header.Set("X-API-KEY", account.APIKey)

		w := httptest.NewRecorder()
		invoiceHandler.GetInvoices()(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []*service.InvoiceOutput
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		if len(response) != 2 {
			t.Errorf("expected 2 invoices, got %d", len(response))
		}
	})

	t.Run("Get Invoice by ID", func(t *testing.T) {
		// Create an invoice first
		input := service.InvoiceCreateInput{
			AccountID:   account.ID,
			Amount:      150.00,
			Description: "Invoice for ID test",
			PaymentType: "pix",
		}

		created, err := invoiceSvc.Create(context.Background(), input)
		if err != nil {
			t.Fatalf("failed to create test invoice: %v", err)
		}

		// Test getting invoice by ID
		req := httptest.NewRequest(http.MethodGet, "/invoices/"+created.ID, nil)
		req.Header.Set("X-API-KEY", account.APIKey)

		w := httptest.NewRecorder()
		invoiceHandler.GetInvoices()(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response service.InvoiceOutput
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		if response.ID != created.ID {
			t.Errorf("expected invoice ID %s, got %s", created.ID, response.ID)
		}
	})

	t.Run("Missing API Key", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/invoices?account_id="+account.ID, nil)
		// No X-API-KEY header

		w := httptest.NewRecorder()
		invoiceHandler.GetInvoices()(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Invalid Invoice Creation", func(t *testing.T) {
		// Test with invalid amount
		input := service.InvoiceCreateInput{
			AccountID:   account.ID,
			Amount:      -50.00,
			Description: "Invalid invoice",
			PaymentType: "credit_card",
		}

		inputJSON, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/invoices", bytes.NewBuffer(inputJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		invoiceHandler.PostInvoices()(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Method Not Allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/invoices", nil)
		req.Header.Set("X-API-KEY", account.APIKey)

		w := httptest.NewRecorder()
		invoiceHandler.GetInvoices()(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestInvoiceService_Integration(t *testing.T) {
	setupTestTables(t)
	defer cleanupTestTables(t)

	// Create test account
	account := createTestAccount(t)

	// Create invoice service
	invoiceSvc := service.NewInvoiceService(testDB)

	t.Run("Create and Retrieve Invoice", func(t *testing.T) {
		input := service.InvoiceCreateInput{
			AccountID:      account.ID,
			Amount:         100.50,
			Description:    "Test invoice",
			PaymentType:    "credit_card",
			CardLastDigits: "1234",
		}

		// Create invoice
		created, err := invoiceSvc.Create(context.Background(), input)
		if err != nil {
			t.Fatalf("failed to create invoice: %v", err)
		}

		// Retrieve by ID
		retrieved, err := invoiceSvc.GetByID(context.Background(), created.ID)
		if err != nil {
			t.Fatalf("failed to retrieve invoice: %v", err)
		}

		if retrieved.ID != created.ID {
			t.Errorf("expected invoice ID %s, got %s", created.ID, retrieved.ID)
		}
	})

	t.Run("Update Invoice Status", func(t *testing.T) {
		input := service.InvoiceCreateInput{
			AccountID:   account.ID,
			Amount:      200.00,
			Description: "Invoice for status update",
			PaymentType: "debit_card",
		}

		// Create invoice
		created, err := invoiceSvc.Create(context.Background(), input)
		if err != nil {
			t.Fatalf("failed to create invoice: %v", err)
		}

		// Update status to approved
		err = invoiceSvc.UpdateStatus(context.Background(), created.ID, domain.StatusApproved)
		if err != nil {
			t.Fatalf("failed to update status: %v", err)
		}

		// Verify status was updated
		retrieved, err := invoiceSvc.GetByID(context.Background(), created.ID)
		if err != nil {
			t.Fatalf("failed to retrieve updated invoice: %v", err)
		}

		if retrieved.Status != "approved" {
			t.Errorf("expected status 'approved', got '%s'", retrieved.Status)
		}
	})
}
