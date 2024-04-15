package router

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
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

	r.Use(httprate.LimitByIP(20, 1*time.Minute))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Route("/api/v1", func(r chi.Router) {

		r.Get("/health", h.CheckHealth)

		r.Route("/tags", func(r chi.Router) {

			r.With(jwtauth.Verifier(tokenAuth), jwtauth.Authenticator(tokenAuth)).Group(func(r chi.Router) {

				r.Post("/", h.CreateTag)
				r.Put("/", h.UpdateTag)
				r.Delete("/{id}", h.DeleteTag)

			})

			r.Get("/", h.GetAllTags)
			r.Get("/{id}", h.GetTagByID)
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

			r.Get("/", h.GetAllArticles)

			r.With(jwtauth.Verifier(tokenAuth), jwtauth.Authenticator(tokenAuth)).Group(func(r chi.Router) {

				r.Get("/all", h.GetAllArticlesByUser)
				r.Get("/search", h.SearchArticle)
				r.Post("/", h.CreateArticle)
				r.Post("/{id}", h.PublishArticle)

			})
		})

	})

	return r
}
