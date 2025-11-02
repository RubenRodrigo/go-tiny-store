package carthttp

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/platform/api/middleware"
)

type CartHandler struct {
	ErrorHandler middleware.ErrorHandlerFunc
}

func NewHandler() *CartHandler {
	return &CartHandler{
		ErrorHandler: middleware.DefaultErrorHandler,
	}
}

func (h *CartHandler) Get(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *CartHandler) Create(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *CartHandler) AddProduct(w http.ResponseWriter, r *http.Request) error {

	return nil
}
