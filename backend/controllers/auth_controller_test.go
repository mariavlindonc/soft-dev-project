package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
