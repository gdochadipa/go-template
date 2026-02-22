package port

import (
	"context"

	"github.com/user/go-templates/template-grpc-ddd/internal/core/domain"
)

// UserService defines the input port for user operations.
// This is what the adapter (gRPC handler) will call.
type UserService interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}
