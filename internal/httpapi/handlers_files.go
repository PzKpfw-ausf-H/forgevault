package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/PzKpfw-ausf-H/forgevault/internal/service"
	"github.com/go-chi/chi/v5"
)

type FilesHandler struct {
	svc *service.FileService
}

func NewFilesHandler(svc *service.FileService) *FilesHandler {
	return &FilesHandler{svc: svc}
}

func (fh *FilesHandler) UploadURL(w http.ResponseWriter, r *http.Request) {
	var req UploadURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "invalid json", "")
		return
	}

	uid, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "missing auth context", "")
		return
	}

	assetID := domain.AssetID(chi.URLParam(r, "id"))

	version, storageKey, uploadURL, expiresInSec, err := fh.svc.GetUploadURL(r.Context(), uid, assetID, req.Filename, req.ContentType)
	if err != nil {
		//errors mapping
		if errors.Is(err, repo.ErrValidation) {
			writeError(w, http.StatusBadRequest, ErrCodeValidation, "validation", "")
			return
		}
		if errors.Is(err, repo.ErrBadRequest) {
			writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "bad request", "")
			return
		}
		if errors.Is(err, repo.ErrUnauthorized) {
			writeError(w, http.StatusForbidden, ErrCodeUnauthorized, "unauthorized", "")
			return
		}
		if errors.Is(err, repo.ErrNotFound) {
			writeError(w, http.StatusNotFound, ErrCodeNotFound, "not found", "")
			return
		}

		writeError(w, http.StatusInternalServerError, ErrCodeInternal, "internal error", "")
		return
	}

	writeJSON(w, http.StatusOK, toUploadResponse(version, storageKey, uploadURL, expiresInSec))
}

func (fh *FilesHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	var req ConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "invalid json", "")
		return
	}

	uid, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "missing auth context", "")
		return
	}

	assetID := domain.AssetID(chi.URLParam(r, "id"))

	_, err := fh.svc.ConfirmUpload(r.Context(), uid, assetID, req.Version, req.Filename, req.ContentType, req.SizeBytes, req.StorageKey, req.Checksum)
	if err != nil {
		//errors mapping
		if errors.Is(err, repo.ErrValidation) {
			writeError(w, http.StatusBadRequest, ErrCodeValidation, "validation", "")
			return
		}
		if errors.Is(err, repo.ErrBadRequest) {
			writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "bad request", "")
			return
		}
		if errors.Is(err, repo.ErrUnauthorized) {
			writeError(w, http.StatusForbidden, ErrCodeUnauthorized, "unauthorized", "")
			return
		}
		if errors.Is(err, repo.ErrNotFound) {
			writeError(w, http.StatusNotFound, ErrCodeNotFound, "not found", "")
			return
		}

		writeError(w, http.StatusInternalServerError, ErrCodeInternal, "internal error", "")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (fh *FilesHandler) DownloadURL(w http.ResponseWriter, r *http.Request) {
	uid, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "missing auth context", "")
		return
	}
	assetID := domain.AssetID(chi.URLParam(r, "id"))
	version, err := strconv.Atoi(chi.URLParam(r, "version"))
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "invalid version", "")
		return
	}
	downloadUrl, expiresIn, err := fh.svc.GetDownloadURL(r.Context(), uid, assetID, version)
	if err != nil {
		//errors mapping
		if errors.Is(err, repo.ErrValidation) {
			writeError(w, http.StatusBadRequest, ErrCodeValidation, "validation", "")
			return
		}
		if errors.Is(err, repo.ErrBadRequest) {
			writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "bad request", "")
			return
		}
		if errors.Is(err, repo.ErrUnauthorized) {
			writeError(w, http.StatusForbidden, ErrCodeUnauthorized, "unauthorized", "")
			return
		}
		if errors.Is(err, repo.ErrNotFound) {
			writeError(w, http.StatusNotFound, ErrCodeNotFound, "not found", "")
			return
		}

		writeError(w, http.StatusInternalServerError, ErrCodeInternal, "internal error", "")
		return
	}

	writeJSON(w, http.StatusOK, toDownloadResponse(downloadUrl, expiresIn))
}
