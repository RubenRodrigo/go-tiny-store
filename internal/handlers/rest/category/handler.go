package categoryhttp

import (
	"encoding/json"
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/application/category"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	httputil "github.com/RubenRodrigo/go-tiny-store/pkg/httputils"
	"github.com/RubenRodrigo/go-tiny-store/pkg/validation"
	"github.com/gorilla/mux"
)

type Handler struct {
	categoryService category.Service
}

func NewCategoryHandler(categoryService category.Service) *Handler {
	return &Handler{
		categoryService: categoryService,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) error {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.ErrRequestInvalidBody
	}

	// Validate the request
	if validationErrors := validation.Validate(req); len(validationErrors.Errors) > 0 {
		return validationErrors
	}

	category, err := h.categoryService.Create(req.Name)
	if err != nil {
		return err
	}

	resp := CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	httputil.RespondWithJSON(w, http.StatusOK, resp)

	return nil
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) error {
	var req UpdateCategoryRequest

	params := mux.Vars(r)
	id := params["id"]

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.ErrRequestInvalidBody
	}

	// Validate the request
	if validationErrors := validation.Validate(req); len(validationErrors.Errors) > 0 {
		return validationErrors
	}

	category, err := h.categoryService.Update(req.Name, id)
	if err != nil {
		return err
	}

	resp := CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	httputil.RespondWithJSON(w, http.StatusOK, resp)

	return nil
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) error {
	categories, err := h.categoryService.List()

	if err != nil {
		return err
	}

	resp := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		resp[i] = CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		}
	}

	httputil.RespondWithJSON(w, http.StatusOK, resp)

	return nil
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) error {
	categories, err := h.categoryService.List()

	if err != nil {
		return err
	}

	resp := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		resp[i] = CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		}
	}

	httputil.RespondWithJSON(w, http.StatusOK, resp)

	return nil
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	id := params["id"]

	err := h.categoryService.Delete(id)

	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, nil)

	return nil
}
