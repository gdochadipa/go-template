package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/user/go-templates/template-grpc-ddd/internal/core/domain"
	"github.com/user/go-templates/template-grpc-ddd/internal/core/port"
)

type UserService struct {
	repo   port.UserRepository
	logger *slog.Logger
}

func NewUserService(repo port.UserRepository, logger *slog.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	s.logger.InfoContext(ctx, "fetching user", "id", id)
	return s.repo.Get(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	s.logger.InfoContext(ctx, "creating user", "email", user.Email)

	if user.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	user.ID = uuid.New().String()
	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
