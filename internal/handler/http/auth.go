package http

import (
	"account-service/internal/domain/grant"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"account-service/internal/service/auth"
	"account-service/pkg/server/response"
)

type Auth struct {
	authService *auth.Service
}

func NewAuthHandler(authService *auth.Service) *Auth {
	return &Auth{
		authService: authService,
	}
}

func (h *Auth) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/sign-up", h.signUp)
	r.Post("/sign-in", h.signIn)

	return r
}

func (h *Auth) signUp(w http.ResponseWriter, r *http.Request) {
	req := grant.Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	accessToken, err := h.authService.SignUp(r.Context(), req)
	switch {
	case errors.Is(err, grant.ErrUserExist):
		response.BadRequest(w, r, err, nil)
	case errors.Is(err, nil):
		response.OK(w, r, accessToken)
	default:
		response.InternalServerError(w, r, err)
	}
	return
}

func (h *Auth) signIn(w http.ResponseWriter, r *http.Request) {
	req := grant.Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	accessToken, err := h.authService.SignIn(r.Context(), req)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		response.BadRequest(w, r, errors.New("Не найден пользователь"), nil)
	case errors.Is(err, nil):
		response.OK(w, r, accessToken)
		return
	default:
		response.InternalServerError(w, r, err)
		return
	}
	return
}
