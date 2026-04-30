package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/PzKpfw-ausf-H/forgevault/internal/service"
	"github.com/go-chi/chi/v5"
)

type AssetsHandler struct {
	svc *service.AssetService
}

func NewAssetsHandler(svc *service.AssetService) *AssetsHandler {
	return &AssetsHandler{svc: svc}
}

func (h *AssetsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidJSON, "invalid json body", "")
		return
	}

	svcReq := service.CreateAssetRequest{
		Title:       req.Title,
		Type:        req.Type,
		Tags:        make([]string, len(req.Tags)),
		Description: req.Description,
	}
	copy(svcReq.Tags, req.Tags)

	a, err := h.svc.Create(r.Context(), svcReq)
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
		return
	}
	writeJSON(w, http.StatusCreated, toAssetResponse(a))

}

func (h *AssetsHandler) List(w http.ResponseWriter, r *http.Request) {
	filter, err := parseFilter(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeValidation, err.Error(), "")
		return
	}
	items, err := h.svc.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrCodeInternal, "internal error", "")
		return
	}
	assetsRespone := []AssetResponse{}
	for _, asset := range items {
		assetsRespone = append(assetsRespone, toAssetResponse(asset))
	}
	writeJSON(w, http.StatusOK, assetsRespone)

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

	writeJSON(w, http.StatusOK, toAssetResponse(asset))
}

func (h *AssetsHandler) Patch(w http.ResponseWriter, r *http.Request) {
	var patch PatchAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidJSON, "invalid patch request", "")
		return
	}

	idStr := chi.URLParam(r, "id")
	id := domain.AssetID(idStr)

	asset, err := h.svc.Patch(r.Context(), id, service.PatchAssetRequest(patch))
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			writeError(w, http.StatusNotFound, ErrCodeNotFound, "asset not found", "")
			return
		}

		if errors.Is(err, domain.ErrInvalidTitle) || errors.Is(err, domain.ErrInvalidAssetType) || errors.Is(err, domain.ErrInvalidTagLen) || errors.Is(err, domain.ErrInvalidTag) {
			writeError(w, http.StatusBadRequest, ErrCodeValidation, "validation error", "")
			return
		}
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

func parseFilter(r *http.Request) (repo.AssetFilter, error) {
	q := r.URL.Query()

	var filter repo.AssetFilter

	if t := q.Get("type"); t != "" {
		at := domain.AssetType(strings.ToLower(t))
		filter.Type = &at
	}

	if t := q.Get("tag"); t != "" {
		tag := strings.ToLower(strings.TrimSpace(t))
		filter.Tag = &tag
	}

	if t := q.Get("titleSub"); t != "" {
		titleSub := strings.ToLower(strings.TrimSpace(t))
		filter.TitleSub = &titleSub
	}

	if a := q.Get("authorId"); a != "" {
		authorID := domain.UserID(strings.TrimSpace(a))
		filter.AuthorID = &authorID
	}

	if l := q.Get("limit"); l != "" {
		lim, err := strconv.Atoi(l)
		if err != nil {
			return repo.AssetFilter{}, fmt.Errorf("error parsing limit")
		}
		filter.Limit = lim
	}

	if o := q.Get("offset"); o != "" {
		if offset, err := strconv.Atoi(o); err == nil {
			filter.Offset = offset
		} else {
			return repo.AssetFilter{}, fmt.Errorf("error parsing filter")
		}
	}

	return filter, nil
}
