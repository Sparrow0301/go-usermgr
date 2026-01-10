package contextx

import (
	"context"

	"usermgmt/internal/types"
)

type contextKey string

const (
	claimsKey contextKey = "authClaims"
)

// WithClaims stores JWT claims into context.
func WithClaims(ctx context.Context, claims *types.JwtClaims) context.Context {
	if claims == nil {
		return ctx
	}
	return context.WithValue(ctx, claimsKey, claims)
}

// ClaimsFromContext extracts JWT claims from context.
func ClaimsFromContext(ctx context.Context) (*types.JwtClaims, bool) {
	if ctx == nil {
		return nil, false
	}
	val, ok := ctx.Value(claimsKey).(*types.JwtClaims)
	return val, ok && val != nil
}

// MustGetClaims returns claims or nil if absent.
func MustGetClaims(ctx context.Context) *types.JwtClaims {
	claims, _ := ClaimsFromContext(ctx)
	return claims
}
