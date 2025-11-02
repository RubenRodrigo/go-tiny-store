package producthttp

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/product"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/api/middleware"
)

type ProductHandler struct {
	productService product.Service
	ErrorHandler   middleware.ErrorHandlerFunc
}

func NewProductHandler(productService product.Service) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		ErrorHandler:   middleware.DefaultErrorHandler,
	}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *ProductHandler) Disable(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *ProductHandler) Like(w http.ResponseWriter, r *http.Request) error {

	return nil
}
