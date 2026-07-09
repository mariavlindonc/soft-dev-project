package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"backend/domain"
	"backend/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestRegister(t *testing.T) {
	t.Run("valid request returns 201", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		ctrl := NewAuthController(mockSvc)

		mockSvc.On("Register", mock.MatchedBy(func(input services.RegisterInput) bool {
			return input.Email == "a@test.com"
		})).Return(&domain.User{ID: 1, Name: "Alice", Email: "a@test.com", Role: "client"}, nil)
		mockSvc.On("GenerateToken", uint(1), "client").Return("jwt-token-123", nil)

		r := setupRouter()
		r.POST("/register", ctrl.Register)

		body := `{"name":"Alice","email":"a@test.com","password":"secure123"}`
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		var resp registerResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, uint(1), resp.User.ID)
		assert.Equal(t, "Alice", resp.User.Name)
	})

	t.Run("missing fields returns 400", func(t *testing.T) {
		ctrl := NewAuthController(new(MockAuthService))
		r := setupRouter()
		r.POST("/register", ctrl.Register)

		body := `{"name":""}`
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("duplicate email returns 400", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		ctrl := NewAuthController(mockSvc)

		mockSvc.On("Register", mock.Anything).Return(nil, services.ErrInvalidInput)

		r := setupRouter()
		r.POST("/register", ctrl.Register)

		body := `{"name":"A","email":"dup@test.com","password":"secure123"}`
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLogin(t *testing.T) {
	t.Run("valid credentials returns 200 with token", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		ctrl := NewAuthController(mockSvc)

		mockSvc.On("Login", services.LoginInput{
			Email:    "a@test.com",
			Password: "correct",
		}).Return("jwt-token-123", &domain.User{ID: 1, Name: "Alice", Email: "a@test.com", Role: "client"}, nil)

		r := setupRouter()
		r.POST("/login", ctrl.Login)

		body := `{"email":"a@test.com","password":"correct"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp loginResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "jwt-token-123", resp.Token)
	})

	t.Run("invalid credentials returns 401", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		ctrl := NewAuthController(mockSvc)

		mockSvc.On("Login", mock.Anything).Return("", nil, errors.New("unauthorized"))

		r := setupRouter()
		r.POST("/login", ctrl.Login)

		body := `{"email":"a@test.com","password":"wrong"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("missing password returns 400", func(t *testing.T) {
		ctrl := NewAuthController(new(MockAuthService))
		r := setupRouter()
		r.POST("/login", ctrl.Login)

		body := `{"email":"a@test.com"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRegisterTokenError(t *testing.T) {
	t.Run("GenerateToken error returns 500", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		ctrl := NewAuthController(mockSvc)

		mockSvc.On("Register", mock.Anything).Return(&domain.User{ID: 1, Name: "A", Email: "a@test.com", Role: "client"}, nil)
		mockSvc.On("GenerateToken", uint(1), "client").Return("", fmt.Errorf("token error"))

		r := setupRouter()
		r.POST("/register", ctrl.Register)

		body := `{"name":"A","email":"a@test.com","password":"secure123"}`
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("server error returns 500", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		ctrl := NewAuthController(mockSvc)

		mockSvc.On("Register", mock.Anything).Return(nil, fmt.Errorf("db error"))

		r := setupRouter()
		r.POST("/register", ctrl.Register)

		body := `{"name":"A","email":"a@test.com","password":"secure123"}`
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserFriendlyError(t *testing.T) {
	t.Run("email already registered", func(t *testing.T) {
		msg := userFriendlyError(fmt.Errorf("email is already registered"))
		assert.Equal(t, "El correo electrónico ya está registrado", msg)
	})

	t.Run("password too short", func(t *testing.T) {
		msg := userFriendlyError(fmt.Errorf("password must be at least 8 characters"))
		assert.Equal(t, "La contraseña debe tener al menos 8 caracteres", msg)
	})

	t.Run("unknown error", func(t *testing.T) {
		msg := userFriendlyError(fmt.Errorf("something else"))
		assert.Contains(t, msg, "Error interno")
	})
}

func TestHelpers(t *testing.T) {
	t.Run("isNotFound with not found error", func(t *testing.T) {
		assert.True(t, isNotFound(services.ErrNotFound))
	})

	t.Run("isNotFound with other error", func(t *testing.T) {
		assert.False(t, isNotFound(errors.New("not found")))
	})

	t.Run("optionalString empty returns nil", func(t *testing.T) {
		assert.Nil(t, optionalString(""))
	})

	t.Run("optionalString non-empty returns pointer", func(t *testing.T) {
		s := optionalString("hello")
		require.NotNil(t, s)
		assert.Equal(t, "hello", *s)
	})

	t.Run("timePtrToString nil returns nil", func(t *testing.T) {
		assert.Nil(t, timePtrToString(nil))
	})

	t.Run("timePtrToString non-nil returns formatted string", func(t *testing.T) {
		ts := time.Date(2026, 6, 15, 20, 0, 0, 0, time.UTC)
		s := timePtrToString(&ts)
		require.NotNil(t, s)
		assert.Contains(t, *s, "2026")
	})
}
