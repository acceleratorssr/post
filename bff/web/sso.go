package web

import (
	"github.com/gin-gonic/gin"
	ssov1 "post/api/proto/gen/sso/v1"
	"post/pkg/gin-extra"
)

type SSOHandler struct {
	sso ssov1.AuthServiceClient
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenReq struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

func (s *SSOHandler) RegisterRoutes(engine *gin.Engine, mw gin.HandlerFunc) {
	ssoGroup := engine.Group("/sso")
	ssoGroup.POST("/login", s.Login)

	ssoGroup.POST("/refresh", mw, gin_extra.WrapWithReq[RefreshTokenReq](s.Refresh)) // 由 web jwtAOP 检查token
	ssoGroup.POST("/logout", mw, gin_extra.WrapNilReq(s.Logout))
}

func (s *SSOHandler) Login(ctx *gin.Context) {
	// todo 添加bloom，快速判断用户是否存在
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		gin_extra.FailWithMessage(ctx, gin_extra.InvalidArgument, err.Error())
		return
	}

	login, err := s.sso.Login(ctx, &ssov1.LoginRequest{
		Username:  req.Username,
		Password:  req.Password,
		UserAgent: ctx.Request.UserAgent(),
		Code:      req.Code,
	})
	if err != nil {
		gin_extra.FailWithError(ctx, err, gin_extra.NotFound)
	}

	gin_extra.OKWithData(ctx, LoginResp{
		AccessToken:  login.AccessToken,
		RefreshToken: login.RefreshToken,
	})
}

// Refresh 安全性：考虑同一token连续n次请求刷新，触发2fa验证
func (s *SSOHandler) Refresh(ctx *gin.Context, request RefreshTokenReq) (*gin_extra.Response, error) {
	token, err := s.sso.RefreshToken(ctx, &ssov1.RefreshTokenRequest{
		RefreshToken: ctx.Value("token").(string),
		UserInfo: &ssov1.UserInfo{
			Uid:      ctx.Value("uid").(uint64),
			Username: request.Username,
			Nickname: request.Nickname,
		},
	})
	if err != nil {
		return &gin_extra.Response{
			Code: gin_extra.Unauthenticated,
			Msg:  "请重新登录",
		}, err
	}

	return &gin_extra.Response{
		Data: token.AccessToken,
		Msg:  "刷新成功",
	}, nil
}

func (s *SSOHandler) Logout(ctx *gin.Context) (*gin_extra.Response, error) {
	_, err := s.sso.Logout(ctx, &ssov1.LogoutRequest{
		ExpiredAt:    ctx.Value("x_exp").(int64),
		RefreshToken: ctx.Value("x_token").(string),
	})
	if err != nil {
		return &gin_extra.Response{
			Code: gin_extra.Unauthenticated,
			Msg:  "退出失败",
		}, err
	}

	_, err = s.sso.Logout(ctx, &ssov1.LogoutRequest{
		ExpiredAt:    ctx.Value("exp").(int64),
		RefreshToken: ctx.Value("token").(string),
	})
	if err != nil {
		return &gin_extra.Response{
			Code: gin_extra.Unauthenticated,
			Msg:  "退出失败",
		}, err
	}

	return &gin_extra.Response{
		Msg: "退出成功",
	}, nil
}

func NewSSOHandler(sso ssov1.AuthServiceClient) *SSOHandler {
	return &SSOHandler{
		sso: sso,
	}
}
