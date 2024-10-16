package service

import (
	"context"
	"github.com/zhenghaoz/gorse/client"
)

type RecommendService interface {
	GetRecommend(ctx context.Context, userId, writeBackType, writeBackDelay string, n, offset int, category ...string) ([]string, error)
	GetNeighbors(ctx context.Context, itemId, userId string, n, offset int) ([]client.Score, error)
	GetUser(ctx context.Context, userId string) (client.User, error)
	GetUsers(ctx context.Context, cursor string, n int) (client.Users, error)
	GetItemByID(ctx context.Context, itemId string) (client.Item, error)
}

type recommendService struct {
	gorse *client.GorseClient
}

func (r *recommendService) GetRecommend(ctx context.Context, userId, writeBackType, writeBackDelay string, n, offset int, category ...string) ([]string, error) {
	if len(category) == 1 {
		// 单类别推荐
		rc, err := r.gorse.GetItemRecommendWithCategory(ctx, userId, category[0], writeBackType, writeBackDelay, n, offset)
		if err != nil {
			return nil, err
		}
		return rc, nil
	}
	rc, err := r.gorse.GetItemRecommend(ctx, userId, category, writeBackType, writeBackDelay, n, offset) // 多类别推荐
	if err != nil {
		return nil, err
	}
	return rc, nil
}

func (r *recommendService) GetNeighbors(ctx context.Context, itemId, userId string, n, offset int) ([]client.Score, error) {
	scores, err := r.gorse.GetItemNeighbors(ctx, itemId, userId, n, offset)
	if err != nil {
		return nil, err
	}
	return scores, nil
}

func (r *recommendService) GetUser(ctx context.Context, userId string) (client.User, error) {
	return r.gorse.GetUser(ctx, userId)
}

func (r *recommendService) GetUsers(ctx context.Context, cursor string, n int) (client.Users, error) {
	users, err := r.gorse.GetUsers(ctx, cursor, n)
	if err != nil {
		return client.Users{}, err
	}
	return users, nil
}

func (r *recommendService) GetItemByID(ctx context.Context, itemId string) (client.Item, error) {
	return r.gorse.GetItem(ctx, itemId)
}

func NewRecommendService(gorse *client.GorseClient) RecommendService {
	return &recommendService{
		gorse: gorse,
	}
}
