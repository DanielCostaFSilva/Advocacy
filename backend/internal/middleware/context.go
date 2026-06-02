package middleware

import (
	"context"
	"time"
)

type TokenClaims struct {
	UserID    string
	Email     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type authCtxKey string

const claimsKey authCtxKey = "auth_claims"

func GetClaims(ctx context.Context) *TokenClaims {
	claims, _ := ctx.Value(claimsKey).(*TokenClaims)
	return claims
}

func GetUserID(ctx context.Context) string {
	if claims := GetClaims(ctx); claims != nil {
		return claims.UserID
	}
	return ""
}

func GetUserEmail(ctx context.Context) string {
	if claims := GetClaims(ctx); claims != nil {
		return claims.Email
	}
	return ""
}
