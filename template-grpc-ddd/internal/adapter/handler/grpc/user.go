package handler

import (
	"context"

	userv1 "github.com/user/go-templates/template-grpc-ddd/gen/go/user/v1"
	"github.com/user/go-templates/template-grpc-ddd/internal/core/domain"
	"github.com/user/go-templates/template-grpc-ddd/internal/core/port"
)

type UserHandler struct {
	userv1.UnimplementedUserServiceServer
	svc port.UserService
}

func NewUserHandler(svc port.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := h.svc.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserResponse{
		User: mapDomainToProto(user),
	}, nil
}

func (h *UserHandler) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	user := &domain.User{
		Name:  req.Name,
		Email: req.Email,
	}

	createdUser, err := h.svc.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &userv1.CreateUserResponse{
		User: mapDomainToProto(createdUser),
	}, nil
}

func mapDomainToProto(u *domain.User) *userv1.User {
	if u == nil {
		return nil
	}
	return &userv1.User{
		Id:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
