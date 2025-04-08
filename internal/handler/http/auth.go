package http

import (
	"account-service/internal/domain/billing"
	"account-service/internal/domain/grant"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
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

	r.Get("/pay", h.pay)
	r.Post("/callback", h.callback)
	r.Post("/createPayment", h.createBilling)
	r.Get("/cards", h.getCards)

	return r
}

func (h *Auth) createBilling(w http.ResponseWriter, r *http.Request) {
	req := billing.Entity{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	id, err := h.authService.CreatePayment(r.Context(), req)
	switch {
	case errors.Is(err, grant.ErrUserExist):
		response.BadRequest(w, r, err, nil)
	case errors.Is(err, nil):
		response.OK(w, r, response.Object{Success: true, Data: id})
	default:
		response.InternalServerError(w, r, err)
	}
	return
}

func (h *Auth) getCards(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("userId")
	cards, err := h.authService.GetCards(r.Context(), id)
	switch {
	case errors.Is(err, nil):
		response.OK(w, r, cards)
	default:
		response.InternalServerError(w, r, err)
	}
	return
}

func (h *Auth) pay(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := h.authService.Pay(r.Context(), w, id)
	switch {
	case errors.Is(err, grant.ErrUserExist):
		response.BadRequest(w, r, err, nil)
	case errors.Is(err, nil):
		response.OK(w, r, "")
	default:
		response.InternalServerError(w, r, err)
	}
	return
}

func (h *Auth) callback(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	err = h.authService.Callback(r.Context(), id, body)
	switch {
	case errors.Is(err, grant.ErrUserExist):
		response.BadRequest(w, r, err, nil)
	case errors.Is(err, nil):
		response.OK(w, r, "")
	default:
		response.InternalServerError(w, r, err)
	}
	return
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
