package authhttp

import (
	"encoding/json"
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/auth"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/api/httputil"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/api/middleware"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"github.com/RubenRodrigo/go-tiny-store/pkg/validation"
)

type AuthHandler struct {
	authService  auth.Service
	ErrorHandler middleware.ErrorHandlerFunc
}

func NewAuthHandler(authService auth.Service) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		ErrorHandler: middleware.DefaultErrorHandler,
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) error {
	var req RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.ErrRequestInvalidBody
	}

	// Validate the request
	if validationErrors := validation.Validate(req); len(validationErrors.Errors) > 0 {
		return validationErrors
	}

	user, err := h.authService.RegisterUser(
		req.Email,
		req.Username,
		req.Password,
		req.FirstName,
		req.LastName,
	)

	if err != nil {
		return err
	}

	resp := RegisterUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	httputil.RespondWithJSON(w, http.StatusOK, resp)

	return nil
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) error {
	var req LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.ErrRequestInvalidBody
	}

	// Validate the request
	if validationErrors := validation.Validate(req); len(validationErrors.Errors) > 0 {
		return validationErrors
	}

	user, token, err := h.authService.LoginUser(
		req.Email,
		req.Password,
	)
	if err != nil {
		return err
	}

	resp := LoginUserResponse{
		Token:     token,
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	httputil.RespondWithJSON(w, http.StatusOK, resp)

	return nil
}

// TODO: Pending to implement
func (h *AuthHandler) LogOutUser(w http.ResponseWriter, r *http.Request) error {
	var req LogOutUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.ErrRequestInvalidBody
	}

	// Validate the request
	if validationErrors := validation.Validate(req); len(validationErrors.Errors) > 0 {
		return validationErrors
	}

	err := h.authService.LogOutUser(
		req.Token,
	)

	if err != nil {
		return apperrors.ErrAuthTokenInvalid
	}

	httputil.RespondWithJSON(w, http.StatusOK, nil)

	return nil
}
