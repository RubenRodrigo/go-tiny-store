package webhookhttp

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/infraestructure/api/middleware"
)

type WebhookHandler struct {
	ErrorHandler middleware.ErrorHandlerFunc
}

func NewHandler() *WebhookHandler {
	return &WebhookHandler{}
}

func (h *WebhookHandler) WebhookStripe(w http.ResponseWriter, r *http.Request) error { return nil }
