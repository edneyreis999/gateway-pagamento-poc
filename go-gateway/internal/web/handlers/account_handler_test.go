package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	"github.com/devfullcycle/imersao22/go-gateway/internal/service"
)

type fakeSvc struct {
	create      func(ctx context.Context, in service.AccountCreateInput) (*service.AccountOutput, error)
	getByAPIKey func(ctx context.Context, apiKey string) (*service.AccountOutput, error)
}

func (f *fakeSvc) Create(ctx context.Context, in service.AccountCreateInput) (*service.AccountOutput, error) {
	return f.create(ctx, in)
}
func (f *fakeSvc) GetByAPIKey(ctx context.Context, apiKey string) (*service.AccountOutput, error) {
	return f.getByAPIKey(ctx, apiKey)
}

func TestAccountHandler_Create(t *testing.T) {
	svc := &fakeSvc{
		create: func(ctx context.Context, in service.AccountCreateInput) (*service.AccountOutput, error) {
			return &service.AccountOutput{ID: "1", Name: in.Name, Email: in.Email, APIKey: "k", Balance: 0, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
		},
	}
	h := NewAccountHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body, _ := json.Marshal(service.AccountCreateInput{Name: "Acme", Email: "acme@example.com"})
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestAccountHandler_GetByAPIKey_NotFound(t *testing.T) {
	svc := &fakeSvc{
		getByAPIKey: func(ctx context.Context, apiKey string) (*service.AccountOutput, error) {
			return nil, domain.ErrAccountNotFound
		},
	}
	h := NewAccountHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/accounts/any", nil)
	req.Header.Set("X-API-KEY", "nope")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestAccountHandler_GetByAPIKey_MissingHeader(t *testing.T) {
	svc := &fakeSvc{}
	h := NewAccountHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/accounts/any", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAccountHandler_Create_InvalidJSON(t *testing.T) {
	svc := &fakeSvc{
		create: func(ctx context.Context, in service.AccountCreateInput) (*service.AccountOutput, error) {
			return nil, errors.New("should not be called")
		},
	}
	h := NewAccountHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte("{")))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
