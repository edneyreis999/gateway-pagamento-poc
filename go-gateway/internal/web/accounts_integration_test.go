package web

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// helper to spin up test server with sqlmock DB
func newTestServer(t *testing.T) (*httptest.Server, sqlmock.Sqlmock, *sql.DB) {
	b, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	handler := ConfigureRoutes(b)
	ts := httptest.NewServer(handler)
	return ts, mock, b
}

func TestAccount_CreateAndGet_Success(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO accounts (id, name, email, api_key, balance, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)")).
		WithArgs(sqlmock.AnyArg(), "John Doe", "john@example.com", sqlmock.AnyArg(), 0.0, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	body := bytes.NewBufferString(`{"name":"John Doe","email":"john@example.com"}`)
	resp, err := http.Post(ts.URL+"/accounts", "application/json", body)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 got %d", resp.StatusCode)
	}
	var created map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()
	apiKey, _ := created["api_key"].(string)
	if apiKey == "" {
		t.Fatalf("expected api_key in response")
	}

	now := time.Now().UTC()
	rows := sqlmock.NewRows([]string{"id", "name", "email", "api_key", "balance", "created_at", "updated_at"}).
		AddRow(created["id"], created["name"], created["email"], apiKey, 0.0, now, now)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, api_key, balance, created_at, updated_at FROM accounts WHERE api_key = $1")).
		WithArgs(apiKey).WillReturnRows(rows)

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/accounts", nil)
	req.Header.Set("X-API-KEY", apiKey)
	getResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", getResp.StatusCode)
	}
	getResp.Body.Close()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestAccount_Create_InvalidJSON(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	resp, err := http.Post(ts.URL+"/accounts", "application/json", bytes.NewBufferString("{"))
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

func TestAccount_Create_InvalidName(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	resp, err := http.Post(ts.URL+"/accounts", "application/json", bytes.NewBufferString(`{"name":"A","email":"a@b.com"}`))
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

func TestAccount_Get_Unauthorized(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/accounts", nil)
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

func TestAccount_Get_NotFound(t *testing.T) {
	ts, mock, db := newTestServer(t)
	defer ts.Close()
	defer db.Close()

	apiKey := "does-not-exist"
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, api_key, balance, created_at, updated_at FROM accounts WHERE api_key = $1")).
		WithArgs(apiKey).WillReturnError(sql.ErrNoRows)

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/accounts", nil)
	req.Header.Set("X-API-KEY", apiKey)
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
