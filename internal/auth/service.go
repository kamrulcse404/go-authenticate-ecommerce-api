package auth

import (
	"context"
	"ecommerce-api/internal/platform/security"
	"ecommerce-api/internal/user"
	"errors"
	"strings"
)

var (
	ErrNameRequired     = errors.New("name is required")
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*user.User, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if req.Name == "" {
		return nil, ErrNameRequired
	}

	if req.Email == "" {
		return nil, ErrEmailRequired
	}

	if req.Password == "" {
		return nil, ErrPasswordRequired
	}

	if len(req.Password) < 8 {
		return nil, ErrPasswordTooShort
	}

	passwordHash, err := security.HashPassword(req.Password)

	if err != nil {
		return nil, err
	}

	return s.repo.CreateUser(ctx, req.Name, req.Email, passwordHash)
}
