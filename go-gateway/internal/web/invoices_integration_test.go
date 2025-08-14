package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"testing"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
)

// ============================================================================
// TESTES DE AUTENTICAÇÃO COM MIDDLEWARE
// ============================================================================

func TestInvoice_Create_WithoutAuth(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	// Create invoice without X-API-KEY header
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/invoices", bytes.NewBufferString(`{"amount":1000.00,"description":"Test invoice","payment_type":"credit_card","card_last_digits":"1234"}`))
	req.Header.Set("Content-Type", "application/json")
	// Não seta X-API-KEY header

	client := &http.Client{}
	invoiceResp, err := client.Do(req)
	if err != nil {
		t.Fatalf("create invoice: %v", err)
	}
	if invoiceResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", invoiceResp.StatusCode)
	}

	// Verify error message
	var errorResp map[string]any
	_ = json.NewDecoder(invoiceResp.Body).Decode(&errorResp)
	invoiceResp.Body.Close()

	if errorResp["error"] != "X-API-KEY header is required" {
		t.Errorf("expected error message 'X-API-KEY header is required', got: %v", errorResp["error"])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestInvoice_Create_WithInvalidAuth(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	// Mock GetByAPIKey call for auth middleware (falha com API key inválida)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, api_key, balance, created_at, updated_at FROM accounts WHERE api_key = $1")).
		WithArgs("invalid-api-key").WillReturnError(domain.ErrAccountNotFound)

	// Create invoice with invalid API key
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/invoices", bytes.NewBufferString(`{"amount":1000.00,"description":"Test invoice","payment_type":"credit_card","card_last_digits":"1234"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "invalid-api-key")

	client := &http.Client{}
	invoiceResp, err := client.Do(req)
	if err != nil {
		t.Fatalf("create invoice: %v", err)
	}
	if invoiceResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", invoiceResp.StatusCode)
	}

	// Verify error message
	var errorResp map[string]any
	_ = json.NewDecoder(invoiceResp.Body).Decode(&errorResp)
	invoiceResp.Body.Close()

	if errorResp["error"] != "Invalid API key" {
		t.Errorf("expected error message 'Invalid API key', got: %v", errorResp["error"])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestInvoice_Get_WithoutAuth(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	// Get invoices without X-API-KEY header
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/invoices", nil)
	// Não seta X-API-KEY header

	client := &http.Client{}
	invoiceResp, err := client.Do(req)
	if err != nil {
		t.Fatalf("get invoices: %v", err)
	}
	if invoiceResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", invoiceResp.StatusCode)
	}

	// Verify error message
	var errorResp map[string]any
	_ = json.NewDecoder(invoiceResp.Body).Decode(&errorResp)
	invoiceResp.Body.Close()

	if errorResp["error"] != "X-API-KEY header is required" {
		t.Errorf("expected error message 'X-API-KEY header is required', got: %v", errorResp["error"])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
