package api

import (
	"fmt"
	"log"

	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/platform/api/middleware"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/config"
	authhttp "github.com/RubenRodrigo/go-tiny-store/internal/transport/auth/http"
	categoryhttp "github.com/RubenRodrigo/go-tiny-store/internal/transport/category/http"
	producthttp "github.com/RubenRodrigo/go-tiny-store/internal/transport/product/http"
	userhttp "github.com/RubenRodrigo/go-tiny-store/internal/transport/user/http"
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

	r := mux.NewRouter()

	// API subrouter
	api := r.PathPrefix("/api").Subrouter()

	// Status route
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	// Public routes
	api.HandleFunc("/auth/register", middleware.WithErrorHandling(auth_handler.RegisterUser, auth_handler.ErrorHandler)).Methods("POST")
	api.HandleFunc("/auth/login", middleware.WithErrorHandling(auth_handler.LoginUser, auth_handler.ErrorHandler)).Methods("POST")
	api.HandleFunc("/auth/logout", middleware.WithErrorHandling(auth_handler.LogOutUser, auth_handler.ErrorHandler)).Methods("POST")
	api.HandleFunc("/categories", middleware.WithErrorHandling(category_handler.List, category_handler.ErrorHandler)).Methods("GET")

	// Protected routes
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware(s.jwtManager))
	protected.HandleFunc("/users/{id}", middleware.WithErrorHandling(user_handler.GetUser, user_handler.ErrorHandler)).Methods("GET")
	protected.HandleFunc("/users", middleware.WithErrorHandling(user_handler.ListUsers, user_handler.ErrorHandler)).Methods("GET")
	protected.HandleFunc("/categories", middleware.WithErrorHandling(category_handler.Create, category_handler.ErrorHandler)).Methods("POST")
	protected.HandleFunc("/categories/{id}", middleware.WithErrorHandling(category_handler.Update, category_handler.ErrorHandler)).Methods("PUT")
	protected.HandleFunc("/products", middleware.WithErrorHandling(product_handler.Create, product_handler.ErrorHandler)).Methods("POST")

	// Add middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecoveryMiddleware)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	log.Printf("REST server starting on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
