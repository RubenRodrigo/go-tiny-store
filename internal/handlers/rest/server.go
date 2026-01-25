package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/application"
	authhttp "github.com/RubenRodrigo/go-tiny-store/internal/handlers/rest/auth"
	carthttp "github.com/RubenRodrigo/go-tiny-store/internal/handlers/rest/cart"
	categoryhttp "github.com/RubenRodrigo/go-tiny-store/internal/handlers/rest/category"
	checkouthttp "github.com/RubenRodrigo/go-tiny-store/internal/handlers/rest/checkout"
	"github.com/RubenRodrigo/go-tiny-store/internal/handlers/rest/middleware"
	orderhttp "github.com/RubenRodrigo/go-tiny-store/internal/handlers/rest/order"
	producthttp "github.com/RubenRodrigo/go-tiny-store/internal/handlers/rest/product"
	userhttp "github.com/RubenRodrigo/go-tiny-store/internal/handlers/rest/user"
	"github.com/RubenRodrigo/go-tiny-store/internal/handlers/rest/webhook"
	infraAuth "github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/auth"
	"github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/config"
	"github.com/gorilla/mux"
)

type Server struct {
	router       *mux.Router
	services     application.Services
	config       *config.ServerConfig
	tokenService infraAuth.TokenService
}

func NewServer(services application.Services, cfg *config.ServerConfig, tokenService infraAuth.TokenService) *Server {
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
	handlers := s.initializeHandlers()

	// Apply global middleware
	s.router.Use(middleware.LoggingMiddleware)
	s.router.Use(middleware.RecoveryMiddleware)
	// s.router.Use(middleware.CORS())

	// Health check
	s.router.HandleFunc("/health", s.healthCheck).Methods("GET")

	// API v1 routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Setup route groups
	s.setupPublicRoutes(api, handlers)
	s.setupProtectedRoutes(api, handlers)
	s.setupManagerRoutes(api, handlers)
}

// Handler collection
type Handlers struct {
	Auth     *authhttp.Handler
	User     *userhttp.Handler
	Category *categoryhttp.Handler
	Product  *producthttp.Handler
	Order    *orderhttp.Handler
	Cart     *carthttp.Handler
	Checkout *checkouthttp.Handler
	Webhook  *webhook.Handler
}

func (s *Server) initializeHandlers() *Handlers {
	return &Handlers{
		Auth:     authhttp.NewAuthHandler(*s.services.Auth),
		User:     userhttp.NewUserHandler(*s.services.User),
		Category: categoryhttp.NewCategoryHandler(*s.services.Category),
		Product:  producthttp.NewProductHandler(*s.services.Product),
		Order:    orderhttp.NewOrderHandler(),
		Cart:     carthttp.NewHandler(),
		Checkout: checkouthttp.NewHandler(),
		Webhook:  webhook.NewHandler(),
	}
}

func (s *Server) setupPublicRoutes(api *mux.Router, h *Handlers) {
	// Auth routes
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", s.handle(h.Auth.RegisterUser)).Methods("POST")
	auth.HandleFunc("/login", s.handle(h.Auth.LoginUser)).Methods("POST")
	auth.HandleFunc("/refresh", s.handle(h.Auth.RefreshToken)).Methods("POST")
	auth.HandleFunc("/forgot-password", s.handle(h.Auth.ForgotPassword)).Methods("POST")
	auth.HandleFunc("/reset-password", s.handle(h.Auth.ResetPassword)).Methods("POST")

	// Webhook routes (public but should verify Stripe signature)
	webhooks := api.PathPrefix("/webhooks").Subrouter()
	webhooks.HandleFunc("/stripe", s.handle(h.Webhook.WebhookStripe)).Methods("POST")

	// Public product viewing (no auth required per your requirements)
	products := api.PathPrefix("/products").Subrouter()
	products.HandleFunc("", s.handle(h.Product.List)).Methods("GET")
	products.HandleFunc("/{id}", s.handle(h.Product.Get)).Methods("GET")
	products.HandleFunc("/category/{categoryId}", s.handle(h.Product.GetByCategory)).Methods("GET")
}

func (s *Server) setupProtectedRoutes(api *mux.Router, h *Handlers) {
	// Create protected subrouter with auth middleware
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(s.tokenService))

	// User routes (client operations)
	users := protected.PathPrefix("/users").Subrouter()
	users.HandleFunc("/me", s.handle(h.User.GetCurrentUser)).Methods("GET")
	users.HandleFunc("/me", s.handle(h.User.UpdateProfile)).Methods("PUT")

	// Product interactions (client operations)
	products := protected.PathPrefix("/products").Subrouter()
	products.HandleFunc("/{id}/like", s.handle(h.Product.Like)).Methods("POST")

	// Cart routes (client operations)
	cart := protected.PathPrefix("/cart").Subrouter()
	cart.HandleFunc("", s.handle(h.Cart.Get)).Methods("GET")
	cart.HandleFunc("/items", s.handle(h.Cart.AddProduct)).Methods("POST")
	cart.HandleFunc("/items/{productId}", s.handle(h.Cart.RemoveProduct)).Methods("DELETE")
	cart.HandleFunc("", s.handle(h.Cart.Clear)).Methods("DELETE")

	// Order routes (client operations)
	orders := protected.PathPrefix("/orders").Subrouter()
	orders.HandleFunc("", s.handle(h.Order.ListMyOrders)).Methods("GET")
	orders.HandleFunc("/{id}", s.handle(h.Order.Get)).Methods("GET")

	// Checkout routes
	checkout := protected.PathPrefix("/checkout").Subrouter()
	checkout.HandleFunc("", s.handle(h.Checkout.Create)).Methods("POST")
}

func (s *Server) setupManagerRoutes(api *mux.Router, h *Handlers) {
	// Create manager subrouter with auth + role middleware
	manager := api.PathPrefix("/manager").Subrouter()
	manager.Use(middleware.AuthMiddleware(s.tokenService))
	// manager.Use(middleware.RequireRole("manager"))

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

	// Order management (view all orders)
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

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	log.Printf("🚀 REST server starting on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
