package httpapi

import (
	"encoding/json"
	"net/http"

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
		writeError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "login error", "")
		return
	}

	email := req.Email
	password := req.Password

	accessToken, refreshToken, expiresIn, refreshExpiresIn, err := uh.svc.Login(r.Context(), email, password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "invalid token", "")
		return
	}

	writeJSON(w, http.StatusOK, toTokenResponse(accessToken, refreshToken, expiresIn, refreshExpiresIn))
}

func (uh *UsersHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "invalid email or password", "")
		return
	}
	user, err := uh.svc.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "invalid email or passowrd", "")
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (uh *UsersHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "invalid refresh token", "")
		return
	}

	newAccess, newRefresh, accessExpiresIn, refreshExpiresIn, err := uh.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		writeError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "error while refreshing", "")
		return
	}

	writeJSON(w, http.StatusOK, toTokenResponse(newAccess, newRefresh, accessExpiresIn, refreshExpiresIn))
}

func (uh *UsersHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "bad request", "")
		return
	}

	if err := uh.svc.Logout(r.Context(), req.RefreshToken); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "bad request", "")
		return
	}

	writeJSON(w, http.StatusNoContent, "no content")
}
