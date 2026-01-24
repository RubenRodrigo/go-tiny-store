package app

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/auth"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/category"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/product"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/user"
	authadapter "github.com/RubenRodrigo/go-tiny-store/internal/domain/user/adapters/auth"
	"github.com/RubenRodrigo/go-tiny-store/internal/infraestructure/api"
	"github.com/RubenRodrigo/go-tiny-store/internal/infraestructure/config"
	"github.com/RubenRodrigo/go-tiny-store/internal/infraestructure/db"
	"github.com/RubenRodrigo/go-tiny-store/pkg/consts"
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
	authConfig := consts.NewAuthConfig()
	jwtManager := jwt.NewJWTManager([]byte(a.config.Auth.JWT_SECRET), authConfig)

	// Initialize repositories
	userRepo := user.NewRepository(database)
	categoryRepo := category.NewRepository(database)
	productRepo := product.NewRepository(database)

	// Initialize services
	userService := user.NewService(userRepo)
	categoryService := category.NewService(categoryRepo)
	productService := product.NewService(productRepo)

	// Initialize adapters
	userAuthAdapter := authadapter.NewAdapter(userService)

	// Initialize compose services
	authService := auth.NewService(userAuthAdapter, jwtManager)

	// Initialize services
	services := service.Services{
		User:     userService,
		Auth:     authService,
		Category: categoryService,
		Product:  productService,
	}

	// Initialize REST server
	a.restServer = api.NewServer(services, &a.config.Server, jwtManager)

	return nil
}

func (a *App) Start() error {
	// Start the REST server (which also serves GraphQL)
	return a.restServer.Start()
}
