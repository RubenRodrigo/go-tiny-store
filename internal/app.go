package app

import (
	"time"

	"github.com/RubenRodrigo/go-tiny-store/internal/adapters/email"
	gormadapter "github.com/RubenRodrigo/go-tiny-store/internal/adapters/persistence/gorm"
	"github.com/RubenRodrigo/go-tiny-store/internal/adapters/security"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/authapp"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/categoryapp"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/productapp"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/userapp"
	"github.com/RubenRodrigo/go-tiny-store/internal/delivery/http"
	"github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/config"
	"github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/database"
)

// App represents the application
type App struct {
	config     *config.Config
	restServer *http.Server
}

// New creates a new application instance
func New() *App {
	cfg := config.Load()

	return &App{
		config: cfg,
	}
}

// Initialize sets up all dependencies
func (a *App) Initialize() error {
	// Setup database connection
	db, err := database.SetupDatabase(&a.config.Database)
	if err != nil {
		return err
	}

	// Initialize adapters (implementations of ports)
	// Security adapters
	passwordHasher := security.NewBcryptHasher()
	tokenHasher := security.NewSHA256TokenHasher()
	tokenService := security.NewJWTService(security.JWTConfig{
		Secret:          a.config.Auth.JWT_SECRET,
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "tiny-store-api",
	})

	// Email adapter
	emailSender := email.NewSendgridSender(email.SendgridConfig{
		APIKey:    "key",
		FromEmail: "admin@admin.com",
	})

	// Repository adapters (GORM implementations)
	userRepo := gormadapter.NewUserRepository(db)
	refreshTokenRepo := gormadapter.NewRefreshTokenRepository(db)
	passwordResetTokenRepo := gormadapter.NewPasswordResetTokenRepository(db)
	categoryRepo := gormadapter.NewCategoryRepository(db)
	productRepo := gormadapter.NewProductRepository(db)

	// Initialize application services (use cases)
	userService := userapp.NewService(userRepo)
	categoryService := categoryapp.NewService(categoryRepo)
	productService := productapp.NewService(productRepo)
	authService := authapp.NewService(
		userRepo,
		refreshTokenRepo,
		passwordResetTokenRepo,
		tokenService,
		passwordHasher,
		tokenHasher,
		emailSender,
	)

	// Create services container
	services := http.Services{
		User:     userService,
		Auth:     authService,
		Category: categoryService,
		Product:  productService,
	}

	// Initialize HTTP server (delivery layer)
	a.restServer = http.NewServer(services, &a.config.Server, tokenService)

	return nil
}

// Start starts the application
func (a *App) Start() error {
	return a.restServer.Start()
}
