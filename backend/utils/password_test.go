package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("securePass123")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "securePass123", hash)
}

func TestCheckPassword(t *testing.T) {
	password := "mySecretP@ss"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	t.Run("correct password returns true", func(t *testing.T) {
		assert.True(t, CheckPassword(password, hash))
	})

	t.Run("wrong password returns false", func(t *testing.T) {
		assert.False(t, CheckPassword("wrongPassword", hash))
	})

	t.Run("empty password returns false", func(t *testing.T) {
		assert.False(t, CheckPassword("", hash))
	})
}

func TestHashPasswordUnique(t *testing.T) {
	hash1, _ := HashPassword("samePassword")
	hash2, _ := HashPassword("samePassword")
	assert.NotEqual(t, hash1, hash2, "bcrypt should produce different salts")
}
