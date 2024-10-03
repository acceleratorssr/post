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

func (a *ArticleServiceServer) GetAuthorArticle(ctx context.Context, request *articlev1.GetAuthorArticleRequest) (*articlev1.GetAuthorArticleResponse, error) {
	article, err := a.svc.GetAuthorModelsByID(ctx, request.GetAid())
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

func (a *ArticleServiceServer) ListSelf(ctx context.Context, request *articlev1.ListSelfRequest) (*articlev1.ListSelfResponse, error) {
	list, err := a.svc.List(ctx, request.GetUid(), int(request.GetLimit()), int(request.GetOffset()))
	if err != nil {
		// log
		return nil, status.Errorf(codes.Unknown, "未知错误")
	}
	return &articlev1.ListSelfResponse{
		Data: toDTO(list...),
	}, nil
}

func (a *ArticleServiceServer) Save(ctx context.Context, request *articlev1.SaveRequest) (*articlev1.SaveResponse, error) {
	_, err := a.svc.Save(ctx, ToDomainArticle(request.GetData()))
	if err != nil {
		// log
		return nil, status.Errorf(codes.Unknown, "未知错误")
	}

	return &articlev1.SaveResponse{}, nil
}

func (a *ArticleServiceServer) Publish(ctx context.Context, request *articlev1.PublishRequest) (*articlev1.PublishResponse, error) {
	_, err := a.svc.Publish(ctx, ToDomainArticle(request.GetData()))
	if err != nil {
		// log
		return nil, status.Errorf(codes.Unknown, "未知错误")
	}

	return &articlev1.PublishResponse{}, nil
}

func (a *ArticleServiceServer) Withdraw(ctx context.Context, request *articlev1.WithdrawRequest) (*articlev1.WithdrawResponse, error) {
	err := a.svc.Withdraw(ctx, ToDomainArticle(request.GetData()))
	if err != nil {
		// log
		return nil, status.Errorf(codes.Unknown, "未知错误")
	}

	return &articlev1.WithdrawResponse{}, nil
}

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

func (a *ArticleServiceServer) GetPublishedByIDS(ctx context.Context, request *articlev1.GetPublishedByIDSRequest) (*articlev1.GetPublishedByIDSResponse, error) {
	arts, err := a.svc.GetPublishedByIDS(ctx, request.GetAids())
	if err != nil {
		// log
		return nil, status.Errorf(codes.Unknown, "未知错误")
	}

	return &articlev1.GetPublishedByIDSResponse{
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
			Author: &articlev1.Author{
				ID:   article.Author.Id,
				Name: article.Author.Name,
			},
			Content: article.Content,
			Ctime:   article.Ctime,
			Status:  int32(article.Status),
			Title:   article.Title,
			Utime:   article.Utime,
		}
		protoArticles = append(protoArticles, protoArticle)
	}

	return protoArticles
}

func ToDomainArticle(article *articlev1.Article) *domain.Article {
	return &domain.Article{
		Author: domain.Author{
			Id:   article.GetAuthor().GetID(),
			Name: article.GetAuthor().GetName(),
		},
		Title:   article.GetTitle(),
		Content: article.GetContent(),
		Status:  domain.StatusType(article.Status),
		Ctime:   article.GetCtime(),
		Utime:   article.GetUtime(),
	}
}

func NewArticleServiceServer(svc service.ArticleService) *ArticleServiceServer {
	return &ArticleServiceServer{
		svc: svc,
	}
}
