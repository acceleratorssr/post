package grpc

import (
	"context"
	"google.golang.org/grpc"
	"post/api/proto/gen/common"
	ssov1 "post/api/proto/gen/sso/v1"
	userv1 "post/api/proto/gen/user/v1"
	"post/user/domain"
	"post/user/service"
)

type UserServiceServer struct {
	userv1.UnimplementedUserServiceServer
	svc           service.UserService
	ssoGrpcClient ssov1.AuthServiceClient
}

func NewUserServiceServer(svc service.UserService) *UserServiceServer {
	return &UserServiceServer{
		svc: svc,
	}
}

func (u *UserServiceServer) Register(server *grpc.Server) {
	userv1.RegisterUserServiceServer(server, u)
}

func (u *UserServiceServer) CreateUser(ctx context.Context, request *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	// sso
	_, err := u.ssoGrpcClient.Register(ctx, &ssov1.RegisterRequest{
		Username:  request.User.Username,
		Password:  request.User.Password,
		UserAgent: request.User.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	err = u.svc.CreateUser(ctx, u.ToDomain(request.GetUser()))

	if err != nil {
		return nil, err
	}
	return &userv1.CreateUserResponse{}, nil
}

func (u *UserServiceServer) GetUserInfoByUsername(ctx context.Context, request *userv1.GetUserInfoByUsernameRequest) (*userv1.GetUserInfoByUsernameResponse, error) {
	data, err := u.svc.GetUserInfoByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}
	return &userv1.GetUserInfoByUsernameResponse{
		User: u.ToDTO(data),
	}, nil
}

func (u *UserServiceServer) UpdateUser(ctx context.Context, request *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	err := u.svc.UpdateUser(ctx, u.ToDomain(request.GetUser()))
	if err != nil {
		return nil, err
	}
	return &userv1.UpdateUserResponse{}, nil
}

func (u *UserServiceServer) ToDomain(user *userv1.User) *domain.User {
	return &domain.User{
		Username:    user.Username,
		Nickname:    user.Nickname,
		Permissions: int(user.Permissions),
	}
}

func (u *UserServiceServer) ToDTO(user *domain.User) *userv1.User {
	return &userv1.User{
		Username:    user.Username,
		Nickname:    user.Nickname,
		Permissions: common.Permissions(user.Permissions),
	}
}
