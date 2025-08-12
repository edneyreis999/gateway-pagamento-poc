package web

import (
	"context"
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

// Server wraps the HTTP server and router configuration.
type Server struct {
	port   string
	router http.Handler
	server *http.Server
}

// NewServer builds a Server with routes configured using the provided DB and port.
func NewServer(db *sql.DB, port string) *Server {
	return &Server{
		port:   port,
		router: ConfigureRoutes(db),
	}
}

// Start starts the HTTP server and blocks until it exits.
func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}
	return s.server.ListenAndServe()
}

// Stop gracefully shuts down the server.
func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}
