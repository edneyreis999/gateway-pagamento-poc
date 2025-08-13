package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	"github.com/devfullcycle/imersao22/go-gateway/internal/service"
)

// InvoiceServicePort defines only the methods needed by the handler.
// It matches methods in service.InvoiceService.
type InvoiceServicePort interface {
	Create(ctx context.Context, in service.InvoiceCreateInput) (*service.InvoiceOutput, error)
	GetByID(ctx context.Context, id string) (*service.InvoiceOutput, error)
	GetByAccountID(ctx context.Context, accountID string) ([]*service.InvoiceOutput, error)
}

// InvoiceHandler handles HTTP requests for invoices.
type InvoiceHandler struct {
	svc InvoiceServicePort
}

func NewInvoiceHandler(svc InvoiceServicePort) *InvoiceHandler {
	return &InvoiceHandler{svc: svc}
}

// PostInvoices returns a handler for POST /invoices
func (h *InvoiceHandler) PostInvoices() http.HandlerFunc {
	return h.handleInvoices
}

// GetInvoices returns a handler for GET /invoices (via X-API-KEY)
func (h *InvoiceHandler) GetInvoices() http.HandlerFunc {
	return h.handleInvoices
}

// GetInvoiceByID returns a handler for GET /invoices/{id}
func (h *InvoiceHandler) GetInvoiceByID() http.HandlerFunc {
	return h.handleInvoiceByID
}

// RegisterRoutes registers the HTTP handlers on a mux.
func (h *InvoiceHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/invoices", h.handleInvoices)
	mux.HandleFunc("/invoices/", h.handleInvoiceByID)
}

// POST /invoices, GET /invoices (list by account)
func (h *InvoiceHandler) handleInvoices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createInvoice(w, r)
	case http.MethodGet:
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing X-API-KEY header"})
			return
		}

		// For now, we'll get invoices by account ID from query param
		// In a real implementation, you'd get the account from the API key first
		accountID := r.URL.Query().Get("account_id")
		if accountID == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "account_id query parameter is required"})
			return
		}

		out, err := h.svc.GetByAccountID(r.Context(), accountID)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, domain.ErrInvoiceNotFound) {
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

// GET /invoices/{id}
func (h *InvoiceHandler) handleInvoiceByID(w http.ResponseWriter, r *http.Request) {
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

	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid invoice ID"})
		return
	}

	invoiceID := pathParts[2]
	if invoiceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invoice ID is required"})
		return
	}

	out, err := h.svc.GetByID(r.Context(), invoiceID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrInvoiceNotFound) {
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

func (h *InvoiceHandler) createInvoice(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var in service.InvoiceCreateInput
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
		case errors.Is(err, domain.ErrInvalidDescription),
			errors.Is(err, domain.ErrInvalidPaymentType),
			errors.Is(err, domain.ErrInvoiceNegativeValue):
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
