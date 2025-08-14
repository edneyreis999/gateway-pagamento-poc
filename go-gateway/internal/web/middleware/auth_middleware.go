package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/devfullcycle/imersao22/go-gateway/internal/domain"
	"github.com/devfullcycle/imersao22/go-gateway/internal/service"
)

type AuthMiddleware struct {
	accountService *service.AccountService
}

func NewAuthMiddleware(accountService *service.AccountService) *AuthMiddleware {
	return &AuthMiddleware{
		accountService: accountService,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "X-API-KEY header is required"})
			return
		}

		_, err := m.accountService.GetByAPIKey(r.Context(), apiKey)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")

			if err == domain.ErrAccountNotFound {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid API key"})
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}

		// Authentication successful, proceed to next handler
		next.ServeHTTP(w, r)
	})
}
