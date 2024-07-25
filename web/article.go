package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
	"post/domain"
	"post/pkg/gin_ex"
	"post/service"
	"post/user"
	"post/utils"
	"strconv"
	"time"
)

//var _ post.handler = (*ArticleHandler)(nil)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc     service.ArticleService
	like    service.LikeService
	ObjType string
}

func NewArticleHandler(svc service.ArticleService, like service.LikeService) *ArticleHandler {
	return &ArticleHandler{
		svc:     svc,
		like:    like,
		ObjType: "article",
	}
}

func (a *ArticleHandler) Test(ctx context.Context) (utils.Response, error) {
	// 注意此处传入的是context.context，而不是gin.Context
	a.svc.Save(ctx, domain.Article{
		Title:   "test",
		Content: "test",
	})
	// 复用
	span := trace.SpanFromContext(ctx)
	span.AddEvent("---50%---")

	if rand.Int31n(100)%2 == 0 {

		return utils.Response{
			Code: utils.UserInvalidInput,
			Msg:  "fail",
		}, nil
	}
	return utils.Response{
		Code: 200,
		Msg:  "ok",
	}, nil
}

func (a *ArticleHandler) RegisterRoutes(s *gin.Engine) {
	s.POST("/test",
		gin_ex.WrapNilReq(a.Test))
	articles := s.Group("/articles")
	articles.POST("/save", a.Save)
	articles.POST("/publish", a.Publish)
	articles.POST("/withdraw", a.Withdraw)
	articles.POST("/list",
		gin_ex.WrapClaimsAndReq[ReqList](a.List))
	articles.GET("/detail/:id",
		a.DetailSelf)

	reader := articles.Group("/reader")
	reader.GET("/:id", a.Detail)

	reader.POST("/like",
		gin_ex.WrapClaimsAndReq[LikeReq](a.Like))
	reader.POST("/collect",
		gin_ex.WrapClaimsAndReq[CollectReq](a.Collect))
}

func (a *ArticleHandler) Collect(ctx context.Context, req CollectReq, claims user.ClaimsUser) (utils.Response, error) {
	var err error

	err = a.like.Collect(ctx, a.ObjType, req.ObjID, claims.Id)

	if err != nil {
		return utils.Response{
			Code: domain.ErrSystem.ToInt(),
			Msg:  "系统错误",
			Data: nil,
		}, err
	}
	return utils.Response{
		Msg: "collect successful",
	}, nil
}

func (a *ArticleHandler) Like(ctx context.Context, req LikeReq, claims user.ClaimsUser) (utils.Response, error) {
	var err error
	if req.Liked {
		err = a.like.Like(ctx, a.ObjType, req.ID, claims.Id)
	} else {
		err = a.like.UnLike(ctx, a.ObjType, req.ID, claims.Id)
	}

	if err != nil {
		return utils.Response{
			Code: domain.ErrSystem.ToInt(),
			Msg:  "系统错误",
			Data: nil,
		}, err
	}
	return utils.Response{
		Msg: "like successful",
	}, nil
}

func (a *ArticleHandler) Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	artId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.FailWithMessage(domain.ErrSystem, "id必须为数字", ctx)
		//log
		return
	}
	claim, ok := ctx.Get("userClaims")
	if !ok {
		utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
		return
	}

	claims, ok := claim.(user.ClaimsUser)
	if !ok {
		utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
		return
	}

	art, err := a.svc.GetPublishedByID(ctx.Request.Context(), artId, claims.Id)
	if err != nil {
		utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
		return
	}

	// 增加阅读计数，也可以放middleware里
	//go func() {
	//	AsyErr := a.like.IncrReadCount(ctx.Request.Context(), a.ObjType, artId)
	//	if AsyErr != nil {
	//		// log
	//	}
	//}()

	utils.OK(ArticleVO{
		Content: art.Content,
		ID:      art.ID,
		Status:  art.Status.ToUint8(),
		Title:   art.Title,
		Ctime:   art.Ctime.Format(time.DateTime),
		Utime:   art.Utime.Format(time.DateTime),
	}, "success", ctx)
}

func (a *ArticleHandler) DetailSelf(ctx *gin.Context) {
	id := ctx.Param("id")
	artId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.FailWithMessage(domain.ErrSystem, "id必须为数字", ctx)
		//log
		return
	}

	claim, ok := ctx.Get("userClaims")
	if !ok {
		utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
		return
	}

	claims, ok := claim.(user.ClaimsUser)
	if !ok {
		utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
		return
	}

	art, err := a.svc.GetAuthorModelsByID(ctx.Request.Context(), artId)
	if err != nil {
		utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
		return
	}

	// 高危，即查询他人私有文章
	if claims.Id != art.ID {
		utils.FailWithMessage(domain.ErrSystem, "无权限", ctx)
		// 监控
		return
	}

	utils.OK(ArticleVO{
		Content: art.Content,
		ID:      art.ID,
		Status:  art.Status.ToUint8(),
		Title:   art.Title,
		Ctime:   art.Ctime.Format(time.DateTime),
		Utime:   art.Utime.Format(time.DateTime),
	}, "success", ctx)
}

func (a *ArticleHandler) List(ctx context.Context, req ReqList, claims user.ClaimsUser) (utils.Response, error) {
	res, err := a.svc.List(ctx, claims.Id, req.Limit, req.Offset)
	if err != nil {
		return utils.Response{
			Code: domain.ErrSystem.ToInt(),
			Msg:  err.Error(),
		}, nil
	}
	return utils.Response{
		Data: res,
	}, nil
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
