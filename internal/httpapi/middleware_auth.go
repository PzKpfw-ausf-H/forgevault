package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/PzKpfw-ausf-H/forgevault/internal/auth"
	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
)

type ctxKey string

const ctxKeyUserID ctxKey = "userID"

func UserIDFromContext(ctx context.Context) (domain.UserID, bool) {
	v := ctx.Value(ctxKeyUserID)
	uid, ok := v.(domain.UserID)
	return uid, ok
}

func RequireAuth(tm *auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" {
				writeError(w, http.StatusUnauthorized, "unauthorized", "missing authorization header", "")
				return
			}
			const prefix = "Bearer "
			if !strings.HasPrefix(h, prefix) {
				writeError(w, http.StatusUnauthorized, "unauthorized", "invalid authorization scheme", "")
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(h, prefix))
			if token == "" {
				writeError(w, http.StatusUnauthorized, "unauthorized", "empty token", "")
				return
			}

			userID, err := tm.Parse(token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "unauthorized", "invalid token", "")
				return
			}
			ctx := context.WithValue(r.Context(), ctxKeyUserID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
