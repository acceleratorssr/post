package web

import (
	"github.com/gin-gonic/gin"
	ssov1 "post/api/proto/gen/sso/v1"
	"post/pkg/gin_ex"
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

	ssoGroup.POST("/refresh", mw, gin_ex.WrapWithReq[RefreshTokenReq](s.Refresh)) // 由 web jwtAOP 检查token
	ssoGroup.POST("/logout", mw, gin_ex.WrapWithReq[LoginReq](s.Logout))
}

func (s *SSOHandler) Login(ctx *gin.Context) {
	// todo 添加bloom，快速判断用户是否存在
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		gin_ex.FailWithMessage(ctx, gin_ex.InvalidArgument, err.Error())
		return
	}

	login, err := s.sso.Login(ctx, &ssov1.LoginRequest{
		Username:  req.Username,
		Password:  req.Password,
		UserAgent: ctx.Request.UserAgent(),
		Code:      req.Code,
	})
	if err != nil {
		gin_ex.FailWithError(ctx, err, gin_ex.NotFound)
	}

	gin_ex.OKWithData(ctx, LoginResp{
		AccessToken:  login.AccessToken,
		RefreshToken: login.RefreshToken,
	})
}

func (s *SSOHandler) Refresh(ctx *gin.Context, request RefreshTokenReq) (*gin_ex.Response, error) {
	token, err := s.sso.RefreshToken(ctx, &ssov1.RefreshTokenRequest{
		RefreshToken: ctx.Value("token").(string),
		UserInfo: &ssov1.UserInfo{
			Username: request.Username,
			Nickname: request.Nickname,
		},
	})
	if err != nil {
		return &gin_ex.Response{
			Code: gin_ex.Unauthenticated,
			Msg:  "请重新登录",
		}, err
	}

	return &gin_ex.Response{
		Data: token.AccessToken,
		Msg:  "刷新成功",
	}, nil
}

func (s *SSOHandler) Logout(ctx *gin.Context, request LoginReq) (*gin_ex.Response, error) {
	_, err := s.sso.Logout(ctx, &ssov1.LogoutRequest{
		RefreshToken: ctx.Value("token").(string),
	})
	if err != nil {
		return &gin_ex.Response{
			Code: gin_ex.Unauthenticated,
			Msg:  "退出失败",
		}, err
	}

	return &gin_ex.Response{
		Msg: "退出成功",
	}, nil
}

func NewSSOHandler(sso ssov1.AuthServiceClient) *SSOHandler {
	return &SSOHandler{
		sso: sso,
	}
}
