package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/PzKpfw-ausf-H/forgevault/internal/service"
)

type UsersHandler struct {
	svc *service.UserService
}

func NewUsersHandler(svc *service.UserService) *UsersHandler {
	return &UsersHandler{svc: svc}
}

func (uh *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "invalid json", "")
		return
	}

	email := req.Email
	password := req.Password

	accessToken, refreshToken, expiresIn, refreshExpiresIn, err := uh.svc.Login(r.Context(), email, password)
	if err != nil {
		log.Printf("LOGIN ERROR: %+v", err)
		if errors.Is(err, repo.ErrUnauthorized) {
			writeError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "invalid credentials", "")
			return
		}
		writeError(w, http.StatusInternalServerError, ErrCodeInternal, "internal error", "")
		return
	}

	writeJSON(w, http.StatusOK, toTokenResponse(accessToken, refreshToken, expiresIn, refreshExpiresIn))
}

func (uh *UsersHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "invalid json", "")
		return
	}
	user, err := uh.svc.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("REGISTER ERROR: %+v", err)

		if errors.Is(err, repo.ErrConflict) {
			writeError(w, http.StatusConflict, ErrCodeConflict, "conflict", "")
			return
		}
		if errors.Is(err, repo.ErrUnauthorized) {
			writeError(w, http.StatusBadRequest, ErrCodeValidation, "invalid email or password", "")
			return
		}
		writeError(w, http.StatusInternalServerError, ErrCodeInternal, "internal error", "")
		return
	}

	writeJSON(w, http.StatusCreated, toRegisterResponse(user))
}

func (uh *UsersHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "invalid json", "")
		return
	}

	if req.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, ErrCodeValidation, "refreshToken is required", "")
		return
	}

	newAccess, newRefresh, accessExpiresIn, refreshExpiresIn, err := uh.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, repo.ErrUnauthorized) {
			writeError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "error while refreshing", "")
			return
		}

		writeError(w, http.StatusInternalServerError, ErrCodeInternal, "internal error", "")
		return
	}

	writeJSON(w, http.StatusOK, toTokenResponse(newAccess, newRefresh, accessExpiresIn, refreshExpiresIn))
}

func (uh *UsersHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "invalid json", "")
		return
	}

	if req.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, ErrCodeValidation, "refreshToken is required", "")
		return
	}

	_ = uh.svc.Logout(r.Context(), req.RefreshToken)

	w.WriteHeader(http.StatusNoContent)
}
