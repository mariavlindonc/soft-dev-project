package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault(t *testing.T) {
	t.Run("returns env value when set", func(t *testing.T) {
		os.Setenv("TEST_EXISTS", "hello")
		defer os.Unsetenv("TEST_EXISTS")
		assert.Equal(t, "hello", getEnvOrDefault("TEST_EXISTS", "fallback"))
	})

	t.Run("returns fallback when env not set", func(t *testing.T) {
		os.Unsetenv("TEST_MISSING")
		assert.Equal(t, "fallback", getEnvOrDefault("TEST_MISSING", "fallback"))
	})

	t.Run("returns empty fallback when env not set", func(t *testing.T) {
		os.Unsetenv("TEST_MISSING_EMPTY")
		assert.Equal(t, "", getEnvOrDefault("TEST_MISSING_EMPTY", ""))
	})

	t.Run("empty env returns fallback", func(t *testing.T) {
		os.Setenv("TEST_EMPTY", "")
		defer os.Unsetenv("TEST_EMPTY")
		assert.Equal(t, "fallback", getEnvOrDefault("TEST_EMPTY", "fallback"))
	})
}
