package handlers

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/application/authapp"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/categoryapp"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/productapp"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/userapp"
	"github.com/RubenRodrigo/go-tiny-store/internal/delivery/http/handlers/auth"
	"github.com/RubenRodrigo/go-tiny-store/internal/delivery/http/handlers/cart"
	"github.com/RubenRodrigo/go-tiny-store/internal/delivery/http/handlers/category"
	"github.com/RubenRodrigo/go-tiny-store/internal/delivery/http/handlers/checkout"
	"github.com/RubenRodrigo/go-tiny-store/internal/delivery/http/handlers/order"
	"github.com/RubenRodrigo/go-tiny-store/internal/delivery/http/handlers/product"
	"github.com/RubenRodrigo/go-tiny-store/internal/delivery/http/handlers/user"
	"github.com/RubenRodrigo/go-tiny-store/internal/delivery/http/handlers/webhook"
)

// Handlers contains all HTTP handlers organized by feature
type Handlers struct {
	Auth     *auth.Handler
	User     *user.Handler
	Category *category.Handler
	Product  *product.Handler
	Order    *order.Handler
	Cart     *cart.Handler
	Checkout *checkout.Handler
	Webhook  *webhook.Handler
}

// NewHandlers creates all handlers with their dependencies
func NewHandlers(
	authService *authapp.Service,
	userService *userapp.Service,
	categoryService *categoryapp.Service,
	productService *productapp.Service,
) *Handlers {
	return &Handlers{
		Auth:     auth.NewHandler(authService),
		User:     user.NewHandler(userService),
		Category: category.NewHandler(categoryService),
		Product:  product.NewHandler(productService),
		Order:    order.NewHandler(),
		Cart:     cart.NewHandler(),
		Checkout: checkout.NewHandler(),
		Webhook:  webhook.NewHandler(),
	}
}
