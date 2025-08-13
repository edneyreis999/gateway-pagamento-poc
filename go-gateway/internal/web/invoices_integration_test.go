package web

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInvoice_CreateAndGet_Success(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	// First create an account
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO accounts (id, name, email, api_key, balance, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
		WithArgs(sqlmock.AnyArg(), "John Doe", "john@example.com", sqlmock.AnyArg(), 0.0, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	accountBody := bytes.NewBufferString(`{"name":"John Doe","email":"john@example.com"}`)
	accountResp, err := http.Post(ts.URL+"/accounts", "application/json", accountBody)
	if err != nil {
		t.Fatalf("create account: %v", err)
	}
	if accountResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 got %d", accountResp.StatusCode)
	}
	var accountCreated map[string]any
	_ = json.NewDecoder(accountResp.Body).Decode(&accountCreated)
	accountResp.Body.Close()
	accountID, _ := accountCreated["id"].(string)
	apiKey, _ := accountCreated["api_key"].(string)
	if accountID == "" || apiKey == "" {
		t.Fatalf("expected id and api_key in response")
	}

	// Now create an invoice
	// First mock the GetByAPIKey call that InvoiceService makes
	now := time.Now().UTC()
	accountRows := sqlmock.NewRows([]string{"id", "name", "email", "api_key", "balance", "created_at", "updated_at"}).
		AddRow(accountID, "John Doe", "john@example.com", apiKey, 0.0, now, now)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, api_key, balance, created_at, updated_at FROM accounts WHERE api_key = $1")).
		WithArgs(apiKey).WillReturnRows(accountRows)

	// Then mock the invoice creation
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO invoices (id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)")).
		WithArgs(sqlmock.AnyArg(), accountID, 100.50, "pending", "Test invoice", "credit_card", "1234", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	invoiceBody := bytes.NewBufferString(`{"api_key":"` + apiKey + `","account_id":"` + accountID + `","amount":100.50,"description":"Test invoice","payment_type":"credit_card","card_last_digits":"1234"}`)
	invoiceResp, err := http.Post(ts.URL+"/invoices", "application/json", invoiceBody)
	if err != nil {
		t.Fatalf("create invoice: %v", err)
	}
	if invoiceResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 got %d", invoiceResp.StatusCode)
	}
	var invoiceCreated map[string]any
	_ = json.NewDecoder(invoiceResp.Body).Decode(&invoiceCreated)
	invoiceResp.Body.Close()
	invoiceID, _ := invoiceCreated["id"].(string)
	if invoiceID == "" {
		t.Fatalf("expected invoice id in response")
	}

	// Now get the invoice by ID
	now = time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "account_id", "amount", "status", "description", "payment_type", "card_last_digits", "created_at", "updated_at"}).
		AddRow(invoiceID, accountID, 100.50, "pending", "Test invoice", "credit_card", "1234", now, now)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at FROM invoices WHERE id = $1")).
		WithArgs(invoiceID).WillReturnRows(rows)

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/invoices/"+invoiceID, nil)
	req.Header.Set("X-API-KEY", apiKey)
	getResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get invoice: %v", err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", getResp.StatusCode)
	}
	getResp.Body.Close()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestInvoice_Create_InvalidJSON(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	resp, err := http.Post(ts.URL+"/invoices", "application/json", bytes.NewBufferString("{"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
	}
	resp.Body.Close()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestInvoice_Create_InvalidAmount(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	// Mock the GetByAPIKey call that InvoiceService makes
	now := time.Now().UTC()
	accountRows := sqlmock.NewRows([]string{"id", "name", "email", "api_key", "balance", "created_at", "updated_at"}).
		AddRow("test-id", "Test Account", "test@example.com", "test-api-key", 0.0, now, now)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, api_key, balance, created_at, updated_at FROM accounts WHERE api_key = $1")).
		WithArgs("test-api-key").WillReturnRows(accountRows)

	resp, err := http.Post(ts.URL+"/invoices", "application/json", bytes.NewBufferString(`{"api_key":"test-api-key","account_id":"test-id","amount":-50.00,"description":"Test invoice","payment_type":"credit_card"}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
	}
	resp.Body.Close()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestInvoice_Create_InvalidDescription(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	// Mock the GetByAPIKey call that InvoiceService makes
	now := time.Now().UTC()
	accountRows := sqlmock.NewRows([]string{"id", "name", "email", "api_key", "balance", "created_at", "updated_at"}).
		AddRow("test-id", "Test Account", "test@example.com", "test-api-key", 0.0, now, now)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, api_key, balance, created_at, updated_at FROM accounts WHERE api_key = $1")).
		WithArgs("test-api-key").WillReturnRows(accountRows)

	resp, err := http.Post(ts.URL+"/invoices", "application/json", bytes.NewBufferString(`{"api_key":"test-api-key","account_id":"test-id","amount":100.00,"description":"Te","payment_type":"credit_card"}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
	}
	resp.Body.Close()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestInvoice_Get_Unauthorized(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/invoices?account_id=test-id", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", resp.StatusCode)
	}
	resp.Body.Close()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestInvoice_Get_NotFound(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	invoiceID := "does-not-exist"
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at FROM invoices WHERE id = $1")).
		WithArgs(invoiceID).WillReturnError(sql.ErrNoRows)

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/invoices/"+invoiceID, nil)
	req.Header.Set("X-API-KEY", "test-api-key")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", resp.StatusCode)
	}
	resp.Body.Close()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestInvoice_GetByAccountID_Success(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	accountID := "acc-1"
	apiKey := "test-api-key"

	// Mock the query to return invoices for the account
	now := time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "account_id", "amount", "status", "description", "payment_type", "card_last_digits", "created_at", "updated_at"}).
		AddRow("inv-1", accountID, 100.00, "pending", "Invoice 1", "credit_card", "1234", now, now).
		AddRow("inv-2", accountID, 200.00, "approved", "Invoice 2", "debit_card", "5678", now, now)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at FROM invoices WHERE account_id = $1 ORDER BY created_at DESC")).
		WithArgs(accountID).WillReturnRows(rows)

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/invoices?account_id="+accountID, nil)
	req.Header.Set("X-API-KEY", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}
	resp.Body.Close()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestInvoice_GetByAccountID_Empty(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	accountID := "acc-2"
	apiKey := "test-api-key"

	// Mock the query to return no invoices
	rows := sqlmock.NewRows([]string{"id", "account_id", "amount", "status", "description", "payment_type", "card_last_digits", "created_at", "updated_at"})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at FROM invoices WHERE account_id = $1 ORDER BY created_at DESC")).
		WithArgs(accountID).WillReturnRows(rows)

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/invoices?account_id="+accountID, nil)
	req.Header.Set("X-API-KEY", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}
	resp.Body.Close()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
