package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "用户不存在")
		}
		return nil, status.Errorf(codes.Unknown, "未知错误")
	}
	return &userv1.GetUserInfoByUsernameResponse{
		User: u.ToDTO(data),
	}, nil
}

func (u *UserServiceServer) UpdateUser(ctx context.Context, request *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	err := u.svc.UpdateUser(ctx, &domain.UserInfo{
		Nickname: request.GetUserInfo().GetNickname(),
	})
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

func NewUserServiceServer(svc service.UserService, ssoGrpcClient ssov1.AuthServiceClient) *UserServiceServer {
	return &UserServiceServer{
		svc:           svc,
		ssoGrpcClient: ssoGrpcClient,
	}
}
