package userhttp

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/user"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/api/httputil"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/api/middleware"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService  user.Service
	ErrorHandler middleware.ErrorHandlerFunc
}

func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{
		userService:  userService,
		ErrorHandler: middleware.DefaultErrorHandler,
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	id := params["id"]

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, user)

	return nil
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.userService.ListUsers()

	if err != nil {
		return err
	}

	httputil.RespondWithJSON(w, http.StatusOK, users)

	return nil
}
