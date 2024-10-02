package cache

import (
	"context"
	"errors"
	"post/article/domain"
	"sync/atomic"
	"time"
)

// LocalCacheForRank 避免数据竞争，此处应该也可以直接用具体类型，因为没有并发环境
// 注意本地缓存永远是优先读写，因为不涉及网络传输，几乎不会error
// 分布式环境下的本地缓存，实例间需要有扩散机制或者订阅机制
// 本地缓存也可以作为兜底的机制，如果发现数据库和缓存挂了，直接读本地缓存，可以不考虑过期时间
type LocalCacheForRank struct {
	topN      atomic.Value
	topNBrief atomic.Value
	ttl       atomic.Value
}

func NewLocalCacheForRank() *LocalCacheForRank {
	return &LocalCacheForRank{}
}

func (l *LocalCacheForRank) SetTopN(ctx context.Context, arts []domain.Article) error {
	l.topN.Store(arts)
	l.ttl.Store(time.Now().Add(time.Minute * 70))
	return nil
}
func (l *LocalCacheForRank) SetTopNBrief(ctx context.Context, arts []domain.Article) error {
	for i, _ := range arts {
		arts[i].ID = 0
		arts[i].Content = ""
		arts[i].Utime = 0
		arts[i].Status = 0
	}
	l.topNBrief.Store(arts)
	return nil
}

func (l *LocalCacheForRank) GetTopNBrief(ctx context.Context) ([]domain.Article, error) {
	arts := l.topN.Load().([]domain.Article)

	ttl := l.ttl.Load().(time.Time)

	// ttl早于time.Now()，即过时则错误
	if len(arts) == 0 || ttl.Before(time.Now()) {
		return nil, errors.New("localCache 挂了或者过期了")
	}
	return arts, nil
}
