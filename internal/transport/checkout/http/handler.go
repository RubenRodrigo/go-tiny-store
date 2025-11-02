package checkouthttp

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/platform/api/middleware"
)

type CheckoutHandler struct {
	ErrorHandler middleware.ErrorHandlerFunc
}

func NewHandler() *CheckoutHandler {
	return &CheckoutHandler{
		ErrorHandler: middleware.DefaultErrorHandler,
	}
}

func (h *CheckoutHandler) Create(w http.ResponseWriter, r *http.Request) error {
	return nil
}
