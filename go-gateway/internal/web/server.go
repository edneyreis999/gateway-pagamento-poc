package web

import (
	"database/sql"
	"net/http"

	"github.com/devfullcycle/imersao22/go-gateway/internal/service"
	"github.com/devfullcycle/imersao22/go-gateway/internal/web/handlers"
	"github.com/go-chi/chi/v5"
)

// ConfigureRoutes wires HTTP routes using chi mux and provided dependencies.
func ConfigureRoutes(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	// Services
	svc := service.NewAccountService(db)

	// Handlers
	h := handlers.NewAccountHandler(svc)

	// Routes
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", h.PostAccounts()) // POST /accounts
		r.Get("/", h.GetAccounts())   // GET /accounts (by X-API-KEY)
	})

	return r
}
