package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/handlers/auth_handler"
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/handlers/user_handler"
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/routes"
	"github.com/RubenRodrigo/go-tiny-store/internal/config"
	"github.com/RubenRodrigo/go-tiny-store/internal/service"
	"github.com/gorilla/mux"
)

type Server struct {
	router      *mux.Router
	userService service.UserService
	authService service.AuthService
	config      *config.ServerConfig
}

func NewServer(userService service.UserService, authService service.AuthService, cfg *config.ServerConfig) *Server {
	server := &Server{
		router:      mux.NewRouter(),
		userService: userService,
		authService: authService,
		config:      cfg,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	auth_handler := auth_handler.NewAuthHandler(s.authService)
	user_handler := user_handler.NewUserHandler(s.userService)
	s.router = routes.SetupRoutes(auth_handler, user_handler)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	log.Printf("REST server starting on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
