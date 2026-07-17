package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestParseOrigins(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty string returns wildcard", "", []string{"*"}},
		{"single origin", "https://example.com", []string{"https://example.com"}},
		{"multiple origins comma separated", "https://a.com,https://b.com", []string{"https://a.com", "https://b.com"}},
		{"trims spaces", " https://a.com , https://b.com ", []string{"https://a.com", "https://b.com"}},
		{"filters empty parts from trailing comma", "https://a.com,,https://b.com,", []string{"https://a.com", "https://b.com"}},
		{"only commas returns wildcard", ",,,", []string{"*"}},
		{"wildcard origin", "*", []string{"*"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOrigins(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsOriginAllowed(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		allowed  []string
		expected bool
	}{
		{"exact match", "https://a.com", []string{"https://a.com", "https://b.com"}, true},
		{"no match", "https://c.com", []string{"https://a.com", "https://b.com"}, false},
		{"wildcard allows all", "https://anything.com", []string{"*"}, true},
		{"empty allowed list", "https://a.com", []string{}, false},
		{"empty origin not matched", "", []string{"https://a.com"}, false},
		{"empty origin with wildcard", "", []string{"*"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOriginAllowed(tt.origin, tt.allowed)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(securityHeadersMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "max-age=31536000; includeSubDomains", w.Header().Get("Strict-Transport-Security"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "default-src 'self'", w.Header().Get("Content-Security-Policy"))
}

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("preflight OPTIONS returns 204", func(t *testing.T) {
		os.Unsetenv("CORS_ALLOWED_ORIGINS")
		router := gin.New()
		router.Use(corsMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "https://example.com")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, PATCH, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Authorization, Content-Type", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("allowed origin echoes back", func(t *testing.T) {
		os.Setenv("CORS_ALLOWED_ORIGINS", "https://allowed.com")
		defer os.Unsetenv("CORS_ALLOWED_ORIGINS")
		router := gin.New()
		router.Use(corsMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "https://allowed.com")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "https://allowed.com", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("disallowed origin not echoed", func(t *testing.T) {
		os.Setenv("CORS_ALLOWED_ORIGINS", "https://allowed.com")
		defer os.Unsetenv("CORS_ALLOWED_ORIGINS")
		router := gin.New()
		router.Use(corsMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "https://evil.com")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("no origin header with wildcard", func(t *testing.T) {
		os.Unsetenv("CORS_ALLOWED_ORIGINS")
		router := gin.New()
		router.Use(corsMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("request continues for non-OPTIONS", func(t *testing.T) {
		os.Unsetenv("CORS_ALLOWED_ORIGINS")
		router := gin.New()
		router.Use(corsMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "reached")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "reached")
	})
}
