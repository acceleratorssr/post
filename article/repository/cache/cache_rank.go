package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"post/article/domain"
	"strconv"
	"time"
)

type RankCache interface {
	SetTopN(ctx context.Context, arts []domain.Article) error
	SetTopNBrief(ctx context.Context, arts []domain.Article) error
	GetTopNBrief(ctx context.Context) ([]domain.Article, error)
}

type RankRedisCache struct {
	client redis.Cmdable
}

func NewRankCache(client redis.Cmdable) RankCache {
	return &RankRedisCache{
		client: client,
	}
}
func (r *RankRedisCache) GetTopNBrief(ctx context.Context) ([]domain.Article, error) {
	val, err := r.client.Get(ctx, r.keyTopNBrief("article", 0)).Result()
	if err != nil {
		panic(err)
	}

	var arts []domain.Article
	err = json.Unmarshal([]byte(val), &arts)
	if err != nil {
		return nil, err
	}

	return arts, nil
}

func (r *RankRedisCache) SetTopNBrief(ctx context.Context, arts []domain.Article) error {
	for i, _ := range arts {
		arts[i].ID = 0
		arts[i].Content = ""
		arts[i].Utime = 0
		arts[i].Status = 0
	}
	val, err := json.Marshal(arts)
	err = r.client.Set(ctx, r.keyTopNBrief("article", 0), val, 70*time.Minute).Err()
	if err != nil {
		// log
	}
	return nil
}

// SetTopN
// todo 从查询角度，可考虑在缓存中设置热榜标识（不存mysql），查询发现热榜则先local，再redis，最后数据库
// 每个文章单独设缓存
func (r *RankRedisCache) SetTopN(ctx context.Context, arts []domain.Article) error {
	for i, _ := range arts {
		val, err := json.Marshal(arts[i])
		err = r.client.Set(ctx, r.keyTopN("article", arts[i].ID), val, 70*time.Minute).Err()
		if err != nil {
			// log
		}
	}
	return nil
}
func (r *RankRedisCache) keyTopN(ObjType string, ObjID uint64) string {
	return ObjType + "_rank:" + strconv.FormatUint(ObjID, 10)
}

func (r *RankRedisCache) keyTopNBrief(ObjType string, ObjID uint64) string {
	return ObjType + "_rank_brief:" + strconv.FormatUint(ObjID, 10)
}
