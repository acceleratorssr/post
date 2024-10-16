package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	noopv1 "post/api/proto/gen/no-op/v1"
	ssov1 "post/api/proto/gen/sso/v1"
	userv1 "post/api/proto/gen/user/v1"
	ch "post/pkg/grpc-extra/balancer/CHBL"
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
	resp, err := ch.RegisterWithKey(ctx, request.GetUser().GetUsername(), &ssov1.RegisterRequest{
		UserInfo: &ssov1.UserInfo{
			Username: request.GetUser().GetUsername(),
			Nickname: request.GetUser().GetNickname(),
		},
		Password:  request.GetUser().GetPassword(),
		UserAgent: request.GetUser().GetUserAgent(),
		Code:      request.GetCode(),
	}, u.ssoGrpcClient.Register)
	if err != nil {
		return nil, err
	}

	user := u.ToDomain(request.GetUser())
	user.UID = resp.Uid
	err = u.svc.CreateUser(ctx, user) // todo 考虑异步化

	if err != nil {
		return nil, status.Errorf(codes.Internal, "用户注册成功，保存信息失败")
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
		return nil, status.Errorf(codes.Unknown, "获取信息失败: %v", err)
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
		return nil, status.Errorf(codes.Internal, "SSO 生成短token失败: %v", err)
	}
	return &userv1.UpdateUserResponse{}, nil
}

func (u *UserServiceServer) ToDomain(user *userv1.User) *domain.User {
	return &domain.User{
		UID:      user.Uid,
		Username: user.Username,
		Nickname: user.Nickname,
	}
}

func (u *UserServiceServer) ToDTO(user *domain.User) *userv1.User {
	return &userv1.User{
		Uid:      user.UID,
		Username: user.Username,
		Nickname: user.Nickname,
	}
}

func NewUserServiceServer(svc service.UserService, ssoGrpcClient ssov1.AuthServiceClient) *UserServiceServer {
	// 使 gRPC 立刻与服务节点建立 HTTP2.0 长连接
	// 防止第一次请求时，服务节点尚未完全就绪
	// 导致影响负载均衡中，通过服务节点数量计算的相关逻辑
	// 如果不使用该 noop 请求，则第一次请求数据的分布会和之后的不一样
	_, _ = ch.RegisterWithKey(context.Background(), "",
		&noopv1.NoOpRequest{}, ssoGrpcClient.NoOp)

	return &UserServiceServer{
		svc:           svc,
		ssoGrpcClient: ssoGrpcClient,
	}
}
