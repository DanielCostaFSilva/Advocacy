package hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBcryptHasher_ShouldHashAndVerify(t *testing.T) {
	hasher := NewBcryptHasher()
	verifier := NewBcryptVerifier()

	password := "mySecurePassword123"
	hash, err := hasher.Hash(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	err = verifier.Compare(password, hash)
	assert.NoError(t, err)
}

func TestBcryptVerifier_ShouldReturnErrorOnWrongPassword(t *testing.T) {
	hasher := NewBcryptHasher()
	verifier := NewBcryptVerifier()

	hash, err := hasher.Hash("correctpassword")
	require.NoError(t, err)

	err = verifier.Compare("wrongpassword", hash)
	assert.Error(t, err)
}

func TestBcryptHasher_ShouldProduceDifferentHashes(t *testing.T) {
	hasher := NewBcryptHasher()

	hash1, err := hasher.Hash("samepassword")
	require.NoError(t, err)

	hash2, err := hasher.Hash("samepassword")
	require.NoError(t, err)

	assert.NotEqual(t, hash1, hash2)
}
