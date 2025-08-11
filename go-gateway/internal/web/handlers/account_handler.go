package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	"github.com/devfullcycle/imersao22/go-gateway/internal/service"
)

// AccountServicePort defines only the methods needed by the handler.
// It matches methods in service.AccountService.
type AccountServicePort interface {
	Create(ctx context.Context, in service.AccountCreateInput) (*service.AccountOutput, error)
	GetByAPIKey(ctx context.Context, apiKey string) (*service.AccountOutput, error)
}

// AccountHandler handles HTTP requests for accounts.
type AccountHandler struct {
	svc AccountServicePort
}

func NewAccountHandler(svc AccountServicePort) *AccountHandler {
	return &AccountHandler{svc: svc}
}

// PostAccounts returns a handler for POST /accounts
func (h *AccountHandler) PostAccounts() http.HandlerFunc {
	return h.handleAccounts
}

// GetAccounts returns a handler for GET /accounts (via X-API-KEY)
func (h *AccountHandler) GetAccounts() http.HandlerFunc {
	return h.handleAccounts
}

// RegisterRoutes registers the HTTP handlers on a mux.
func (h *AccountHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/accounts", h.handleAccounts)
	mux.HandleFunc("/accounts/", h.handleAccountByID)
}

// POST /accounts, GET /accounts (list not implemented, will 405)
func (h *AccountHandler) handleAccounts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createAccount(w, r)
	case http.MethodGet:
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing X-API-KEY header"})
			return
		}
		out, err := h.svc.GetByAPIKey(r.Context(), apiKey)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, domain.ErrAccountNotFound) {
				status = http.StatusNotFound
			}
			w.WriteHeader(status)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(out)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
	}
}

// GET /accounts/{id}
func (h *AccountHandler) handleAccountByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}
	apiKey := r.Header.Get("X-API-KEY")
	if apiKey == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing X-API-KEY header"})
		return
	}
	out, err := h.svc.GetByAPIKey(r.Context(), apiKey)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrAccountNotFound) {
			status = http.StatusNotFound
		}
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}

func (h *AccountHandler) createAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var in service.AccountCreateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
		return
	}
	out, err := h.svc.Create(r.Context(), in)
	if err != nil {
		status := http.StatusBadRequest
		// Domain validation errors should map to 400; others 500
		switch {
		case errors.Is(err, domain.ErrInvalidName), errors.Is(err, domain.ErrInvalidEmail):
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(out)
}
