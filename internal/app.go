package app

import (
	"time"

	"github.com/RubenRodrigo/go-tiny-store/internal/application"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/auth"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/category"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/integrations"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/product"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/user"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/repository"
	api "github.com/RubenRodrigo/go-tiny-store/internal/handlers/rest"
	infraAuth "github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/auth"
	"github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/config"
	"github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/db"
	"github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/email"
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
	// Setup db
	db, err := db.SetupDatabase(&a.config.Database)
	if err != nil {
		return err
	}

	// Initialize Token Service (JWT implementation)
	tokenService := infraAuth.NewJWTManager(infraAuth.JWTConfig{
		Secret:          a.config.Auth.JWT_SECRET,
		AccessTokenTTL:  15 * time.Minute,   // 15 minutes
		RefreshTokenTTL: 7 * 24 * time.Hour, // 7 days
		Issuer:          "tiny-store-api",
	})

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	passwordResetTokenRepo := repository.NewPasswordResetTokenRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)

	// Initialize ports
	emailSender := email.NewSendgridManager(email.SendgridConfig{
		APIKey:    "key",
		FromEmail: "admin@admin.com",
	})

	// Initialize integrations
	emailService := integrations.NewEmailService(emailSender)

	// Initialize services
	userService := user.NewService(userRepo)
	categoryService := category.NewService(categoryRepo)
	productService := product.NewService(productRepo)
	authService := auth.NewService(userRepo, passwordResetTokenRepo, tokenService, emailService)

	// Initialize services
	services := application.Services{
		User:     userService,
		Auth:     authService,
		Category: categoryService,
		Product:  productService,
	}

	// Initialize REST server
	a.restServer = api.NewServer(services, &a.config.Server, tokenService)

	return nil
}

func (a *App) Start() error {
	// Start the REST server (which also serves GraphQL)
	return a.restServer.Start()
}
