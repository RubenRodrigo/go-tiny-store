package app

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/auth"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/category"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/product"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/user"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/api"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/config"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/db"
	"github.com/RubenRodrigo/go-tiny-store/pkg/jwt"
	"github.com/RubenRodrigo/go-tiny-store/pkg/service"
)

type App struct {
	config     *config.Config
	restServer *api.Server
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
	jwtManager := jwt.NewJWTManager([]byte(a.config.Auth.JWT_SECRET))

	// Initialize repositories
	userRepo := user.NewRepository(database)
	categoryRepo := category.NewRepository(database)
	productRepo := product.NewRepository(database)

	// Initialize services
	services := service.Services{
		User:     user.NewService(userRepo),
		Auth:     auth.NewService(userRepo, jwtManager),
		Category: category.NewService(categoryRepo),
		Product:  product.NewService(productRepo),
	}

	// Initialize REST server
	a.restServer = api.NewServer(services, &a.config.Server, jwtManager)

	return nil
}

func (a *App) Start() error {
	// Start the REST server (which also serves GraphQL)
	return a.restServer.Start()
}
