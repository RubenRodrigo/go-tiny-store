package order

import "net/http"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ListMyOrders(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *Handler) ListAllOrders(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) error {
	return nil
}
