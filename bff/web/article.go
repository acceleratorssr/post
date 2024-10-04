package web

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
	articlev1 "post/api/proto/gen/article/v1"
	intrv1 "post/api/proto/gen/intr/v1"
	"post/pkg/gin_ex"
	"strconv"
)

type ArticleHandler struct {
	svc     articlev1.ArticleServiceClient
	like    intrv1.LikeServiceClient
	ObjType string
}

func NewArticleHandler(art articlev1.ArticleServiceClient, like intrv1.LikeServiceClient) *ArticleHandler {
	return &ArticleHandler{
		svc:     art,
		like:    like,
		ObjType: "article",
	}
}

func (a *ArticleHandler) Test(ctx *gin.Context) (*gin_ex.Response, error) {
	_, err := a.svc.Save(ctx, &articlev1.SaveRequest{
		Data: &articlev1.Article{
			Title:   "test",
			Content: "test",
		},
	})
	if err != nil {
		return nil, err
	}

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
	articles.Use(mw)
	articles.POST("/save", gin_ex.WrapWithReq[Req](a.Save))                   //保存文章
	articles.POST("/publish", gin_ex.WrapWithReq[Req](a.Publish))             // 发布文章
	articles.POST("/withdraw", gin_ex.WrapWithReq[ReqOnlyWithID](a.Withdraw)) // 撤回已发布文章
	articles.POST("/list_self", gin_ex.WrapWithReq[ReqList](a.ListSelf))      // 获取当前用户未发布文章列表
	articles.GET("/detail/:id", a.DetailSelf)                                 // 获取未发布文章内容

	reader := s.Group("/reader")
	reader.Use(mw)
	reader.GET("/:id", a.Detail) // 获取发布文章内容
	reader.POST("/list_publish", gin_ex.WrapWithReq[ReqList](a.ListPublished))

	reader.POST("/like", gin_ex.WrapWithReq[LikeReq](a.Like))          // 点赞
	reader.POST("/collect", gin_ex.WrapWithReq[CollectReq](a.Collect)) //收藏
}

func (a *ArticleHandler) Collect(ctx *gin.Context, req CollectReq) (*gin_ex.Response, error) {
	var err error

	_, err = a.like.Collect(ctx, &intrv1.CollectRequest{
		ObjID:   req.ObjID,
		ObjType: a.ObjType,
		Uid:     a.getUserID(ctx),
	})

	if err != nil {
		return &gin_ex.Response{
			Code: gin_ex.System,
			Msg:  "收藏失败",
		}, err
	}
	return &gin_ex.Response{
		Code: gin_ex.OK,
		Msg:  "收藏成功",
	}, nil
}

// Like todo 添加like等测试
func (a *ArticleHandler) Like(ctx *gin.Context, req LikeReq) (*gin_ex.Response, error) {
	var err error
	id := a.getUserID(ctx)
	if req.Liked {
		_, err = a.like.Like(ctx, &intrv1.LikeRequest{
			ObjID:   req.ObjID,
			ObjType: a.ObjType,
			Uid:     id,
		})
	} else {
		_, err = a.like.UnLike(ctx, &intrv1.UnLikeRequest{
			ObjID:   req.ObjID,
			ObjType: a.ObjType,
			Uid:     id,
		})
	}

	if err != nil {
		return &gin_ex.Response{
			Code: gin_ex.System,
			Msg:  "点赞相关操作失败",
		}, err
	}
	return &gin_ex.Response{
		Code: gin_ex.OK,
		Msg:  "点赞成功",
	}, nil
}

func (a *ArticleHandler) Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	artId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		gin_ex.FailWithMessage(ctx, gin_ex.InvalidArgument, "id必须为数字")
		//log
		return
	}

	art, err := a.svc.GetPublishedByID(ctx, &articlev1.GetPublishedByIDRequest{
		Aid: artId,
		Uid: a.getUserID(ctx),
	})
	if err != nil {
		gin_ex.FailWithMessage(ctx, gin_ex.System, err.Error())
		return
	}

	gin_ex.OKWithDataAndMsg(ctx, a.toVO(art.Data), "成功")
}

func (a *ArticleHandler) DetailSelf(ctx *gin.Context) {
	id := ctx.Param("id")
	artId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		gin_ex.FailWithMessage(ctx, gin_ex.InvalidArgument, "id必须为数字")
		//log
		return
	}

	art, err := a.svc.GetAuthorArticle(ctx.Request.Context(), &articlev1.GetAuthorArticleRequest{
		Uid: a.getUserID(ctx),
		Aid: artId,
	})
	if err != nil {
		// log err.Error()
		gin_ex.FailWithMessage(ctx, gin_ex.System, "查询失败")
		return
	}

	gin_ex.OKWithDataAndMsg(ctx, a.toVO(art.GetData()), "成功")
}

