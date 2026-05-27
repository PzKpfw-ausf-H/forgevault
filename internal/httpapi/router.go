package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/PzKpfw-ausf-H/forgevault/internal/auth"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *AssetsHandler, tm *auth.TokenManager, ah *UsersHandler, fh *FilesHandler) http.Handler {
	r := chi.NewRouter()

	MountSwagger(r)

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
			r.Post("/{id}/files/upload-url", fh.UploadURL)
			r.Post("/{id}/files/confirm", fh.Confirm)
			r.Get("/{id}/files/{version}/download-url", fh.DownloadURL)
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

func MountSwagger(r chi.Router) {
	r.Get("/openapi/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "openapi/openapi.yaml")
	})

	//swagger UI static files
	fs := http.FileServer(http.Dir("web/swagger-ui"))
	r.Handle("/swagger/*", http.StripPrefix("/swagger/", fs))
	//redirect
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
}
