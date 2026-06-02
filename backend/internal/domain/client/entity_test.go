package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient_ShouldCreateSuccessfully(t *testing.T) {
	client, err := NewClient("John Doe", "12345678901", "john@example.com", "11999999999")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotEmpty(t, client.ID)
	assert.Equal(t, "John Doe", client.Name)
	assert.Equal(t, "12345678901", client.CPF)
	assert.Equal(t, "john@example.com", client.Email)
	assert.Equal(t, "11999999999", client.Phone)
	assert.False(t, client.CreatedAt.IsZero())
	assert.False(t, client.UpdatedAt.IsZero())
}

func TestNewClient_ShouldReturnErrInvalidNameWhenEmpty(t *testing.T) {
	client, err := NewClient("", "12345678901", "john@example.com", "11999999999")
	assert.ErrorIs(t, err, ErrInvalidName)
	assert.Nil(t, client)
}

func TestNewClient_ShouldReturnErrInvalidNameWhenLessThan3(t *testing.T) {
	client, err := NewClient("Jo", "12345678901", "john@example.com", "11999999999")
	assert.ErrorIs(t, err, ErrInvalidName)
	assert.Nil(t, client)
}

func TestNewClient_ShouldReturnErrInvalidCPFWhenEmpty(t *testing.T) {
	client, err := NewClient("John Doe", "", "john@example.com", "11999999999")
	assert.ErrorIs(t, err, ErrInvalidCPF)
	assert.Nil(t, client)
}

func TestNewClient_ShouldReturnErrInvalidCPFWhenLessThan11Digits(t *testing.T) {
	client, err := NewClient("John Doe", "1234567890", "john@example.com", "11999999999")
	assert.ErrorIs(t, err, ErrInvalidCPF)
	assert.Nil(t, client)
}

func TestNewClient_ShouldReturnErrInvalidCPFWhenMoreThan11Digits(t *testing.T) {
	client, err := NewClient("John Doe", "123456789012", "john@example.com", "11999999999")
	assert.ErrorIs(t, err, ErrInvalidCPF)
	assert.Nil(t, client)
}

func TestNewClient_ShouldNormalizeCPFWithPunctuation(t *testing.T) {
	client, err := NewClient("John Doe", "123.456.789-00", "john@example.com", "11999999999")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "12345678900", client.CPF)
}

func TestNewClient_ShouldReturnErrInvalidEmailWhenInvalidFormat(t *testing.T) {
	client, err := NewClient("John Doe", "12345678901", "invalid-email", "11999999999")
	assert.ErrorIs(t, err, ErrInvalidEmail)
	assert.Nil(t, client)
}

func TestNewClient_ShouldAllowEmptyEmail(t *testing.T) {
	client, err := NewClient("John Doe", "12345678901", "", "11999999999")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Empty(t, client.Email)
}

func TestNewClient_ShouldAllowEmptyPhone(t *testing.T) {
	client, err := NewClient("John Doe", "12345678901", "john@example.com", "")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Empty(t, client.Phone)
}
