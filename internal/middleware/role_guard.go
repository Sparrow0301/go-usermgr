package middleware

import (
	"net/http"
	"strings"

	"usermgmt/internal/errorx"
	"usermgmt/pkg/contextx"
	"usermgmt/pkg/response"
)

// NewRoleGuard creates a middleware that ensures the user has one of the required roles.
func NewRoleGuard(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	required := make(map[string]struct{})
	for _, role := range roles {
		role = strings.TrimSpace(role)
		if role == "" {
			continue
		}
		required[strings.ToLower(role)] = struct{}{}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if len(required) == 0 {
				next(w, r)
				return
			}

			claims := contextx.MustGetClaims(r.Context())
			if claims == nil {
				response.Error(w, r, errorx.ErrInvalidCredentials.Status, errorx.ErrInvalidCredentials.Code, errorx.ErrInvalidCredentials.Message, nil)
				return
			}

			for _, role := range claims.Roles {
				if _, ok := required[strings.ToLower(role)]; ok {
					next(w, r)
					return
				}
			}

			response.Error(w, r, errorx.ErrForbidden.Status, errorx.ErrForbidden.Code, errorx.ErrForbidden.Message, nil)
		}
	}
}
