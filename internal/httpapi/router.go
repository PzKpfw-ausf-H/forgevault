package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/PzKpfw-ausf-H/forgevault/internal/auth"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *AssetsHandler, tm *auth.TokenManager, ah *UsersHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Route("/assets", func(r chi.Router) {
		//public
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)

		//protected
		r.Group(func(r chi.Router) {
			r.Use(RequireAuth(tm))
			r.Post("/", h.Create)
			r.Patch("/{id}", h.Patch)
			r.Delete("/{id}", h.Delete)
		})
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", ah.Register)
		r.Post("/login", ah.Login)
		r.Post("/refresh", ah.Refresh)
		r.Post("/logout", ah.Logout)
	})

	return r
}
