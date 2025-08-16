package app

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest"
	"github.com/RubenRodrigo/go-tiny-store/internal/config"
	"github.com/RubenRodrigo/go-tiny-store/internal/db"
	"github.com/RubenRodrigo/go-tiny-store/internal/lib"
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

	// Initialize Create JWT manager
	jwtManager := lib.NewJWTManager([]byte(a.config.Auth.JWT_SECRET))

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)
	categoryRepo := repository.NewCategoryRepository(database)

	// Initialize services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, jwtManager)
	categoryService := service.NewCategoryService(categoryRepo)

	// Initialize REST server
	a.restServer = rest.NewServer(userService, authService, categoryService, &a.config.Server, jwtManager)

	return nil
}

func (a *App) Start() error {
	// Start the REST server (which also serves GraphQL)
	return a.restServer.Start()
}
