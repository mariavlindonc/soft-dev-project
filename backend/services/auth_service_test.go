package services

import (
	"errors"
	"testing"

	"backend/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(pswd string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func TestRegister(t *testing.T) {
	t.Run("success creates user", func(t *testing.T) {
		userDAO := new(MockUserDAO)
		svc := NewAuthService(userDAO)

		userDAO.On("FindByEmail", "new@test.com").Return(nil, nil)
		userDAO.On("Create", mock.MatchedBy(func(u *domain.User) bool {
			return u.Email == "new@test.com" && u.Name == "Test User"
		})).Return(nil)

		user, err := svc.Register(RegisterInput{
			Name:     "Test User",
			Email:    "new@test.com",
			Password: "securePass123",
		})
		require.NoError(t, err)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, "client", user.Role)
		assert.NotEmpty(t, user.PasswordHash)
		userDAO.AssertExpectations(t)
	})

	t.Run("short password returns ErrInvalidInput", func(t *testing.T) {
		svc := NewAuthService(new(MockUserDAO))
		_, err := svc.Register(RegisterInput{
			Name:     "User",
			Email:    "u@test.com",
			Password: "short",
		})
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("duplicate email returns ErrInvalidInput", func(t *testing.T) {
		userDAO := new(MockUserDAO)
		svc := NewAuthService(userDAO)

		existing := &domain.User{Email: "dup@test.com"}
		userDAO.On("FindByEmail", "dup@test.com").Return(existing, nil)

		_, err := svc.Register(RegisterInput{
			Name:     "Dup",
			Email:    "dup@test.com",
			Password: "securePass123",
		})
		assert.ErrorIs(t, err, ErrInvalidInput)
		userDAO.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	hash := hashPassword("myPassword")

	t.Run("success returns token", func(t *testing.T) {
		userDAO := new(MockUserDAO)
		svc := NewAuthService(userDAO)

		userDAO.On("FindByEmail", "a@test.com").Return(&domain.User{
			ID:           1,
			Email:        "a@test.com",
			PasswordHash: hash,
		}, nil)

		token, _, err := svc.Login(LoginInput{Email: "a@test.com", Password: "myPassword"})
		require.NoError(t, err)
		assert.NotEmpty(t, token)
		userDAO.AssertExpectations(t)
	})

	t.Run("wrong password returns ErrUnauthorized", func(t *testing.T) {
		userDAO := new(MockUserDAO)
		svc := NewAuthService(userDAO)

		userDAO.On("FindByEmail", "u@test.com").Return(&domain.User{
			Email:        "u@test.com",
			PasswordHash: hash,
		}, nil)

		_, _, err := svc.Login(LoginInput{Email: "u@test.com", Password: "wrongPassword"})
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("unknown email returns ErrUnauthorized", func(t *testing.T) {
		userDAO := new(MockUserDAO)
		svc := NewAuthService(userDAO)

		userDAO.On("FindByEmail", "no@test.com").Return(nil, errors.New("not found"))

		_, _, err := svc.Login(LoginInput{Email: "no@test.com", Password: "any"})
		assert.ErrorIs(t, err, ErrUnauthorized)
	})
}
