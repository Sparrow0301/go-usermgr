package middleware

import (
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"usermgmt/internal/errorx"
	"usermgmt/pkg/contextx"
	"usermgmt/pkg/response"
	"usermgmt/pkg/security"
)

// AuthMiddleware validates JWT tokens from the Authorization header.
type AuthMiddleware struct {
	secret string
}

// NewAuthMiddleware creates a JWT middleware with the provided secret.
func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{secret: secret}
}

// Handle enforces bearer tokens and injects claims into the request context.
func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeUnauthorized(w, r)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeUnauthorized(w, r)
			return
		}

		claims, err := security.ParseToken(parts[1], m.secret)

		if err != nil {
			logx.WithContext(r.Context()).Errorf("parse token failed: %v", err)
			writeUnauthorized(w, r)
			return
		}

		ctx := contextx.WithClaims(r.Context(), claims)
		next(w, r.WithContext(ctx))
	}
}

func writeUnauthorized(w http.ResponseWriter, r *http.Request) {
	response.Error(w, r, errorx.ErrInvalidCredentials.Status, errorx.ErrInvalidCredentials.Code, errorx.ErrInvalidCredentials.Message, nil)
}
