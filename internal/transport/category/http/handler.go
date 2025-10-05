package categoryhttp

import (
	"encoding/json"
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/category"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/api/httputil"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/api/middleware"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"github.com/RubenRodrigo/go-tiny-store/pkg/validation"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	categoryService category.Service
	ErrorHandler    middleware.ErrorHandlerFunc
}

func NewCategoryHandler(categoryService category.Service) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		ErrorHandler:    middleware.DefaultErrorHandler,
	}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) error {
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

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) error {
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

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) error {
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

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	id := params["id"]

	err := h.categoryService.Delete(id)

	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, nil)

	return nil
}
