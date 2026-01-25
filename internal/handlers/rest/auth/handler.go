package authhttp

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/application/auth"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	httputil "github.com/RubenRodrigo/go-tiny-store/pkg/httputils"
	"github.com/RubenRodrigo/go-tiny-store/pkg/validation"
)

type Handler struct {
	authService auth.Service
}

func NewAuthHandler(authService auth.Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) error {
	var req SignUpRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	user, err := h.authService.SignUp(
		req.Email,
		req.Username,
		req.Password,
		req.FirstName,
		req.LastName,
	)

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

	user, err := h.authService.SignIn(
		req.Email,
		req.Password,
	)
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

	// Validate the request
	if validationErrors := validation.Validate(req); len(validationErrors.Errors) > 0 {
		return validationErrors
	}

	err := h.authService.SignOut(
		req.Token,
	)

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

	// Validate the request
	if validationErrors := validation.Validate(req); len(validationErrors.Errors) > 0 {
		return validationErrors
	}

	err := h.authService.ForgotPassword(
		req.Email,
	)

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

	// Validate the request
	if validationErrors := validation.Validate(req); len(validationErrors.Errors) > 0 {
		return validationErrors
	}

	err := h.authService.ResetPassword(
		req.Token,
		req.Password,
	)

	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, nil)

	return nil
}
