package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/jay-bhogayata/blogapi/handlers"
	"github.com/jay-bhogayata/blogapi/middleware"
)

func NewRouter(h *handlers.Handlers) *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.Logging)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", h.CheckHealth)
	})

	return r
}
