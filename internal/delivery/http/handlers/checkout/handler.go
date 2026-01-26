package checkout

import "net/http"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) error {
	return nil
}
