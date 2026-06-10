package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWT(t *testing.T) {
	token, err := GenerateJWT(42, "admin")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	t.Run("valid token returns claims", func(t *testing.T) {
		token, err := GenerateJWT(1, "client")
		require.NoError(t, err)

		claims, err := ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, uint(1), claims.UserID)
		assert.Equal(t, "client", claims.Role)
	})

	t.Run("invalid signature is rejected", func(t *testing.T) {
		claims := &Claims{
			UserID: 99,
			Role:   "admin",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		tampered := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tamperedString, err := tampered.SignedString([]byte("different-secret"))
		require.NoError(t, err)

		_, err = ValidateToken(tamperedString)
		assert.Error(t, err)
	})

	t.Run("malformed token is rejected", func(t *testing.T) {
		_, err := ValidateToken("not.a.token")
		assert.Error(t, err)
	})

	t.Run("empty string is rejected", func(t *testing.T) {
		_, err := ValidateToken("")
		assert.Error(t, err)
	})
}
