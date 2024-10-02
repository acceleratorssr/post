package web

import (
	"github.com/gin-gonic/gin"
	ssov1 "post/api/proto/gen/sso/v1"
	userv1 "post/api/proto/gen/user/v1"
	"post/pkg/gin_ex"
)

type UserHandler struct {
	svc userv1.UserServiceClient
	sso ssov1.AuthServiceClient
}

type Bind2FAReq struct {
	Username string `json:"username"`
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Code     string `json:"code"`
}

type RegisterResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u *UserHandler) RegisterRoutes(engine *gin.Engine, mw gin.HandlerFunc) {
	userGroup := engine.Group("/user")
	userGroup.POST("/bind2FA", gin_ex.WrapWithReq[Bind2FAReq](u.Bind2FA))
	userGroup.POST("/register", u.Register)
}

func (u *UserHandler) Register(ctx *gin.Context) {
	var req RegisterReq
	if err := ctx.Bind(&req); err != nil {
		gin_ex.FailWithMessage(ctx, gin_ex.InvalidArgument, err.Error())
		return
	}

	agent := ctx.GetHeader("User-Agent")
	user, err := u.svc.CreateUser(ctx, &userv1.CreateUserRequest{
		User: &userv1.User{
			Username:  req.Username,
			Password:  req.Password,
			Nickname:  req.Nickname,
			UserAgent: agent,
		},
		Code: req.Code,
	})
	if err != nil {
		gin_ex.FailWithMessage(ctx, gin_ex.System, err.Error())
		return
	}

	gin_ex.OKWithDataAndMsg(ctx, RegisterResp{
		AccessToken:  user.GetAccessToken(),
		RefreshToken: user.GetRefreshToken(),
	}, "注册成功")
}

func (u *UserHandler) Bind2FA(ctx *gin.Context, req Bind2FAReq) (*gin_ex.Response, error) {
	totp, err := u.sso.BindTotp(ctx, &ssov1.BindTotpRequest{Username: req.Username})
	if err != nil {
		return nil, err
	}
	return &gin_ex.Response{
		Data: totp.QRUrl,
		Msg:  "请扫描二维码启用2FA，有效期为10分钟",
		Code: gin_ex.OK,
	}, nil
}

func NewUserHandler(svc userv1.UserServiceClient, sso ssov1.AuthServiceClient) *UserHandler {
	return &UserHandler{
		svc: svc,
		sso: sso,
	}
}
