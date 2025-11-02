package api

import (
	"fmt"
	"log"

	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/platform/api/middleware"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/config"
	authhttp "github.com/RubenRodrigo/go-tiny-store/internal/transport/auth/http"
	carthttp "github.com/RubenRodrigo/go-tiny-store/internal/transport/cart/http"
	categoryhttp "github.com/RubenRodrigo/go-tiny-store/internal/transport/category/http"
	checkouthttp "github.com/RubenRodrigo/go-tiny-store/internal/transport/checkout/http"
	orderhttp "github.com/RubenRodrigo/go-tiny-store/internal/transport/order/http"
	producthttp "github.com/RubenRodrigo/go-tiny-store/internal/transport/product/http"
	userhttp "github.com/RubenRodrigo/go-tiny-store/internal/transport/user/http"
	webhookhttp "github.com/RubenRodrigo/go-tiny-store/internal/transport/webhook/http"
	"github.com/RubenRodrigo/go-tiny-store/pkg/jwt"
	"github.com/RubenRodrigo/go-tiny-store/pkg/service"
	"github.com/gorilla/mux"
)

type Server struct {
	router     *mux.Router
	services   service.Services
	config     *config.ServerConfig
	jwtManager *jwt.JWTManager
}

func NewServer(services service.Services, cfg *config.ServerConfig, jwtManager *jwt.JWTManager) *Server {
	server := &Server{
		router:     mux.NewRouter(),
		services:   services,
		config:     cfg,
		jwtManager: jwtManager,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	auth_handler := authhttp.NewAuthHandler(*s.services.Auth)
	user_handler := userhttp.NewUserHandler(*s.services.User)
	category_handler := categoryhttp.NewCategoryHandler(*s.services.Category)
	product_handler := producthttp.NewProductHandler(*s.services.Product)
	order_handler := orderhttp.NewOrderHandler()
	cart_handler := carthttp.NewHandler()
	checkout_handler := checkouthttp.NewHandler()
	webhook_handler := webhookhttp.NewHandler()

	r := mux.NewRouter()

	// API subrouter
	api := r.PathPrefix("/api").Subrouter()

	// Status route
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	// Public routes
	// Auth
	api.HandleFunc("/auth/register", middleware.WithErrorHandling(auth_handler.RegisterUser, auth_handler.ErrorHandler)).Methods("POST")
	api.HandleFunc("/auth/login", middleware.WithErrorHandling(auth_handler.LoginUser, auth_handler.ErrorHandler)).Methods("POST")
	api.HandleFunc("/auth/forgot", middleware.WithErrorHandling(auth_handler.LogOutUser, auth_handler.ErrorHandler)).Methods("POST")
	api.HandleFunc("/auth/restore", middleware.WithErrorHandling(auth_handler.LogOutUser, auth_handler.ErrorHandler)).Methods("POST")
	// Webhooks
	api.HandleFunc("/webhooks/stripe", middleware.WithErrorHandling(webhook_handler.WebhookStripe, webhook_handler.ErrorHandler)).Methods("POST")

	// Protected routes
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware(s.jwtManager))
	// Products
	protected.HandleFunc("/products", middleware.WithErrorHandling(product_handler.List, product_handler.ErrorHandler)).Methods("GET")
	protected.HandleFunc("/products/{id}", middleware.WithErrorHandling(product_handler.Get, product_handler.ErrorHandler)).Methods("GET")
	protected.HandleFunc("/products", middleware.WithErrorHandling(product_handler.Create, product_handler.ErrorHandler)).Methods("POST")
	protected.HandleFunc("/products/{id}", middleware.WithErrorHandling(product_handler.Update, product_handler.ErrorHandler)).Methods("PUT")
	protected.HandleFunc("/products/disable", middleware.WithErrorHandling(product_handler.Disable, product_handler.ErrorHandler)).Methods("POST")
	protected.HandleFunc("/products/{id}/like", middleware.WithErrorHandling(product_handler.Like, product_handler.ErrorHandler)).Methods("POST")
	// Orders
	protected.HandleFunc("/orders", middleware.WithErrorHandling(order_handler.List, order_handler.ErrorHandler)).Methods("GET")
	protected.HandleFunc("/orders/{id}", middleware.WithErrorHandling(order_handler.Get, order_handler.ErrorHandler)).Methods("GET")
	// Cart
	protected.HandleFunc("/cart/{id}", middleware.WithErrorHandling(cart_handler.Get, cart_handler.ErrorHandler)).Methods("GET")
	protected.HandleFunc("/cart", middleware.WithErrorHandling(cart_handler.Create, cart_handler.ErrorHandler)).Methods("POST")
	protected.HandleFunc("/cart/add", middleware.WithErrorHandling(cart_handler.AddProduct, cart_handler.ErrorHandler)).Methods("POST")
	// Checkout
	protected.HandleFunc("/checkout/create", middleware.WithErrorHandling(checkout_handler.Create, checkout_handler.ErrorHandler)).Methods("POST")
	// Users
	protected.HandleFunc("/users/{id}", middleware.WithErrorHandling(user_handler.GetUser, user_handler.ErrorHandler)).Methods("GET")
	protected.HandleFunc("/users", middleware.WithErrorHandling(user_handler.ListUsers, user_handler.ErrorHandler)).Methods("GET")
	// Categories
	protected.HandleFunc("/categories", middleware.WithErrorHandling(category_handler.Create, category_handler.ErrorHandler)).Methods("POST")
	protected.HandleFunc("/categories/{id}", middleware.WithErrorHandling(category_handler.Update, category_handler.ErrorHandler)).Methods("PUT")
	protected.HandleFunc("/categories/{id}", middleware.WithErrorHandling(category_handler.Get, category_handler.ErrorHandler)).Methods("GET")

	// Add middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecoveryMiddleware)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	log.Printf("REST server starting on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
