package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTService_Generate_ShouldReturnValidToken(t *testing.T) {
	svc := NewJWTService("my-secret", 60)

	token, err := svc.Generate("user-123", "john@example.com")
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestJWTService_Validate_ShouldReturnClaims(t *testing.T) {
	svc := NewJWTService("my-secret", 60)

	token, err := svc.Generate("user-123", "john@example.com")
	require.NoError(t, err)

	claims, err := svc.Validate(token)
	require.NoError(t, err)
	require.NotNil(t, claims)

	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "john@example.com", claims.Email)
	assert.False(t, claims.IssuedAt.IsZero())
	assert.False(t, claims.ExpiresAt.IsZero())
}

func TestJWTService_Validate_ShouldReturnErrExpiredToken(t *testing.T) {
	svc := NewJWTService("my-secret", 0)

	token, err := svc.Generate("user-123", "john@example.com")
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	claims, err := svc.Validate(token)
	assert.ErrorIs(t, err, ErrExpiredToken)
	assert.Nil(t, claims)
}

func TestJWTService_Validate_ShouldReturnErrInvalidTokenWhenMalformed(t *testing.T) {
	svc := NewJWTService("my-secret", 60)

	claims, err := svc.Validate("this-is-not-a-valid-jwt")
	assert.ErrorIs(t, err, ErrInvalidToken)
	assert.Nil(t, claims)
}

func TestJWTService_Validate_ShouldReturnErrInvalidTokenWhenWrongSignature(t *testing.T) {
	svc1 := NewJWTService("secret-1", 60)
	svc2 := NewJWTService("secret-2", 60)

	token, err := svc1.Generate("user-123", "john@example.com")
	require.NoError(t, err)

	claims, err := svc2.Validate(token)
	assert.ErrorIs(t, err, ErrInvalidToken)
	assert.Nil(t, claims)
}

func TestJWTService_Generate_ShouldReturnErrMissingSecret(t *testing.T) {
	svc := NewJWTService("", 60)

	token, err := svc.Generate("user-123", "john@example.com")
	assert.ErrorIs(t, err, ErrMissingSecret)
	assert.Empty(t, token)
}

func TestJWTService_Validate_ShouldReturnCorrectUserID(t *testing.T) {
	svc := NewJWTService("my-secret", 60)

	token, err := svc.Generate("user-456", "jane@example.com")
	require.NoError(t, err)

	claims, err := svc.Validate(token)
	require.NoError(t, err)
	assert.Equal(t, "user-456", claims.UserID)
}

func TestJWTService_Validate_ShouldReturnCorrectEmail(t *testing.T) {
	svc := NewJWTService("my-secret", 60)

	token, err := svc.Generate("user-456", "jane@example.com")
	require.NoError(t, err)

	claims, err := svc.Validate(token)
	require.NoError(t, err)
	assert.Equal(t, "jane@example.com", claims.Email)
}

func TestJWTService_Validate_ShouldHandleExpiredTokenFromPast(t *testing.T) {
	svc := NewJWTService("my-secret", 60)

	now := time.Now()
	claims := CustomClaims{
		UserID: "user-1",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-1",
			IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)),
		},
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := tokenObj.SignedString([]byte("my-secret"))
	require.NoError(t, err)

	result, err := svc.Validate(tokenStr)
	assert.ErrorIs(t, err, ErrExpiredToken)
	assert.Nil(t, result)
}
