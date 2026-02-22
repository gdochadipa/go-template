package port

import (
	"context"

	"github.com/user/go-templates/template-grpc-ddd/internal/core/domain"
)

// UserRepository defines the output port for persistence.
// This is what the application core will use to save/retrieve data.
type UserRepository interface {
	Get(ctx context.Context, id string) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) error
}
