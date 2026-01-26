package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/application/authapp"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/categoryapp"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/productapp"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/userapp"
	"github.com/RubenRodrigo/go-tiny-store/internal/delivery/http/handlers"
	"github.com/RubenRodrigo/go-tiny-store/internal/delivery/http/middleware"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/auth"
	"github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/config"
	"github.com/gorilla/mux"
)

// Services contains all application services
type Services struct {
	User     *userapp.Service
	Auth     *authapp.Service
	Category *categoryapp.Service
	Product  *productapp.Service
}

// Server represents the HTTP server
type Server struct {
	router       *mux.Router
	services     Services
	config       *config.ServerConfig
	tokenService auth.TokenService
}

// NewServer creates a new HTTP server
func NewServer(services Services, cfg *config.ServerConfig, tokenService auth.TokenService) *Server {
	server := &Server{
		router:       mux.NewRouter(),
		services:     services,
		config:       cfg,
		tokenService: tokenService,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Initialize handlers
	h := s.initializeHandlers()

	// Apply global middleware
	s.router.Use(middleware.LoggingMiddleware)
	s.router.Use(middleware.RecoveryMiddleware)

	// Health check
	s.router.HandleFunc("/health", s.healthCheck).Methods("GET")

	// API v1 routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Setup route groups
	s.setupPublicRoutes(api, h)
	s.setupProtectedRoutes(api, h)
	s.setupManagerRoutes(api, h)
}

func (s *Server) initializeHandlers() *handlers.Handlers {
	return handlers.NewHandlers(
		s.services.Auth,
		s.services.User,
		s.services.Category,
		s.services.Product,
	)
}

func (s *Server) setupPublicRoutes(api *mux.Router, h *handlers.Handlers) {
	// Auth routes
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/sign-up", s.handle(h.Auth.SignUp)).Methods("POST")
	auth.HandleFunc("/sign-in", s.handle(h.Auth.SignIn)).Methods("POST")
	auth.HandleFunc("/refresh", s.handle(h.Auth.RefreshToken)).Methods("POST")
	auth.HandleFunc("/forgot-password", s.handle(h.Auth.ForgotPassword)).Methods("POST")
	auth.HandleFunc("/reset-password", s.handle(h.Auth.ResetPassword)).Methods("POST")

	// Webhook routes
	webhooks := api.PathPrefix("/webhooks").Subrouter()
	webhooks.HandleFunc("/stripe", s.handle(h.Webhook.WebhookStripe)).Methods("POST")

	// Public product viewing
	products := api.PathPrefix("/products").Subrouter()
	products.HandleFunc("", s.handle(h.Product.List)).Methods("GET")
	products.HandleFunc("/{id}", s.handle(h.Product.Get)).Methods("GET")
	products.HandleFunc("/category/{categoryId}", s.handle(h.Product.GetByCategory)).Methods("GET")
}

func (s *Server) setupProtectedRoutes(api *mux.Router, h *handlers.Handlers) {
	// Create protected subrouter with auth middleware
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(s.tokenService))

	// User routes
	users := protected.PathPrefix("/users").Subrouter()
	users.HandleFunc("/me", s.handle(h.User.GetCurrentUser)).Methods("GET")
	users.HandleFunc("/me", s.handle(h.User.UpdateProfile)).Methods("PUT")

	// Product interactions
	products := protected.PathPrefix("/products").Subrouter()
	products.HandleFunc("/{id}/like", s.handle(h.Product.Like)).Methods("POST")

	// Cart routes
	cart := protected.PathPrefix("/cart").Subrouter()
	cart.HandleFunc("", s.handle(h.Cart.Get)).Methods("GET")
	cart.HandleFunc("/items", s.handle(h.Cart.AddProduct)).Methods("POST")
	cart.HandleFunc("/items/{productId}", s.handle(h.Cart.RemoveProduct)).Methods("DELETE")
	cart.HandleFunc("", s.handle(h.Cart.Clear)).Methods("DELETE")

	// Order routes
	orders := protected.PathPrefix("/orders").Subrouter()
	orders.HandleFunc("", s.handle(h.Order.ListMyOrders)).Methods("GET")
	orders.HandleFunc("/{id}", s.handle(h.Order.Get)).Methods("GET")

	// Checkout routes
	checkout := protected.PathPrefix("/checkout").Subrouter()
	checkout.HandleFunc("", s.handle(h.Checkout.Create)).Methods("POST")
}

func (s *Server) setupManagerRoutes(api *mux.Router, h *handlers.Handlers) {
	// Create manager subrouter with auth middleware
	manager := api.PathPrefix("/manager").Subrouter()
	manager.Use(middleware.AuthMiddleware(s.tokenService))

	// Product management
	products := manager.PathPrefix("/products").Subrouter()
	products.HandleFunc("", s.handle(h.Product.Create)).Methods("POST")
	products.HandleFunc("/{id}", s.handle(h.Product.Update)).Methods("PUT")
	products.HandleFunc("/{id}", s.handle(h.Product.Delete)).Methods("DELETE")
	products.HandleFunc("/{id}/disable", s.handle(h.Product.Disable)).Methods("PATCH")
	products.HandleFunc("/{id}/images", s.handle(h.Product.UploadImage)).Methods("POST")

	// Category management
	categories := manager.PathPrefix("/categories").Subrouter()
	categories.HandleFunc("", s.handle(h.Category.Create)).Methods("POST")
	categories.HandleFunc("/{id}", s.handle(h.Category.Update)).Methods("PUT")
	categories.HandleFunc("/{id}", s.handle(h.Category.Get)).Methods("GET")
	categories.HandleFunc("", s.handle(h.Category.List)).Methods("GET")

	// Order management
	orders := manager.PathPrefix("/orders").Subrouter()
	orders.HandleFunc("", s.handle(h.Order.ListAllOrders)).Methods("GET")
	orders.HandleFunc("/{id}", s.handle(h.Order.Get)).Methods("GET")

	// User management
	users := manager.PathPrefix("/users").Subrouter()
	users.HandleFunc("", s.handle(h.User.ListUsers)).Methods("GET")
	users.HandleFunc("/{id}", s.handle(h.User.GetUser)).Methods("GET")
}

// Wrapper to handle errors consistently
func (s *Server) handle(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			middleware.HandleError(w, r, err)
		}
	}
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	log.Printf("🚀 REST server starting on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
