package routes

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/handlers/auth_handler"
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/handlers/category_handler"
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/handlers/user_handler"
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/middleware"
	"github.com/RubenRodrigo/go-tiny-store/internal/lib"
	"github.com/gorilla/mux"
)

func SetupRoutes(auth_handler *auth_handler.AuthHandler, user_handler *user_handler.UserHandler, category_handler *category_handler.CategoryHandler, jwtManager *lib.JWTManager) *mux.Router {
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
	protected.Use(middleware.AuthMiddleware(jwtManager))
	protected.HandleFunc("/users/{id}", middleware.WithErrorHandling(user_handler.GetUser, user_handler.ErrorHandler)).Methods("GET")
	protected.HandleFunc("/users", middleware.WithErrorHandling(user_handler.ListUsers, user_handler.ErrorHandler)).Methods("GET")
	protected.HandleFunc("/categories", middleware.WithErrorHandling(category_handler.Save, category_handler.ErrorHandler)).Methods("POST")

	// Add middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecoveryMiddleware)

	return r
}
