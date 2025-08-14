package web

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/devfullcycle/imersao22/go-gateway/internal/service"
	"github.com/devfullcycle/imersao22/go-gateway/internal/web/handlers"
	"github.com/devfullcycle/imersao22/go-gateway/internal/web/middleware"
	"github.com/go-chi/chi/v5"
)

// ConfigureRoutes wires HTTP routes using chi mux and provided dependencies.
func ConfigureRoutes(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	// Services
	accountSvc := service.NewAccountService(db)
	invoiceSvc := service.NewInvoiceService(db)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(accountSvc)

	// Handlers
	accountH := handlers.NewAccountHandler(accountSvc)
	invoiceH := handlers.NewInvoiceHandler(invoiceSvc)

	// Routes
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", accountH.PostAccounts()) // POST /accounts
		r.Get("/", accountH.GetAccounts())   // GET /accounts
	})

	// Rotas de invoice COM autenticação
	r.Route("/invoices", func(r chi.Router) {
		// Aplicar auth middleware apenas nas rotas de invoice
		r.Use(authMiddleware.Authenticate)

		r.Post("/", invoiceH.PostInvoices())      // POST /invoices
		r.Get("/", invoiceH.GetInvoices())        // GET /invoices
		r.Get("/{id}", invoiceH.GetInvoiceByID()) // GET /invoices/{id}
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
