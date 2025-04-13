package app

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest"
	"github.com/RubenRodrigo/go-tiny-store/internal/config"
	"github.com/RubenRodrigo/go-tiny-store/internal/db"
	"github.com/RubenRodrigo/go-tiny-store/internal/repository"
	"github.com/RubenRodrigo/go-tiny-store/internal/service"
)

type App struct {
	config     *config.Config
	restServer *rest.Server
}

func New() *App {
	cfg := config.Load()

	return &App{
		config: cfg,
	}
}

func (a *App) Initialize() error {
	// Setup database
	database, err := db.SetupDatabase(&a.config.Database)
	if err != nil {
		return err
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)

	// Initialize services
	userService := service.NewUserService(userRepo, "JWT_SECRET")
	authService := service.NewAuthService(userRepo, "JWT_SECRET")

	// Initialize REST server
	a.restServer = rest.NewServer(userService, authService, &a.config.Server)

	return nil
}

func (a *App) Start() error {
	// Start the REST server (which also serves GraphQL)
	return a.restServer.Start()
}
