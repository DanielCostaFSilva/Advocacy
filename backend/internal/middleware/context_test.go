package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetUserID_ShouldReturnUserID(t *testing.T) {
	claims := &TokenClaims{UserID: "user-123"}
	ctx := context.WithValue(context.Background(), claimsKey, claims)

	id := GetUserID(ctx)
	assert.Equal(t, "user-123", id)
}

func TestGetUserID_ShouldReturnEmptyWhenNoClaims(t *testing.T) {
	id := GetUserID(context.Background())
	assert.Empty(t, id)
}

func TestGetUserEmail_ShouldReturnEmail(t *testing.T) {
	claims := &TokenClaims{Email: "john@example.com"}
	ctx := context.WithValue(context.Background(), claimsKey, claims)

	email := GetUserEmail(ctx)
	assert.Equal(t, "john@example.com", email)
}

func TestGetUserEmail_ShouldReturnEmptyWhenNoClaims(t *testing.T) {
	email := GetUserEmail(context.Background())
	assert.Empty(t, email)
}

func TestGetClaims_ShouldReturnClaims(t *testing.T) {
	expected := &TokenClaims{
		UserID:    "user-123",
		Email:     "john@example.com",
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	ctx := context.WithValue(context.Background(), claimsKey, expected)

	claims := GetClaims(ctx)
	assert.NotNil(t, claims)
	assert.Equal(t, expected.UserID, claims.UserID)
	assert.Equal(t, expected.Email, claims.Email)
}

func TestGetClaims_ShouldReturnNilWhenNoClaims(t *testing.T) {
	claims := GetClaims(context.Background())
	assert.Nil(t, claims)
}
