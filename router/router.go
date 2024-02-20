package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"github.com/jay-bhogayata/blogapi/config"
	"github.com/jay-bhogayata/blogapi/handlers"
	"github.com/jay-bhogayata/blogapi/middleware"
)

var tokenAuth *jwtauth.JWTAuth

func NewRouter(cfg *config.Config, h *handlers.Handlers) *chi.Mux {

	tokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	r := chi.NewRouter()

	r.Use(middleware.Logging)

	r.Route("/api/v1", func(r chi.Router) {

		r.Get("/health", h.CheckHealth)

		r.Route("/catagories", func(r chi.Router) {

			r.With(jwtauth.Verifier(tokenAuth), jwtauth.Authenticator(tokenAuth)).Group(func(r chi.Router) {

				r.Post("/", h.CreateCategory)
				r.Put("/", h.UpdateCategory)
				r.Delete("/{id}", h.DeleteCategory)

			})

			r.Get("/", h.GetAllCategories)
			r.Get("/{id}", h.GetCategoryByID)
		})
		r.Route("/accounts", func(r chi.Router) {
			r.Post("/register", h.RegisterUser)
			r.Get("/verify", h.VerifyUser)
			r.Post("/login", h.LoginUser)
			r.Get("/user/{mail}", h.GetUserInfoByUserEmail)
			r.With(jwtauth.Verifier(tokenAuth), jwtauth.Authenticator(tokenAuth)).Group(func(r chi.Router) {
				r.Post("/logout", h.LogoutUser)
				r.Get("/users", h.GetAllUsers)
			})

		})

		r.Route("/article", func(r chi.Router) {

			r.With(jwtauth.Verifier(tokenAuth), jwtauth.Authenticator(tokenAuth)).Group(func(r chi.Router) {

				r.Get("/all", h.GetAllArticlesByUser)
				r.Get("/", h.GetAllArticles)
				r.Post("/", h.CreateArticle)
				r.Post("/{id}", h.PublishArticle)

			})
		})

	})

	return r
}
