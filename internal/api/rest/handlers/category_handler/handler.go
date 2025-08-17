package category_handler

import (
	"encoding/json"
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/httputil"
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/middleware"
	"github.com/RubenRodrigo/go-tiny-store/internal/apperrors"
	"github.com/RubenRodrigo/go-tiny-store/internal/lib"
	"github.com/RubenRodrigo/go-tiny-store/internal/service"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	categoryService service.CategoryService
	ErrorHandler    middleware.ErrorHandlerFunc
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		ErrorHandler:    middleware.DefaultErrorHandler,
	}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var req CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.ErrRequestInvalidBody
	}

	// Validate the request
	if validationErrors := lib.Validate(req); len(validationErrors.Errors) > 0 {
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
	var req CategoryRequest

	params := mux.Vars(r)
	id := params["id"]

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.ErrRequestInvalidBody
	}

	// Validate the request
	if validationErrors := lib.Validate(req); len(validationErrors.Errors) > 0 {
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
