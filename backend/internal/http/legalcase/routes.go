package legalcase

import "github.com/go-chi/chi/v5"

func (h *Handler) Register(r chi.Router) {
	r.Post("/cases", h.Create)
}
