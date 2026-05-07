package httpapi

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int64  `json:"expiresIn"`
}

func toTokenResponse(token string, expiresIn int64) TokenResponse {
	var tr TokenResponse
	tr.AccessToken = token
	tr.ExpiresIn = expiresIn
	tr.TokenType = "Bearer"

	return tr
}
