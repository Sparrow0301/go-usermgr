package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"usermgmt/internal/types"
)

// GenerateToken builds a signed JWT token for the current user.
func GenerateToken(userID uint, roles []string, secret string, expireSeconds time.Duration) (string, time.Time, error) {
	if secret == "" {
		return "", time.Time{}, errors.New("jwt secret missing")
	}
	if expireSeconds == 0 {
		expireSeconds = time.Hour
	}

	expireAt := time.Now().Add(expireSeconds)
	claims := types.JwtClaims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expireAt, nil
}

// ParseToken validates token string and returns claims.
func ParseToken(tokenStr, secret string) (*types.JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &types.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*types.JwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
