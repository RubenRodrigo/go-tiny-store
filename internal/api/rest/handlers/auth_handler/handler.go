package auth_handler

import (
	"encoding/json"
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/httputil"
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/middleware"
	"github.com/RubenRodrigo/go-tiny-store/internal/apperrors"
	"github.com/RubenRodrigo/go-tiny-store/internal/service"
)

type AuthHandler struct {
	authService  service.AuthService
	ErrorHandler middleware.ErrorHandlerFunc
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
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

	// Basic validation
	if req.Email == "" || req.Username == "" || req.Password == "" {
		return apperrors.ErrAuthRequiredFields
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

	// Basic validation
	if req.Email == "" || req.Password == "" {
		return apperrors.ErrAuthRequiredFields
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
