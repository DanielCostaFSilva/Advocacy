package auth

import "github.com/go-chi/chi/v5"

func (h *LoginHandler) Register(r chi.Router) {
	r.Post("/auth/login", h.Login)
}
