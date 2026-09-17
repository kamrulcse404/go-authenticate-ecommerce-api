package auth

import "github.com/go-chi/chi/v5"

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()

	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	
	return r
}
