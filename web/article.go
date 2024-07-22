package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"post/domain"
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

func (a *ArticleHandler) Test(ctx *gin.Context) {
	a.svc.Save(ctx, domain.Article{
		Title:   "test",
		Content: "test",
	})
	ctx.JSON(200, gin.H{
		"message": "ok",
	})
}

func (a *ArticleHandler) RegisterRoutes(s *gin.Engine) {
	s.POST("/test", a.Test)
	articles := s.Group("/articles")
	articles.POST("/save", a.Save)
	articles.POST("/publish", a.Publish)
	articles.POST("/withdraw", a.Withdraw)
	articles.POST("/list",
		WrapClaimsAndReq[ReqList](a.List))
	articles.GET("/detail/:id",
		a.DetailSelf)

	reader := articles.Group("/reader")
	reader.GET("/:id", a.Detail)

	reader.POST("/like",
		WrapClaimsAndReq[LikeReq](a.Like))
	reader.POST("/collect",
		WrapClaimsAndReq[CollectReq](a.Collect))
}

func (a *ArticleHandler) Collect(ctx *gin.Context, req CollectReq, claims user.ClaimsUser) (utils.Response, error) {
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

func (a *ArticleHandler) Like(ctx *gin.Context, req LikeReq, claims user.ClaimsUser) (utils.Response, error) {
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

	art, err := a.svc.GetPublishedByID(ctx, artId, claims.Id)
	if err != nil {
		utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
		return
	}

	// 增加阅读计数，也可以放middleware里
	//go func() {
	//	AsyErr := a.like.IncrReadCount(ctx, a.ObjType, artId)
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

	art, err := a.svc.GetAuthorModelsByID(ctx, artId)
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

func (a *ArticleHandler) List(ctx *gin.Context, req ReqList, claims user.ClaimsUser) (utils.Response, error) {
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

// WrapClaimsAndReq
// TODO 除此之外还可以考虑单独解析claims或者req，解决全部post
func WrapClaimsAndReq[Req any](fn func(*gin.Context, Req, user.ClaimsUser) (utils.Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			err = fmt.Errorf("解析请求失败%w", err)
			utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
			return
		}

		claim, ok := ctx.Get("userClaims")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			err := fmt.Errorf("无法获得 claims:%v", ctx.Request.URL.Path)
			utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
			return
		}

		claims, ok := claim.(user.ClaimsUser)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			err := fmt.Errorf("无法获得 claims:%v", ctx.Request.URL.Path)
			utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
			return
		}

		res, err := fn(ctx, req, claims)

		if err != nil {
			err = fmt.Errorf("业务失败:%w", err)
			utils.FailWithMessage(domain.ErrSystem, err.Error(), ctx)
		}

		utils.OK(res.Data, res.Msg, ctx)
	}
}
