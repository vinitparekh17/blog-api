package router

import (
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

// func NewRouter(cfg *config.Config, h *handlers.Handlers) *chi.Mux {

// tokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

// r := chi.NewRouter()

// r.Use(middleware.Logging)

// r.Use(httprate.LimitByIP(20, 1*time.Minute))

// r.Use(cors.Handler(cors.Options{
// 	AllowedOrigins:   []string{"*"},
// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
// 	AllowCredentials: true,
// }))

// r.Route("/api/v1", func(r chi.Router) {

// 	r.Get("/health", h.CheckHealth)

// 	r.Route("/tags", func(r chi.Router) {

// 		r.With(jwtauth.Verifier(tokenAuth), jwtauth.Authenticator(tokenAuth)).Group(func(r chi.Router) {

// 			r.Post("/", h.CreateTag)
// 			r.Put("/", h.UpdateTag)
// 			r.Delete("/{id}", h.DeleteTag)

// 		})

// 		r.Get("/", h.GetAllTags)
// 		r.Get("/{id}", h.GetTagByID)
// 	})
// 	r.Route("/users", func(r chi.Router) {

// 		r.Post("/signup", h.RegisterUser)
// 		r.Get("/verify", h.VerifyUser)
// 		r.Post("/signin", h.LoginUser)
// 		r.Get("/{mail}", h.GetUserInfoByUserEmail)

// 		r.Get("/", h.GetAllUsers)
// 		r.With(jwtauth.Verifier(tokenAuth), jwtauth.Authenticator(tokenAuth)).Group(func(r chi.Router) {
// 			r.Get("/posts", h.GetAllArticlesByUser)
// 			r.Post("/logout", h.LogoutUser)
// 			r.Delete("/delete", h.DeleteUser)

// 		})

// 	})

// 	r.Route("/posts", func(r chi.Router) {

// 		r.Get("/", h.GetAllArticles)
// 		r.Get("/search", h.SearchArticle)

// 		r.With(jwtauth.Verifier(tokenAuth), jwtauth.Authenticator(tokenAuth)).Group(func(r chi.Router) {
// 			r.Post("/", h.CreateArticle)
// 			r.Put("/publish/{id}", h.PublishArticle)
// 		})
// 	})

// })

// return r
// }
