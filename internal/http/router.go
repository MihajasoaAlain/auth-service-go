package httpapi

import "github.com/go-chi/chi/v5"

func Router(api AuthAPI, corsAllowedOrigins []string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(CORSMiddleware(corsAllowedOrigins))

	r.Post("/auth/register", api.Register)
	r.Post("/auth/login", api.Login)
	r.Post("/auth/refresh", api.Refresh)
	r.Post("/auth/logout", api.Logout)
	r.Post("/auth/verify/start", api.VerifyStart)
	r.Post("/auth/verify/confirm", api.VerifyConfirm)
	r.Post("/auth/password/forgot", api.PasswordForgot)
	r.Post("/auth/password/reset", api.PasswordReset)
	r.Get("/auth/google/start", api.GoogleStart)
	r.Get("/auth/google/callback", api.GoogleCallback)
	r.Get("/auth/github/start", api.GitHubStart)
	r.Get("/auth/github/callback", api.GitHubCallback)

	r.Group(func(pr chi.Router) {
		pr.Use(api.AuthMiddleware)
		pr.Get("/me", api.Me)
	})

	return r
}
