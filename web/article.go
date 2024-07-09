package web

import (
	"github.com/gin-gonic/gin"
	"post/domain"
	"post/service"
	"post/user"
	"post/utils"
)

//var _ post.handler = (*ArticleHandler)(nil)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc service.ArticleService
}

func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}

func (a *ArticleHandler) RegisterRoutes(s *gin.Engine) {
	articles := s.Group("/articles")
	articles.POST("/save", a.Save)
	articles.POST("/publish", a.Publish)
	articles.POST("/withdraw", a.Withdraw)
}

func (a *ArticleHandler) Publish(c *gin.Context) {
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}
	// TODO 检查输入

	claims := a.getUserInfo(c)

	id, err := a.svc.Publish(c, req.toDomain(claims.Id, claims.Name))
	if err != nil {
		utils.Fail(domain.ErrSystem, err.Error(), "系统错误", c)
		// log
		return
	}
	utils.OK(id, "success", c)
}

func (a *ArticleHandler) Save(c *gin.Context) {
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	} //需要使用ShouldBindJSON，如果使用bind则获取不到值
	// TODO 检查输入

	claims := a.getUserInfo(c)

	id, err := a.svc.Save(c, req.toDomain(claims.Id, claims.Name))
	if err != nil {
		utils.FailWithMessage(domain.ErrSystem, "系统错误", c)
		// log
		return
	}

	utils.OK(id, "success", c)
}

func (a *ArticleHandler) Withdraw(c *gin.Context) {
	var req ReqOnlyWithID
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	} //需要使用ShouldBindJSON，如果使用bind则获取不到值
	// TODO 检查输入

	claims := a.getUserInfo(c)
	err := a.svc.Withdraw(c, req.toDomain(claims.Id, claims.Name))
	if err != nil {
		utils.FailWithMessage(domain.ErrSystem, "系统错误", c)
		// log
		return
	}

	utils.OKWithMessage("success", c)
}

func (a *ArticleHandler) getUserInfo(c *gin.Context) *user.ClaimsUser {
	claim := c.MustGet("userClaims")
	claims, ok := claim.(*user.ClaimsUser)
	if !ok {
		utils.FailWithMessage(domain.ErrSystem, "系统错误", c)
		// log
		return nil
	}
	return claims
}
