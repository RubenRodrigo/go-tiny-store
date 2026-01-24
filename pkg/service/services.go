package service

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/auth"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/category"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/product"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/user"
)

type Services struct {
	User     user.Service
	Auth     auth.Service
	Category *category.Service
	Product  *product.Service
}
