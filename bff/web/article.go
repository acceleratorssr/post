package web

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
	intrv1 "post/api/proto/gen/intr/v1"
	"post/article/domain"
	"post/article/service"
	"post/article/user"
	"post/pkg/gin_ex"
	"strconv"
)

type ArticleHandler struct {
	svc     service.ArticleService
	like    intrv1.LikeServiceClient
	ObjType string
}

func NewArticleHandler(svc service.ArticleService, like intrv1.LikeServiceClient) *ArticleHandler {
	return &ArticleHandler{
		svc:     svc,
		like:    like,
		ObjType: "article",
	}
}

func (a *ArticleHandler) Test(ctx *gin.Context) (*gin_ex.Response, error) {
	// 注意此处传入的是context.context，而不是gin.Context
	a.svc.Save(ctx, &domain.Article{
		Title:   "test",
		Content: "test",
	})
	// 复用
	span := trace.SpanFromContext(ctx)
	span.AddEvent("---50%---")

	if rand.Int31n(100)%2 == 0 {

		return &gin_ex.Response{
			Code: gin_ex.InvalidArgument,
			Msg:  "fail",
		}, nil
	}
	return &gin_ex.Response{
		Code: gin_ex.OK,
		Msg:  "ok",
	}, nil
}

func (a *ArticleHandler) RegisterRoutes(s *gin.Engine, mw gin.HandlerFunc) {
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

func (a *ArticleHandler) Collect(ctx *gin.Context, req CollectReq, claims user.ClaimsUser) (gin_ex.Response, error) {
	var err error

	_, err = a.like.Collect(ctx, &intrv1.CollectRequest{
		ObjID:   req.ObjID,
		ObjType: a.ObjType,
		Uid:     claims.Id,
	})

	if err != nil {
		return gin_ex.Response{
			Code: gin_ex.System,
			Msg:  "系统错误",
			Data: nil,
		}, err
	}
	return gin_ex.Response{
		Msg: "collect successful",
	}, nil
}

// Like todo 添加like等测试
func (a *ArticleHandler) Like(ctx *gin.Context, req LikeReq, claims user.ClaimsUser) (gin_ex.Response, error) {
	var err error
	if req.Liked {
		_, err = a.like.Like(ctx, &intrv1.LikeRequest{
			ObjID:   req.ObjID,
			ObjType: a.ObjType,
			Uid:     claims.Id,
		})
	} else {
		_, err = a.like.UnLike(ctx, &intrv1.UnLikeRequest{
			ObjID:   req.ObjID,
			ObjType: a.ObjType,
			Uid:     claims.Id,
		})
	}

	if err != nil {
		return gin_ex.Response{
			Code: gin_ex.System,
			Msg:  "系统错误",
			Data: nil,
		}, err
	}
	return gin_ex.Response{
		Msg: "like successful",
	}, nil
}

func (a *ArticleHandler) Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	artId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		gin_ex.FailWithMessage(ctx, gin_ex.InvalidArgument, "id必须为数字")
		//log
		return
	}

	claim, ok := ctx.Get("userClaims")
	if !ok {
		gin_ex.FailWithMessage(ctx, gin_ex.System, err.Error())
		return
	}

	claims, ok := claim.(*user.ClaimsUser)
	if !ok {
		gin_ex.FailWithMessage(ctx, gin_ex.System, err.Error())
		return
	}

	art, err := a.svc.GetPublishedByID(ctx.Request.Context(), artId, claims.Id)
	if err != nil {
		gin_ex.FailWithMessage(ctx, gin_ex.System, err.Error())
		return
	}

	// 增加阅读计数，也可以放middleware里
	//go func() {
	//	AsyErr := a.like.IncrReadCount(ctx.Request.Context(), a.ObjType, artId)
	//	if AsyErr != nil {
	//		// log
	//	}
	//}()

	gin_ex.OKWithDataAndMsg(ctx, ArticleResp{
		Content: art.Content,
		ID:      art.ID,
		Status:  art.Status.ToUint8(),
		Title:   art.Title,
		Ctime:   art.Ctime,
		Utime:   art.Utime,
	}, "success")
}

func (a *ArticleHandler) DetailSelf(ctx *gin.Context) {
	id := ctx.Param("id")
	artId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		gin_ex.FailWithMessage(ctx, gin_ex.System, "id必须为数字")
		//log
		return
	}

	claim, ok := ctx.Get("userClaims")
	if !ok {
		gin_ex.FailWithMessage(ctx, gin_ex.System, err.Error())
		return
	}

	claims, ok := claim.(user.ClaimsUser)
	if !ok {
		gin_ex.FailWithMessage(ctx, gin_ex.System, err.Error())
		return
	}

	art, err := a.svc.GetAuthorModelsByID(ctx.Request.Context(), artId)
	if err != nil {
		gin_ex.FailWithMessage(ctx, gin_ex.System, err.Error())
		return
	}

	// 高危，即查询他人私有文章
	if claims.Id != art.ID {
		gin_ex.FailWithMessage(ctx, gin_ex.PermissionDenied, "无权限")
		// 监控
		return
	}

	gin_ex.OKWithDataAndMsg(ctx, ArticleResp{
		Content: art.Content,
		ID:      art.ID,
		Status:  art.Status.ToUint8(),
		Title:   art.Title,
		Ctime:   art.Ctime,
		Utime:   art.Utime,
	}, "success")
}

func (a *ArticleHandler) List(ctx *gin.Context, req ReqList, claims user.ClaimsUser) (gin_ex.Response, error) {
	res, err := a.svc.List(ctx, claims.Id, req.Limit, req.Offset)
	if err != nil {
		return gin_ex.Response{
			Code: gin_ex.System,
			Msg:  err.Error(),
		}, nil
	}
	return gin_ex.Response{
		Data: res,
	}, nil
}

func (a *ArticleHandler) Publish(ctx *gin.Context) {
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	}
	// TODO 检查输入

	claims := a.getUserInfo(ctx)

	id, err := a.svc.Publish(ctx, req.toDomain(claims.Id, claims.Name))
	if err != nil {
		gin_ex.Fail(ctx, gin_ex.System, err.Error(), "系统错误")
		// log
		return
	}
	gin_ex.OKWithDataAndMsg(ctx, id, "success")
}

func (a *ArticleHandler) Save(ctx *gin.Context) {
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	} //需要使用ShouldBindJSON，如果使用bind则获取不到值
	// TODO 检查输入

	claims := a.getUserInfo(ctx)

	id, err := a.svc.Save(ctx, req.toDomain(claims.Id, claims.Name))
	if err != nil {
		gin_ex.FailWithMessage(ctx, gin_ex.System, "系统错误")
		// log
		return
	}

	gin_ex.OKWithDataAndMsg(ctx, id, "success")
}

func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	var req ReqOnlyWithID
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	} //需要使用ShouldBindJSON，如果使用bind则获取不到值
	// TODO 检查输入

	claims := a.getUserInfo(ctx)
	err := a.svc.Withdraw(ctx, req.toDomain(claims.Id, claims.Name))
	if err != nil {
		gin_ex.FailWithMessage(ctx, gin_ex.System, "系统错误")
		// log
		return
	}

	gin_ex.OKWithMessage(ctx, "success")
}

func (a *ArticleHandler) getUserInfo(ctx *gin.Context) *user.ClaimsUser {
	claim := ctx.MustGet("userClaims")
	claims, ok := claim.(*user.ClaimsUser)
	if !ok {
		gin_ex.FailWithMessage(ctx, gin_ex.System, "系统错误")
		// log
		return nil
	}
	return claims
}
