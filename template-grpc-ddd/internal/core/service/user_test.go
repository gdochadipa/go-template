package service

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/user/go-templates/template-grpc-ddd/internal/core/domain"
)

// MockUserRepository is a manual mock for the port.UserRepository interface
type MockUserRepository struct {
	GetFunc  func(ctx context.Context, id string) (*domain.User, error)
	SaveFunc func(ctx context.Context, user *domain.User) error
}

func (m *MockUserRepository) Get(ctx context.Context, id string) (*domain.User, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, errors.New("unimplemented")
}

func (m *MockUserRepository) Save(ctx context.Context, user *domain.User) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, user)
	}
	return nil
}

func TestUserService_CreateUser(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	tests := []struct {
		name          string
		inputUser     *domain.User
		mockRepo      *MockUserRepository
		expectedError bool
	}{
		{
			name:      "Success",
			inputUser: &domain.User{Name: "John", Email: "john@example.com"},
			mockRepo: &MockUserRepository{
				SaveFunc: func(ctx context.Context, user *domain.User) error {
					if user.ID == "" {
						return errors.New("ID should be generated")
					}
					return nil
				},
			},
			expectedError: false,
		},
		{
			name:          "MissingName",
			inputUser:     &domain.User{Name: "", Email: "john@example.com"},
			mockRepo:      &MockUserRepository{},
			expectedError: true,
		},
		{
			name:      "RepoError",
			inputUser: &domain.User{Name: "John", Email: "john@example.com"},
			mockRepo: &MockUserRepository{
				SaveFunc: func(ctx context.Context, user *domain.User) error {
					return errors.New("db error")
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewUserService(tt.mockRepo, logger)
			_, err := svc.CreateUser(context.Background(), tt.inputUser)

			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}
