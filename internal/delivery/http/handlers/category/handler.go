package category

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/application/categoryapp"
	httputil "github.com/RubenRodrigo/go-tiny-store/pkg/httputils"
	"github.com/RubenRodrigo/go-tiny-store/pkg/validation"
	"github.com/gorilla/mux"
)

type Handler struct {
	categoryService *categoryapp.Service
}

func NewHandler(categoryService *categoryapp.Service) *Handler {
	return &Handler{
		categoryService: categoryService,
	}
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) error {
	categories, err := h.categoryService.List()
	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, categories)
	return nil
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	id := params["id"]

	category, err := h.categoryService.GetByID(id)
	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, category)
	return nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) error {
	var req CreateCategoryRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	category, err := h.categoryService.Create(req.Name)
	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusCreated, category)
	return nil
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	id := params["id"]

	var req UpdateCategoryRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	category, err := h.categoryService.Update(id, req.Name)
	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, category)
	return nil
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	id := params["id"]

	if err := h.categoryService.Delete(id); err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusNoContent, nil)
	return nil
}
