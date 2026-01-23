package httpapi

import "github.com/go-chi/chi/v5"

func Router(api AuthAPI) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/auth/register", api.Register)
	r.Post("/auth/login", api.Login)

	r.Group(func(pr chi.Router) {
		pr.Use(api.AuthMiddleware)
		pr.Get("/me", api.Me)
	})

	return r
}
