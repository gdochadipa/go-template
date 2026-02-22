package user

import (
	"context"

	userv1 "github.com/user/go-templates/template-grpc-sdk/gen/go/user/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	userv1.UnimplementedUserServiceServer
	logger *zap.Logger
}

func NewService(logger *zap.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	s.logger.Info("fetching user", zap.String("id", req.GetId()))

	// Mock implementation
	return &userv1.GetUserResponse{
		User: &userv1.User{
			Id:    req.GetId(),
			Name:  "John Doe",
			Email: "john@example.com",
		},
	}, nil
}

func (s *Service) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	s.logger.Info("creating user", zap.String("email", req.GetEmail()))

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	return &userv1.CreateUserResponse{
		User: &userv1.User{
			Id:    "new-uuid",
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	}, nil
}
