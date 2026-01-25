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

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) error {
	var req RegisterUserRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		return err
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

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) error {
	var req LoginUserRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	user, accessToken, refreshToken, err := h.authService.LoginUser(
		req.Email,
		req.Password,
	)
	if err != nil {
		return err
	}

	resp := AuthUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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

	user, accessToken, refreshToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return err
	}

	resp := AuthUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}

	httputil.RespondWithJSON(w, http.StatusOK, resp)

	return nil
}

func (h *Handler) LogOutUser(w http.ResponseWriter, r *http.Request) error {
	var req LogOutUserRequest
	if err := validation.DecodeAndValidate(r, &req); err != nil {
		return err
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

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	return nil
}
