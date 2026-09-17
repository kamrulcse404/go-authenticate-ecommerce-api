package auth

import (
	"ecommerce-api/internal/platform/response"
	"encoding/json"
	"errors"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Register(r.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, ErrNameRequired), errors.Is(err, ErrEmailRequired), errors.Is(err, ErrPasswordRequired), errors.Is(err, ErrPasswordTooShort):
			response.Error(w, http.StatusBadRequest, err.Error())

		case errors.Is(err, ErrEmailAlreadyExists):
			response.Error(w, http.StatusConflict, "email already exists")

		default:
			response.Error(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.JSON(w, http.StatusCreated, "user registered successfully", user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Login(r.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, ErrEmailRequired),
			errors.Is(err, ErrPasswordRequired):

			response.Error(w, http.StatusBadRequest, err.Error())

		case errors.Is(err, ErrInvalidCredentials):
			response.Error(w, http.StatusInternalServerError, err.Error())

		default:
			response.Error(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.JSON(w, http.StatusOK, "login successful", user)
}
