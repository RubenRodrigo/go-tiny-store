package application

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/application/auth"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/category"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/product"
	"github.com/RubenRodrigo/go-tiny-store/internal/application/user"
)

type Services struct {
	User     *user.Service
	Auth     *auth.Service
	Category *category.Service
	Product  *product.Service
}
