package services

import (
	"fmt"

	db "backend/dao"
	"backend/domain"
	"backend/utils"
)

// ---------------------------------------------------------------------------
// Interface
// ---------------------------------------------------------------------------

type AuthServiceInterface interface {
	Register(input RegisterInput) (*domain.User, error)
	Login(input LoginInput) (string, error)
}

// ---------------------------------------------------------------------------
// Input types
// ---------------------------------------------------------------------------

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

// ---------------------------------------------------------------------------
// Concrete service
// ---------------------------------------------------------------------------

type AuthService struct {
	userDAO db.UserDAO
}

func NewAuthService(userDAO db.UserDAO) *AuthService {
	return &AuthService{userDAO: userDAO}
}

func (s *AuthService) Register(input RegisterInput) (*domain.User, error) {
	if len(input.Password) < 8 {
		return nil, fmt.Errorf("%w: password must be at least 8 characters", ErrInvalidInput)
	}

	existing, err := s.userDAO.FindByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("checking existing user: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: email is already registered", ErrInvalidInput)
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hash,
		Role:         "client",
	}

	if err := s.userDAO.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(input LoginInput) (string, error) {
	user, err := s.userDAO.FindByEmail(input.Email)
	if err != nil {
		return "", ErrUnauthorized
	}

	if !utils.CheckPassword(input.Password, user.PasswordHash) {
		return "", ErrUnauthorized
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// Compile-time interface check.
var _ AuthServiceInterface = (*AuthService)(nil)
