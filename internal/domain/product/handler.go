package product

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/platform/api/middleware"
)

type ProductHandler struct {
	productService Service
	ErrorHandler   middleware.ErrorHandlerFunc
}

func NewProductHandler(productService Service) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		ErrorHandler:   middleware.DefaultErrorHandler,
	}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) error {

	return nil
}
