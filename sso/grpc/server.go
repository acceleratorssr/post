package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/boombuler/barcode/qr"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"image/png"
	"os"
	"post/api/proto/gen/common"
	ssov1 "post/api/proto/gen/sso/v1"
	"post/sso/config"
	"post/sso/domain"
	"post/sso/service"
	"time"
)

type AuthServiceServer struct {
	ssov1.UnimplementedAuthServiceServer
	svc  service.AuthUserService
	info *config.Info
}

func (a *AuthServiceServer) Logout(ctx context.Context, request *ssov1.LogoutRequest) (*ssov1.LogoutResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthServiceServer) Register(ctx context.Context, request *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	// 展示2fa进行绑定
	user, err := a.createUser(request.Username)
	if err != nil {
		return nil, status.Errorf(status.Code(err), fmt.Sprintf("SSO 2FA失败: %s", err))
	}
	user.Password = a.HashAndSalt(request.Password)
	user.UserAgent = request.UserAgent
	// 保存密码等凭证

	err = a.svc.SaveUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(status.Code(err), fmt.Sprintf("SSO 保存用户信息失败: %s", err))
	}
	return &ssov1.RegisterResponse{
		QRUrl: user.QrcodeURL,
	}, nil
}

func (a *AuthServiceServer) RefreshToken(ctx context.Context, request *ssov1.RefreshTokenRequest) (*ssov1.RefreshTokenResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthServiceServer) Login(ctx context.Context, request *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	// todo 一次查库
	user, err := a.svc.GetInfoByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	if !a.CheckPasswords(user.Password, request.Password) {
		return nil, status.Errorf(codes.Unauthenticated, "SSO 用户或密码错误")
	}

	// 懒得分开了，一起放着验证
	if user.UserAgent != request.UserAgent {
		if !a.validateTOTP(user.TotpSecret, request.Code) {
			return nil, status.Errorf(codes.Unauthenticated, "SSO 2FA验证码错误")
		}
	}

	jwtPayload := &domain.JwtPayload{
		UserID:   user.ID,
		Username: user.Username,
		NickName: user.Nickname,
	}
	accessToken, err := a.generateAccessToken(ctx, jwtPayload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 生成短token失败: %v", err))
	}
	refreshToken, err := a.generateRefreshToken(ctx, jwtPayload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("SSO 生成长token失败: %v", err))
	}

	return &ssov1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AuthServiceServer) ValidateToken(ctx context.Context, request *ssov1.ValidateTokenRequest) (*ssov1.ValidateTokenResponse, error) {
	token, err := a.ParseToken(ctx, request.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, fmt.Sprintf("SSO token验证失败: %v", err))
	}

	return &ssov1.ValidateTokenResponse{
		Valid: true,
		JwtPayload: &ssov1.JwtPayload{
			Userid:      token.UserID,
			Username:    token.Username,
			Nickname:    token.NickName,
			Permissions: common.Permissions(token.Permissions),
		},
		Message: "Token is valid",
	}, nil

}

func (a *AuthServiceServer) RegisterServer(server *grpc.Server) {
	ssov1.RegisterAuthServiceServer(server, a)
}

func NewSSOServiceServer(svc service.AuthUserService, info *config.Info) *AuthServiceServer {
	return &AuthServiceServer{
		svc:  svc,
		info: info,
	}
}

func (a *AuthServiceServer) createUser(username string) (*domain.User, error) {
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "歪比八不",
		AccountName: username,
		Period:      15,
	})
	if err != nil {
		return nil, err
	}

	return &domain.User{
		Username:   username,
		TotpSecret: secret.Secret(),
		QrcodeURL:  secret.URL(),
	}, nil
}

func (a *AuthServiceServer) generateQRCode(user *domain.User) error {
	url := user.QrcodeURL
	qrCode, err := qr.Encode(url, qr.M, qr.Auto)
	if err != nil {
		return err
	}

	file, err := os.Create("qrcode.png")
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, qrCode)
}

// todo 怎么做验证的？
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

func (a *AuthServiceServer) generateRefreshToken(ctx context.Context, user *domain.JwtPayload) (string, error) {
	claims := &domain.Claims{
		JwtPayload: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(a.info.Config.Jwt.LongExpires))),
			Issuer:    a.info.Config.Jwt.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.info.Config.Jwt.Secret))
}

func (a *AuthServiceServer) generateAccessToken(ctx context.Context, user *domain.JwtPayload) (string, error) {
	claims := &domain.Claims{
		JwtPayload: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(a.info.Config.Jwt.LongExpires))),
			Issuer:    a.info.Config.Jwt.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.info.Config.Jwt.Secret))
}

func (a *AuthServiceServer) ParseToken(ctx context.Context, tokenStr string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.info.Config.Jwt.Secret), nil
	})
	if err != nil {
		// log
		return nil, err
	}
	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
