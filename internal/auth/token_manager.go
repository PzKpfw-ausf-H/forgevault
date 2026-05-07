package auth

import (
	"fmt"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secret []byte
	ttl    time.Duration
	issuer string
}

type AuthClaims struct {
	jwt.RegisteredClaims
}

func NewTokenManager(secret []byte, ttl time.Duration, issuer string) *TokenManager {
	return &TokenManager{
		secret: secret,
		ttl:    ttl,
		issuer: issuer,
	}
}

func (tm *TokenManager) New(userID domain.UserID) (string, int64, error) {
	now := time.Now().UTC()
	exp := now.Add(tm.ttl)

	claims := AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   string(userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    tm.issuer,
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(tm.secret)
	if err != nil {
		return "", 0, err
	}

	return signed, int64(tm.ttl.Seconds()), nil
}

func (tm *TokenManager) Parse(token string) (domain.UserID, error) {
	var claims AuthClaims

	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return tm.secret, nil
	})

	if err != nil {
		return "", err
	}

	if !parsedToken.Valid {
		return "", ErrInvalidToken
	}

	if tm.issuer != "" && claims.Issuer != tm.issuer {
		return "", ErrInvalidToken
	}

	if claims.Subject == "" {
		return "", fmt.Errorf("%w: empty subject", ErrInvalidToken)
	}

	userID := domain.UserID(claims.Subject)
	return userID, nil

}
