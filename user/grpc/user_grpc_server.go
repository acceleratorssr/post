package grpc

import (
	"context"
	"google.golang.org/grpc"
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

func (u *UserServiceServer) Register(server *grpc.Server) {
	userv1.RegisterUserServiceServer(server, u)
}

func (u *UserServiceServer) CreateUser(ctx context.Context, request *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	// sso
	resp, err := u.ssoGrpcClient.Register(ctx, &ssov1.RegisterRequest{
		UserInfo: &ssov1.UserInfo{
			Username: request.GetUser().GetUsername(),
			Nickname: request.GetUser().GetNickname(),
		},
		Password:  request.GetUser().GetPassword(),
		UserAgent: request.GetUser().GetUserAgent(),
		Code:      request.GetCode(),
	})
	if err != nil {
		return nil, err // 此处err为 SSO服务 返回的
	}

	err = u.svc.CreateUser(ctx, u.ToDomain(request.GetUser()))

	if err != nil {
		return nil, err
	}
	return &userv1.CreateUserResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
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
		Username: user.Username,
		Nickname: user.Nickname,
	}
}

func (u *UserServiceServer) ToDTO(user *domain.User) *userv1.User {
	return &userv1.User{
		Username: user.Username,
		Nickname: user.Nickname,
	}
}

func NewUserServiceServer(svc service.UserService) *UserServiceServer {
	return &UserServiceServer{
		svc: svc,
	}
}
