package http

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/PzKpfw-ausf-H/forgevault/internal/service"
	"github.com/go-chi/chi/v5"
)

type AssetsHandler struct {
	svc       *service.AssetService
	templates *template.Template
}

func NewAssetsHandler(svc *service.AssetService) *AssetsHandler {
	return &AssetsHandler{svc: svc}
}

func (h *AssetsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req service.CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidJSON, "invalid json body", "")
		return
	}

	a, err := h.svc.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidTitle) || errors.Is(err, domain.ErrInvalidAssetType) {
			writeError(w, http.StatusBadRequest, ErrCodeValidation, err.Error(), "")
			return
		}
		if errors.Is(err, repo.ErrConflict) {
			writeError(w, http.StatusConflict, ErrCodeConflict, "asset already exists", "")
			return
		}

		writeError(w, http.StatusInternalServerError, ErrCodeInternal, "internal error", "")
	}
	writeJSON(w, http.StatusOK, toAssetResponse(a))

}

func (h *AssetsHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context(), repo.AssetFilter{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrCodeInternal, "internal error", "")
		return
	}
	writeJSON(w, http.StatusOK, items)

}

func (h *AssetsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id := domain.AssetID(idStr)

	asset, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			writeError(w, http.StatusNotFound, ErrCodeNotFound, "asset not found", "")
			return
		}

		writeError(w, http.StatusInternalServerError, ErrCodeInternal, "internal error", "")
		return
	}

	writeJSON(w, http.StatusOK, asset)
}

func (h *AssetsHandler) Patch(w http.ResponseWriter, r *http.Request) {
	var patch PatchAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "invalid patch request", "")
		return
	}

	idStr := chi.URLParam(r, "id")
	id := domain.AssetID(idStr)

	asset, err := h.svc.Patch(r.Context(), id, service.PatchAssetRequest(patch))
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrCodeInternal, "internal error", "")
		return
	}

	writeJSON(w, http.StatusOK, toAssetResponse(asset))

}

func (h *AssetsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id := domain.AssetID(idStr)

	err := h.svc.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			writeError(w, http.StatusNotFound, ErrCodeNotFound, "asset not found", "")
			return
		}

		writeError(w, http.StatusInternalServerError, ErrCodeInternal, "internal error", "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
