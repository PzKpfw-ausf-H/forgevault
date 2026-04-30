package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h *AssetsHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Route("/assets", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Patch("/{id}", h.Patch)
		r.Delete("/{id}", h.Delete)
	})

	return r
}
