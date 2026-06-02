package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser_ShouldCreateUserWhenValid(t *testing.T) {
	user, err := NewUser("John Doe", "john@example.com", "hashedpassword123")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "hashedpassword123", user.PasswordHash)
	assert.NotEmpty(t, user.ID)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestNewUser_ShouldReturnErrorWhenInvalidEmail(t *testing.T) {
	user, err := NewUser("John Doe", "invalid-email", "hashedpassword123")

	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrInvalidEmail)
}

func TestNewUser_ShouldReturnErrorWhenEmptyEmail(t *testing.T) {
	user, err := NewUser("John Doe", "", "hashedpassword123")

	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrInvalidEmail)
}

func TestNewUser_ShouldReturnErrorWhenEmptyPasswordHash(t *testing.T) {
	user, err := NewUser("John Doe", "john@example.com", "")

	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrInvalidPasswordHash)
}

func TestNewUser_ShouldReturnErrorWhenShortName(t *testing.T) {
	user, err := NewUser("ab", "john@example.com", "hashedpassword123")

	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrInvalidName)
}

func TestNewUser_ShouldReturnErrorWhenEmptyName(t *testing.T) {
	user, err := NewUser("", "john@example.com", "hashedpassword123")

	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrInvalidName)
}
