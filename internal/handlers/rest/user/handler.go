package userhttp

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/application/user"
	httputil "github.com/RubenRodrigo/go-tiny-store/pkg/httputils"
	"github.com/gorilla/mux"
)

type Handler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	id := params["id"]

	user, err := h.userService.GetById(id)
	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, user)

	return nil
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.userService.ListUsers()

	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, users)

	return nil
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) error {

	return nil
}
