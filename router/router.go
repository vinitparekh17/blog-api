package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/jay-bhogayata/blogapi/handlers"
)

func NewRouter(h *handlers.Handlers) *chi.Mux {

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", h.CheckHealth)
	})

	return r
}
