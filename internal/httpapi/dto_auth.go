package httpapi

import "github.com/PzKpfw-ausf-H/forgevault/internal/domain"

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type TokenResponse struct {
	AccessToken      string `json:"accessToken"`
	TokenType        string `json:"tokenType"`
	ExpiresIn        int64  `json:"expiresIn"`
	RefreshToken     string `json:"refreshToken"`
	RefreshExpiresIn int64  `json:"refreshExpiresIn"`
}

type RegisterResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func toTokenResponse(accessToken, refreshToken string, expiresIn, refreshExpiresIn int64) TokenResponse {
	var tr TokenResponse
	tr.AccessToken = accessToken
	tr.RefreshToken = refreshToken
	tr.ExpiresIn = expiresIn
	tr.RefreshExpiresIn = refreshExpiresIn
	tr.TokenType = "Bearer"

	return tr
}

func toRegisterResponse(u domain.User) RegisterResponse {
	var r RegisterResponse
	r.ID = string(u.ID)
	r.Email = u.Email

	return r
}
