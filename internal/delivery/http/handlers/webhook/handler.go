package webhook

import "net/http"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) WebhookStripe(w http.ResponseWriter, r *http.Request) error {
	return nil
}
