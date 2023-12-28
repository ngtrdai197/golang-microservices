package transport

import (
	apiV1 "codebase/api/user/v1"
	"context"
)

type UserTransport interface {
	apiV1.UserServer
	CreateUser(ctx context.Context, request *apiV1.CreateUserRequest) (*apiV1.CreateUserResponse, error)
}

type userTransport struct {
	apiV1.UnimplementedUserServer
}

func NewUserTransport() UserTransport {
	return &userTransport{}
}

func (s *userTransport) CreateUser(ctx context.Context, request *apiV1.CreateUserRequest) (*apiV1.CreateUserResponse, error) {
	return &apiV1.CreateUserResponse{
		Name:  request.Name,
		Email: request.Email,
		Phone: request.Phone,
	}, nil
}
