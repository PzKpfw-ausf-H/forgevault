package httpapi

import "github.com/PzKpfw-ausf-H/forgevault/internal/service"

type UsersHandler struct {
	svc *service.UserService
}

func NewUsersHandler(svc *service.UserService) *UsersHandler {
	return &UsersHandler{svc: svc}
}
