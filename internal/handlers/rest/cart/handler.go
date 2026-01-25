package carthttp

import (
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *Handler) AddProduct(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *Handler) RemoveProduct(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *Handler) Clear(w http.ResponseWriter, r *http.Request) error {

	return nil
}