func (a *ArticleHandler) ListSelf(ctx *gin.Context, req ReqList) (*gin_ex.Response, error) {
	res, err := a.svc.ListSelf(ctx, &articlev1.ListSelfRequest{
		Uid:       a.getUserID(ctx),
		Limit:     req.Limit,
		LastValue: req.LastValue,
		OrderBy:   req.OrderBy,
		Desc:      req.Desc,
	})
	if err != nil {
		return &gin_ex.Response{
			Code: gin_ex.System,
			Msg:  err.Error(),
		}, nil
	}
	return &gin_ex.Response{
		Data: res,
	}, nil
}

func (a *ArticleHandler) ListPublished(ctx *gin.Context, req ReqList) (*gin_ex.Response, error) {
	res, err := a.svc.ListPublished(ctx, &articlev1.ListPublishedRequest{
		Limit:     req.Limit,
		LastValue: req.LastValue, // 上次查询的最小snow_id
		OrderBy:   req.OrderBy,   // 支持创建时间、更新时间作为排序依据
		Desc:      req.Desc,
	})
	if err != nil {
		// log err.Error()
		return &gin_ex.Response{
			Code: gin_ex.System,
			Msg:  "查询失败",
		}, nil
	}

	return &gin_ex.Response{
		Data: a.toVO(res.Data...),
	}, nil
}

func (a *ArticleHandler) Publish(ctx *gin.Context, req Req) (*gin_ex.Response, error) {
	id, err := a.svc.Publish(ctx, &articlev1.PublishRequest{
		Data: &articlev1.Article{
			ID: req.ID,
			Author: &articlev1.Author{
				ID:   a.getUserID(ctx),
				Name: a.getUsername(ctx),
			},
			Title:   req.Title,
			Content: req.Content,
		},
	})
	if err != nil {
		// log
		return &gin_ex.Response{
			Code: gin_ex.System,
			Msg:  "发布失败",
		}, err
	}
	return &gin_ex.Response{
		Data: id,
		Msg:  "发布成功",
	}, nil
}

func (a *ArticleHandler) Save(ctx *gin.Context, req Req) (*gin_ex.Response, error) {
	// 考虑压缩内容
	id, err := a.svc.Save(ctx, &articlev1.SaveRequest{
		Data: &articlev1.Article{
			ID: req.ID,
			Author: &articlev1.Author{
				ID:   a.getUserID(ctx),
				Name: a.getUsername(ctx),
			},
			Title:   req.Title,
			Content: req.Content,
		},
	})
	if err != nil {
		// log
		return &gin_ex.Response{
			Code: gin_ex.System,
			Msg:  "文章保存失败",
		}, err
	}

	gin_ex.OKWithDataAndMsg(ctx, id, "保存成功")
	return &gin_ex.Response{
		Msg: "保存成功",
	}, nil
}

func (a *ArticleHandler) Withdraw(ctx *gin.Context, req ReqOnlyWithID) (*gin_ex.Response, error) {
	_, err := a.svc.Withdraw(ctx, &articlev1.WithdrawRequest{
		Aid: req.ID,
		Uid: a.getUserID(ctx),
	})
	if err != nil {
		// log
		return &gin_ex.Response{
			Code: gin_ex.System,
			Msg:  "撤回文章失败",
		}, err
	}

	return &gin_ex.Response{
		Msg: "成功",
	}, nil
}

func (a *ArticleHandler) Delete(ctx *gin.Context) {
	// todo 调用草稿、发布、交互，删除相关数据
}

func (a *ArticleHandler) getUsername(ctx *gin.Context) string {
	username := ctx.MustGet("username")
	un, ok := username.(string)
	if !ok {
		gin_ex.FailWithMessage(ctx, gin_ex.System, "token存在问题")
		// log
		return ""
	}
	return un
}

func (a *ArticleHandler) getUserID(ctx *gin.Context) uint64 {
	uid := ctx.MustGet("uid")
	id, ok := uid.(uint64)
	if !ok {
		gin_ex.FailWithMessage(ctx, gin_ex.System, "token存在问题")
		// log
		return 0
	}
	return id
}

func (a *ArticleHandler) getNickname(ctx *gin.Context) string {
	nickname := ctx.MustGet("nickname")
	nn, ok := nickname.(string)
	if !ok {
		gin_ex.FailWithMessage(ctx, gin_ex.System, "token存在问题")
		// log
		return ""
	}
	return nn
}

func (a *ArticleHandler) toVO(art ...*articlev1.Article) []ArticleResp {
	artResp := make([]ArticleResp, len(art))
	for i, v := range art {
		artResp[i] = ArticleResp{
			Content: v.GetContent(),
			ID:      v.GetID(),
			Title:   v.GetTitle(),
			Ctime:   v.GetCtime(),
			Utime:   v.GetUtime(),
		}
	}
	return artResp
}
