package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenProvider interface {
	Generate(userID, email string) (string, error)
	Validate(token string) (*TokenClaims, error)
}

type JWTService struct {
	secret     string
	expiration time.Duration
}

func NewJWTService(secret string, expirationMinutes int) *JWTService {
	return &JWTService{
		secret:     secret,
		expiration: time.Duration(expirationMinutes) * time.Minute,
	}
}

func (s *JWTService) Generate(userID, email string) (string, error) {
	if s.secret == "" {
		return "", ErrMissingSecret
	}

	now := time.Now()
	claims := CustomClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", errors.Join(ErrGenerateToken, err)
	}

	return tokenStr, nil
}

func (s *JWTService) Validate(tokenStr string) (*TokenClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return &TokenClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}
