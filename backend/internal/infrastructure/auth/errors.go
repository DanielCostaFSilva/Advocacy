package auth

import "errors"

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpiredToken   = errors.New("token has expired")
	ErrGenerateToken  = errors.New("failed to generate token")
	ErrMissingSecret  = errors.New("JWT secret is not configured")
)
