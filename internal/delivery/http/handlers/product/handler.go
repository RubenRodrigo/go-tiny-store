package product

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/application/productapp"
	httputil "github.com/RubenRodrigo/go-tiny-store/pkg/httputils"
	"github.com/RubenRodrigo/go-tiny-store/pkg/pagination"
	"github.com/gorilla/mux"
)

// Handler handles product HTTP requests
type Handler struct {
	productService *productapp.Service
}

// NewHandler creates a new product handler
func NewHandler(productService *productapp.Service) *Handler {
	return &Handler{
		productService: productService,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) error {
	// TODO: Implement product creation
	return nil
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) error {
	paginationParams := pagination.ParseParams(r)
	filters := ParseFilters(r)

	result, err := h.productService.List(paginationParams, filters)
	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, result)
	return nil
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	id := params["id"]

	product, err := h.productService.Get(id)
	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, product)
	return nil
}

func (h *Handler) GetByCategory(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	categoryID := params["categoryId"]

	paginationParams := pagination.ParseParams(r)
	filters := ParseFilters(r)
	filters.CategoryID = categoryID

	result, err := h.productService.List(paginationParams, filters)
	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, result)
	return nil
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) error {
	// TODO: Implement product update
	return nil
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) error {
	// TODO: Implement product deletion
	return nil
}

func (h *Handler) Disable(w http.ResponseWriter, r *http.Request) error {
	// TODO: Implement product disable
	return nil
}

func (h *Handler) Like(w http.ResponseWriter, r *http.Request) error {
	// TODO: Implement product like
	return nil
}

func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) error {
	// TODO: Implement image upload
	return nil
}
