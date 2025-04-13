package user_handler

import (
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/httputil"
	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/middleware"
	"github.com/RubenRodrigo/go-tiny-store/internal/service"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService  service.UserService
	ErrorHandler middleware.ErrorHandlerFunc
}

func NewUserHandler(userService service.UserService) *UserHandler {
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
