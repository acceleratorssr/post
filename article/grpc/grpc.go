package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	articlev1 "post/api/proto/gen/article/v1"
	"post/article/domain"
	"post/article/service"
)

type ArticleServiceServer struct {
	articlev1.UnimplementedArticleServiceServer
	svc service.ArticleService
}

// GetAuthorArticle 获取个人未发布文章内容
func (a *ArticleServiceServer) GetAuthorArticle(ctx context.Context, request *articlev1.GetAuthorArticleRequest) (*articlev1.GetAuthorArticleResponse, error) {
	article, err := a.svc.GetAuthorModelsByID(ctx, request.GetAid(), request.GetUid())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "文章不存在")
		}
		// log
		return nil, status.Errorf(codes.Unknown, "未知错误")
	}

	return &articlev1.GetAuthorArticleResponse{
		Data: toDTO([]domain.Article{*article}...)[0],
	}, nil
}

// ListSelf 获取个人前一百条文章内容
func (a *ArticleServiceServer) ListSelf(ctx context.Context, request *articlev1.ListSelfRequest) (*articlev1.ListSelfResponse, error) {
	list, err := a.svc.ListSelf(ctx, request.GetUid(), nil)
	if err != nil {
		// log
		return nil, status.Errorf(codes.Unknown, "未知错误")
	}
	return &articlev1.ListSelfResponse{
		Data: toDTO(list...),
	}, nil
}

func (a *ArticleServiceServer) Save(ctx context.Context, request *articlev1.SaveRequest) (*articlev1.SaveResponse, error) {
	err := a.svc.Save(ctx, ToDomainArticleAuthor(request.GetData()))
	if err != nil {
		// log
		return nil, status.Errorf(codes.Unknown, "未知错误")
	}

	return &articlev1.SaveResponse{}, nil
}

func (a *ArticleServiceServer) Publish(ctx context.Context, request *articlev1.PublishRequest) (*articlev1.PublishResponse, error) {
	err := a.svc.Publish(ctx, ToDomainArticleAuthor(request.GetData()))
	if err != nil {
		// log
		return nil, status.Errorf(codes.Unknown, "未知错误")
	}

	return &articlev1.PublishResponse{}, nil
}

// Withdraw 撤回发布库的文章
func (a *ArticleServiceServer) Withdraw(ctx context.Context, request *articlev1.WithdrawRequest) (*articlev1.WithdrawResponse, error) {
	err := a.svc.Withdraw(ctx, request.GetAid(), request.GetUid())
	if err != nil {
		// log
		return nil, status.Errorf(codes.Unknown, "未知错误")
	}

	return &articlev1.WithdrawResponse{}, nil
}

// GetPublishedByID 获取已发布文章内容
func (a *ArticleServiceServer) GetPublishedByID(ctx context.Context, request *articlev1.GetPublishedByIDRequest) (*articlev1.GetPublishedByIDResponse, error) {
	art, err := a.svc.GetPublishedByID(ctx, request.GetAid(), request.GetUid())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "文章不存在")
		}
		// log
		return nil, status.Errorf(codes.Unknown, "未知错误")
	}

	return &articlev1.GetPublishedByIDResponse{
		Data: toDTO(*art)[0],
	}, nil
}

// ListPublished 获取发布文章列表，无文章内容
func (a *ArticleServiceServer) ListPublished(ctx context.Context, request *articlev1.ListPublishedRequest) (*articlev1.ListPublishedResponse, error) {
	arts, err := a.svc.ListPublished(ctx, &domain.List{
		Limit:     int(request.GetLimit()),
		LastValue: request.GetLastValue(),
		OrderBy:   request.GetOrderBy(),
		Desc:      request.GetDesc(),
	})
	if err != nil {
		// log
		return nil, status.Errorf(codes.Unknown, "未知错误")
	}

	return &articlev1.ListPublishedResponse{
		Data: toDTO(arts...),
	}, nil
}

func (a *ArticleServiceServer) RegisterServer(server *grpc.Server) {
	articlev1.RegisterArticleServiceServer(server, a)
}

func toDTO(articles ...domain.Article) []*articlev1.Article {
	var protoArticles []*articlev1.Article

	for _, article := range articles {
		protoArticle := &articlev1.Article{
			ID: article.ID,
			Author: &articlev1.Author{
				ID:   article.Author.Id,
				Name: article.Author.Name,
			},
			Content: article.Content,
			Ctime:   article.Ctime,
			Title:   article.Title,
			Utime:   article.Utime,
		}
		protoArticles = append(protoArticles, protoArticle)
	}

	return protoArticles
}

func ToDomainArticleAuthor(article *articlev1.Article) *domain.Article {
	return &domain.Article{
		ID: article.GetID(),
		Author: domain.Author{
			Id:   article.GetAuthor().GetID(),
			Name: article.GetAuthor().GetName(),
		},
		Title:   article.GetTitle(),
		Content: article.GetContent(),
		Ctime:   article.GetCtime(),
		Utime:   article.GetUtime(),
	}
}

func NewArticleServiceServer(svc service.ArticleService) *ArticleServiceServer {
	return &ArticleServiceServer{
		svc: svc,
	}
}
