package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	ssov1 "post/api/proto/gen/sso/v1"
	"post/sso/config"
	"post/sso/domain"
	"post/sso/repository"
	"post/sso/service"
	"time"
)

type AuthServiceServer struct {
	ssov1.UnimplementedAuthServiceServer
	issuer     string
	expiration int64 // 用户和TOTP绑定关系的缓存时间
	svc        service.AuthUserService
	cache      repository.SSOCache
	jwtSvc     domain.AuthService
	info       *config.Info
}

func (a *AuthServiceServer) GetPublicKey(ctx context.Context, request *ssov1.PublicKeyRequest) (*ssov1.PublicKeyResponse, error) {
	key := a.jwtSvc.GetPublicKey(ctx)
	if key == "" {
		return nil, status.Errorf(codes.Internal, "SSO 读取公钥失败")
	}
	return &ssov1.PublicKeyResponse{
		PublicKey: key,
	}, nil
}

func (a *AuthServiceServer) BindTotp(ctx context.Context, request *ssov1.BindTotpRequest) (*ssov1.BindTotpResponse, error) {
	if a.svc.FindUsernameExist(ctx, request.GetUsername()) {
		return nil, status.Errorf(codes.AlreadyExists, "SSO 用户已存在")
	}

	key, url, err := a.GenTotpSecret(a.issuer, request.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO TOTP绑定失败: %s", err))
	}

	err = a.cache.SetString(ctx, request.GetUsername(), key, 10*time.Minute)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 缓存用户绑定关系失败: %s", err))
	}

	return &ssov1.BindTotpResponse{
		QRUrl: url,
	}, nil
}

func (a *AuthServiceServer) Register(ctx context.Context, request *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	now := time.Now().UnixMilli()
	secretKey, err := a.cache.GetString(ctx, request.GetUserInfo().GetUsername())
	if err == redis.Nil {
		return nil, status.Errorf(codes.Unauthenticated, fmt.Sprintf("SSO 绑定密钥超时，请重新注册: %s", err))
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 获取用户绑定关系失败: %s", err))
	}

	if !a.validateTOTP(secretKey, request.GetCode()) {
		return nil, status.Errorf(codes.Unauthenticated, "SSO 2FA验证码错误")
	}

	err = a.svc.SaveUser(ctx, &domain.User{
		Username:   request.GetUserInfo().GetUsername(),
		Nickname:   request.GetUserInfo().GetNickname(),
		Password:   request.GetPassword(),
		TotpSecret: secretKey,
		UserAgent:  request.GetUserAgent(),
	}, now)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 保存用户信息失败: %s", err))
	}

	accessToken, err := a.jwtSvc.GenerateAccessToken(ctx, &domain.JwtPayload{
		Username: request.GetUserInfo().GetUsername(),
		NickName: request.GetUserInfo().GetNickname(),
		Ctime:    now,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 生成 access token 失败: %s", err))
	}

	refreshToken, err := a.jwtSvc.GenerateRefreshToken(ctx, &domain.JwtPayload{
		Ctime: now,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 生成 refresh token 失败: %s", err))
	}

	return &ssov1.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AuthServiceServer) UpdatePassword(ctx context.Context, request *ssov1.UpdatePasswordRequest) (*ssov1.UpdatePasswordResponse, error) {
	totpSecret, err := a.svc.GetTotpSecret(ctx, request.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("SSO 用户查找失败: %s", err))
	}

	if !a.validateTOTP(totpSecret, request.GetCode()) {
		return nil, status.Errorf(codes.Unauthenticated, "SSO 2FA验证码错误")
	}

	err = a.svc.SaveUser(ctx, &domain.User{
		Password: request.GetPassword(),
	}, time.Now().UnixMilli())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 系统错误，更新密码失败: %s", err))
	}

	return &ssov1.UpdatePasswordResponse{}, nil
}

func (a *AuthServiceServer) Logout(ctx context.Context, request *ssov1.LogoutRequest) (*ssov1.LogoutResponse, error) {
	handleToken := func(tokenStr string, maxExpires int64) error {
		token, err := a.jwtSvc.ValidateToken(ctx, tokenStr)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return nil
			}
			return status.Errorf(codes.Unauthenticated, fmt.Sprintf("SSO token 验证失败: %v", err))
		}
		return a.cache.SetString(ctx, tokenStr, "",
			time.Duration(max(0, maxExpires-(time.Now().UnixMilli()-token.Ctime))))
	}

	if err := handleToken(request.GetRefreshToken(), a.info.Config.Jwt.LongExpires); err != nil {
		return nil, err
	}

	return &ssov1.LogoutResponse{}, nil
}

