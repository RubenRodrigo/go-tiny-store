package producthttp

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/application/product"
	httputil "github.com/RubenRodrigo/go-tiny-store/pkg/httputils"
	"github.com/RubenRodrigo/go-tiny-store/pkg/pagination"
	"github.com/gorilla/mux"
)

type Handler struct {
	productService product.Service
}

func NewProductHandler(productService product.Service) *Handler {
	return &Handler{
		productService: productService,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) error {
	paginationParams := pagination.ParseParams(r)
	// filterParams := ParseParams(r)

	result, err := h.productService.List(paginationParams)
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

	return nil
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *Handler) Disable(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *Handler) Like(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) error {

	return nil
}
