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
	info       *config.Info
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

	jwtPayload := &domain.JwtPayload{
		Username: request.GetUserInfo().GetUsername(),
		NickName: request.GetUserInfo().GetNickname(),
		Ctime:    now,
	}
	accessToken, err := a.generateAccessToken(jwtPayload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 生成 access token 失败: %s", err))
	}
	refreshToken, err := a.generateRefreshToken(jwtPayload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 生成 refresh token 失败: %s", err))
	}

	return &ssov1.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// UpdateInfo 注：SSO服务保存的用户数据是 凭证 + 必要用户信息（放入jwt），剩下的信息保存在user，将采用kafka异步消费更新数据
// todo 将此处的存储逻辑改为消费者负责异步复制 数据；
// todo 因为保证sso写入需要成功，故此处直接重新生成长短token返回给user
func (a *AuthServiceServer) UpdateInfo(ctx context.Context, request *ssov1.UpdateInfoRequest) (*ssov1.UpdateInfoResponse, error) {
	totpSecret, err := a.svc.GetTotpSecret(ctx, request.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("SSO 查找失败: %s", err))
	}

	if !a.validateTOTP(totpSecret, request.GetCode()) {
		return nil, status.Errorf(codes.Unauthenticated, "SSO 2FA验证码错误")
	}

	err = a.svc.SaveUser(ctx, &domain.User{
		Username: request.GetUsername(),
		Nickname: request.GetNickname(),
		Password: request.GetPassword(),
	}, time.Now().UnixMilli())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 更新用户信息失败: %s", err))
	}
	return &ssov1.UpdateInfoResponse{}, nil
}

func (a *AuthServiceServer) Logout(ctx context.Context, request *ssov1.LogoutRequest) (*ssov1.LogoutResponse, error) {
	handleToken := func(tokenStr string, maxExpires int64) error {
		token, err := a.ParseToken(tokenStr)
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
	_, err := a.ParseToken(request.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, fmt.Sprintf("SSO token验证失败: %v", err))
	}

	token, err := a.generateAccessToken(&domain.JwtPayload{
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

func (a *AuthServiceServer) Login(ctx context.Context, request *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	// todo 一次查库
	user, err := a.svc.GetInfoByUsername(ctx, request.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "SSO 未找到对应用户")
	}

	if !a.CheckPasswords(user.Password, request.GetPassword()) {
		return nil, status.Errorf(codes.Unauthenticated, "SSO 用户或密码错误")
	}

	if user.UserAgent != request.UserAgent {

	}

	jwtPayload := &domain.JwtPayload{
		Username: user.Username,
		NickName: user.Nickname,
		Ctime:    time.Now().UnixMilli(),
	}
	accessToken, err := a.generateAccessToken(jwtPayload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 生成短token失败: %v", err))
	}
	refreshToken, err := a.generateRefreshToken(jwtPayload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 生成长token失败: %v", err))
	}

	return &ssov1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AuthServiceServer) ValidateToken(ctx context.Context, request *ssov1.ValidateTokenRequest) (*ssov1.ValidateTokenResponse, error) {
	token, err := a.ParseToken(request.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, fmt.Sprintf("SSO token验证失败: %v", err))
	}

	return &ssov1.ValidateTokenResponse{
		Valid: true,
		JwtPayload: &ssov1.JwtPayload{
			UserInfo: &ssov1.UserInfo{
				Username: token.Username,
				Nickname: token.NickName,
			},
		},
	}, nil
}

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

//func (a *AuthServiceServer) generateQRCode(user *domain.User) error {
//	url := user.QrcodeURL
//	qrCode, err := qr.Encode(url, qr.M, qr.Auto)
//	if err != nil {
//		return err
//	}
//
//	file, err := os.Create("qrcode.png")
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	return png.Encode(file, qrCode)
//}

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

func (a *AuthServiceServer) generateRefreshToken(user *domain.JwtPayload) (string, error) {
	claims := &domain.Claims{
		JwtPayload: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(a.info.Config.Jwt.LongExpires) * time.Hour)),
			Issuer:    a.info.Config.Jwt.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.info.Config.Jwt.Secret))
}

func (a *AuthServiceServer) generateAccessToken(user *domain.JwtPayload) (string, error) {
	claims := &domain.Claims{
		JwtPayload: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(a.info.Config.Jwt.LongExpires) * time.Hour)),
			Issuer:    a.info.Config.Jwt.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.info.Config.Jwt.Secret))
}

func (a *AuthServiceServer) ParseToken(tokenStr string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.info.Config.Jwt.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func NewSSOServiceServer(svc service.AuthUserService, info *config.Info, cache repository.SSOCache) *AuthServiceServer {
	return &AuthServiceServer{
		issuer: "歪比八不",
		svc:    svc,
		info:   info,
		cache:  cache,
	}
}
