package routes

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/handlers/auth_handler"
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/handlers/user_handler"
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/middleware"
	"github.com/gorilla/mux"
)

func SetupRoutes(auth_handler *auth_handler.AuthHandler, user_handler *user_handler.UserHandler) *mux.Router {

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
	api.HandleFunc("/users/{id}", middleware.WithErrorHandling(user_handler.GetUser, user_handler.ErrorHandler)).Methods("GET")
	api.HandleFunc("/users", middleware.WithErrorHandling(user_handler.ListUsers, user_handler.ErrorHandler)).Methods("GET")

	// Add middleware
	r.Use(loggingMiddleware)
	r.Use(recoveryMiddleware)

	return r
}

// Middleware implementations
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request info
		next.ServeHTTP(w, r)
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
