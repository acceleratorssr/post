package cache

import (
	"context"
	"errors"
	"post/article/domain"
	"sync"
	"sync/atomic"
	"time"
)

// LocalCacheForRank 避免数据竞争，此处应该也可以直接用具体类型，因为没有并发环境
// 注意本地缓存永远是优先读写，因为不涉及网络传输，几乎不会error
// 分布式环境下的本地缓存，实例间需要有扩散机制或者订阅机制
// 本地缓存也可以作为兜底的机制，如果发现数据库和缓存挂了，直接读本地缓存，可以不考虑过期时间
// topN 和 topNBrief 是一个同期缓存，所以需要一个过期时间，否则会无限增长
type LocalCacheForRank struct {
	topN      sync.Map
	keys      sync.Map
	topNBrief atomic.Value
	ttl       atomic.Value
}

func NewLocalCacheForRank() *LocalCacheForRank {
	return &LocalCacheForRank{}
}

// SetTopN 第二次触发时才会淘汰旧缓存
func (l *LocalCacheForRank) SetTopN(ctx context.Context, arts []domain.Article) error {
	newKeys := sync.Map{}

	for _, art := range arts {
		key := domain.GetArtCacheKey(art.ID)
		l.topN.Store(key, art)
		newKeys.Store(key, struct{}{})
	}

	// 过期旧的缓存数据
	l.keys.Range(func(k, _ interface{}) bool {
		key := k.(string)
		if _, ok := newKeys.Load(key); !ok {
			l.topN.Delete(key)
			l.keys.Delete(key)
		}
		return true
	})

	newKeys.Range(func(k, _ interface{}) bool {
		l.keys.Store(k, struct{}{})
		return true
	})

	l.ttl.Store(time.Now().Add(time.Minute * 70))

	return nil
}

func (l *LocalCacheForRank) SetTopNBrief(ctx context.Context, arts []domain.Article) error {
	for i, _ := range arts {
		arts[i].ID = 0
		arts[i].Content = ""
		arts[i].Utime = 0
	}
	l.topNBrief.Store(arts)
	return nil
}

func (l *LocalCacheForRank) GetTopNBrief(ctx context.Context) ([]domain.Article, error) {
	arts := l.topNBrief.Load().([]domain.Article)

	ttl := l.ttl.Load().(time.Time)

	// ttl早于time.Now()，即过时则错误
	if len(arts) == 0 || ttl.Before(time.Now()) {
		return nil, errors.New("localCache 挂了或者过期了")
	}
	return arts, nil
}
