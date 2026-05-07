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

	token, expiresIn, err := uh.svc.Login(r.Context(), email, password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "invalid token", "")
		return
	}

	writeJSON(w, http.StatusOK, toTokenResponse(token, expiresIn))
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
