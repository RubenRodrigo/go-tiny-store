package auth

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/application/authapp"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	httputil "github.com/RubenRodrigo/go-tiny-store/pkg/httputils"
	"github.com/RubenRodrigo/go-tiny-store/pkg/validation"
)

// Handler handles authentication HTTP requests
type Handler struct {
	authService *authapp.Service
}

// NewHandler creates a new auth handler
func NewHandler(authService *authapp.Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) error {
	var req SignUpRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	// Convert HTTP request to application DTO
	dto := authapp.SignUpDTO{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	user, err := h.authService.SignUp(dto)
	if err != nil {
		return err
	}

	resp := AuthUserResponse{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}

	httputil.RespondWithJSON(w, http.StatusOK, resp)
	return nil
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) error {
	var req SignInRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	// Convert HTTP request to application DTO
	dto := authapp.SignInDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.authService.SignIn(dto)
	if err != nil {
		return err
	}

	resp := AuthUserResponse{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}

	httputil.RespondWithJSON(w, http.StatusOK, resp)
	return nil
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	var req RefreshTokenRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	user, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return err
	}

	resp := AuthUserResponse{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}

	httputil.RespondWithJSON(w, http.StatusOK, resp)
	return nil
}

func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) error {
	var req SignOutRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	err := h.authService.SignOut(req.Token)
	if err != nil {
		return apperrors.ErrAuthTokenInvalid
	}

	httputil.RespondWithJSON(w, http.StatusOK, nil)
	return nil
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) error {
	var req ForgotPasswordRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	err := h.authService.ForgotPassword(req.Email)
	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, nil)
	return nil
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	var req ResetPasswordRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	err := h.authService.ResetPassword(req.Token, req.Password)
	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, nil)
	return nil
}