func (a *AuthServiceServer) RefreshToken(ctx context.Context, request *ssov1.RefreshTokenRequest) (*ssov1.RefreshTokenResponse, error) {
	_, err := a.jwtSvc.ValidateToken(ctx, request.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, fmt.Sprintf("SSO refreshToken验证失败: %v", err))
	}

	token, err := a.jwtSvc.GenerateAccessToken(ctx, &domain.JwtPayload{
		Username: request.GetUserInfo().GetUsername(),
		NickName: request.GetUserInfo().GetNickname(),
		Ctime:    time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 生成短token失败: %v", err))
	}

	return &ssov1.RefreshTokenResponse{
		AccessToken: token,
	}, nil
}

// Login 当出现 UserAgent 不一致的情况，则会要求用户提交验证码
// 正常逻辑下，当 UserAgent 不一致时，第一次调用没有验证码，第二次重复调用该方法即可
func (a *AuthServiceServer) Login(ctx context.Context, request *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	user, err := a.svc.GetInfoByUsername(ctx, request.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "SSO 未找到对应用户")
	}

	if !a.CheckPasswords(user.Password, request.GetPassword()) {
		return nil, status.Errorf(codes.Unauthenticated, "SSO 用户或密码错误")
	}

	if user.UserAgent != request.UserAgent {
		if request.GetCode() == "" {
			return nil, status.Errorf(codes.Unauthenticated, "SSO 风险行为，请输入2FA验证码")
		}
		if !a.validateTOTP(user.TotpSecret, request.GetCode()) {
			return nil, status.Errorf(codes.Unauthenticated, "SSO 2FA验证码错误")
		}
	}

	jwtPayload := &domain.JwtPayload{
		Username: user.Username,
		NickName: user.Nickname,
		Ctime:    time.Now().UnixMilli(),
	}
	accessToken, err := a.jwtSvc.GenerateAccessToken(ctx, jwtPayload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 生成短token失败: %v", err))
	}
	refreshToken, err := a.jwtSvc.GenerateRefreshToken(ctx, jwtPayload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 生成长token失败: %v", err))
	}

	return &ssov1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateToken
// <--！ 已弃用 ！-->
// 为防止每次客户端访问其他服务时都要走 SSO 验证，
// 故使用非对称加密的jwt token，由其他服务启动后自动请求SSO服务获取公钥后，自行验证；
//func (a *AuthServiceServer) ValidateToken(ctx context.Context, request *ssov1.ValidateTokenRequest) (*ssov1.ValidateTokenResponse, error) {
//	token, err := a.jwtSvc.ValidateToken(ctx, request.GetToken())
//	if err != nil {
//		return nil, status.Errorf(codes.Unauthenticated, fmt.Sprintf("SSO token验证失败: %v", err))
//	}
//
//	return &ssov1.ValidateTokenResponse{
//		Valid: true,
//		JwtPayload: &ssov1.JwtPayload{
//			UserInfo: &ssov1.UserInfo{
//				Username: token.Username,
//				Nickname: token.NickName,
//			},
//		},
//	}, nil
//}

func (a *AuthServiceServer) RegisterServer(server *grpc.Server) {
	ssov1.RegisterAuthServiceServer(server, a)
}

// GenTotpSecret resp: (secretKey string, sRL string, err error)
func (a *AuthServiceServer) GenTotpSecret(issuer, username string) (string, string, error) {
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: username,
	})
	if err != nil {
		return "", "", err
	}

	return secret.Secret(), secret.URL(), nil
}

func (a *AuthServiceServer) validateTOTP(totpSecret, code string) bool {
	return totp.Validate(code, totpSecret)
}

func (a *AuthServiceServer) HashAndSalt(pwd string) string {
	bytePwd := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

func (a *AuthServiceServer) CheckPasswords(hashedPwd string, rePwd string) bool {
	byteHash := []byte(hashedPwd)
	byteRePwd := []byte(rePwd)

	err := bcrypt.CompareHashAndPassword(byteHash, byteRePwd)
	if err != nil {
		return false
	}
	return true
}

func NewSSOServiceServer(svc service.AuthUserService, info *config.Info, jwtSvc domain.AuthService, cache repository.SSOCache) *AuthServiceServer {
	return &AuthServiceServer{
		issuer: "歪比八不",
		svc:    svc,
		jwtSvc: jwtSvc,
		cache:  cache,
		info:   info,
	}
}
