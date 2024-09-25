package grpc

import (
	"context"
	"google.golang.org/grpc"
	"post/api/proto/gen/common"
	userv1 "post/api/proto/gen/user/v1"
	"post/user/domain"
	"post/user/service"
)

type UserServiceServer struct {
	userv1.UnimplementedUserServiceServer
	svc service.UserService
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
	err := u.svc.CreateUser(ctx, u.ToDomain(request.GetUser()))
	if err != nil {
		return nil, err
	}
	return &userv1.CreateUserResponse{
		Message: "success",
	}, nil
}

func (u *UserServiceServer) GetUserInfoByUsername(ctx context.Context, request *userv1.GetUserInfoByUsernameRequest) (*userv1.GetUserInfoByUsernameResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserServiceServer) UpdateUser(ctx context.Context, request *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserServiceServer) ToDomain(user *userv1.User) *domain.User {
	return &domain.User{
		ID:          user.Id,
		Username:    user.Username,
		Password:    user.Password,
		Nickname:    user.Nickname,
		Permissions: int(user.Permissions),
	}
}

func (u *UserServiceServer) ToDTO(user *domain.User) *userv1.User {
	return &userv1.User{
		Id:          user.ID,
		Username:    user.Username,
		Password:    user.Password,
		Nickname:    user.Nickname,
		Permissions: common.Permissions(user.Permissions),
	}
}
