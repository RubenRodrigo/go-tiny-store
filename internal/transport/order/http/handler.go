package orderhttp

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/infraestructure/api/middleware"
)

type OrderHandler struct {
	ErrorHandler middleware.ErrorHandlerFunc
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		ErrorHandler: middleware.DefaultErrorHandler,
	}
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) error {

	return nil
}
