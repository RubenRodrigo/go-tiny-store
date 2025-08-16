package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/handlers/auth_handler"
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/handlers/category_handler"
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/handlers/user_handler"
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/routes"
	"github.com/RubenRodrigo/go-tiny-store/internal/config"
	"github.com/RubenRodrigo/go-tiny-store/internal/lib"
	"github.com/RubenRodrigo/go-tiny-store/internal/service"
	"github.com/gorilla/mux"
)

type Server struct {
	router          *mux.Router
	userService     service.UserService
	authService     service.AuthService
	categoryService service.CategoryService
	config          *config.ServerConfig
	jwtManager      *lib.JWTManager
}

func NewServer(userService service.UserService, authService service.AuthService, categoryService service.CategoryService, cfg *config.ServerConfig, jwtManager *lib.JWTManager) *Server {
	server := &Server{
		router:          mux.NewRouter(),
		userService:     userService,
		authService:     authService,
		categoryService: categoryService,
		config:          cfg,
		jwtManager:      jwtManager,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	auth_handler := auth_handler.NewAuthHandler(s.authService)
	user_handler := user_handler.NewUserHandler(s.userService)
	category_handler := category_handler.NewCategoryHandler(s.categoryService)

	s.router = routes.SetupRoutes(auth_handler, user_handler, category_handler, s.jwtManager)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	log.Printf("REST server starting on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
