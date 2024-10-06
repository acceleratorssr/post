package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"post/article/domain"
	"post/article/repository/cache/compression"
	"strconv"
	"time"
)

type ArticleCache interface {
	SetListInfo(ctx context.Context, id uint64, arts []domain.Article) error
	GetListInfo(ctx context.Context, id uint64) ([]domain.Article, error)
	DeleteListInfo(ctx context.Context, id uint64) error

	SetListDetailByHashKey(ctx context.Context, arts []domain.Article, key string)
	GetListDetailByHashKey(ctx context.Context, id uint64, key string) (*domain.Article, error)

	GetArticleDetail(ctx context.Context, id uint64) (*domain.Article, error)
	SetArticleDetail(ctx context.Context, art *domain.Article) error
}

type PageCache struct {
	ID     uint64
	Title  string
	Status uint8
	Ctime  int64
	Utime  int64
}

type RedisArticleCache struct {
	client      redis.Cmdable
	compression compression.Compression
}

// todo 考虑延迟读写超时时间，避免pipline因为超时挂掉
// todo 缓存没设计好，hash应该用在频繁变动字段的对象上，文章其实可以直接分id全存string

// SetListDetailByHashKey 以 用户的uid 为hash表分表依据，过期时间为30min，保存文章；（可改为一类用户）
// 可以是个性化推荐的列表，个人文章列表等
// 因为全错误不会造成数据不一致，所以有问题打日志就好
func (r *RedisArticleCache) SetListDetailByHashKey(ctx context.Context, arts []domain.Article, key string) {
	// 创建 pipeline
	pipe := r.client.Pipeline()

	for _, art := range arts {
		compress, err := r.compression.Compressed(&art)
		if err != nil {
			// log，机器出错
			return
		}
		// 使用 pipeline 的 HSet
		pipe.HSet(ctx, key, strconv.FormatUint(art.ID, 10), compress)
	}

	// 设置过期时间
	pipe.Expire(ctx, key, 30*time.Minute)

	// 执行 pipeline
	cmders, err := pipe.Exec(ctx)
	if err != nil {
		// log，pipeline 出错
	}
	go func() {
		// 记录每个命令的错误
		for _, cmdErr := range cmders {
			if cmdErr.Err() != nil {
				// log 错误
			}
		}
	}()
}

func (r *RedisArticleCache) GetListDetailByHashKey(ctx context.Context, aid uint64, key string) (*domain.Article, error) {
	result, err := r.client.HGet(ctx, key, strconv.FormatUint(aid, 10)).Result()
	if err != nil {
		return nil, err
	}

	var art *domain.Article
	err = r.compression.Decompress([]byte(result), art)
	if err != nil {
		// log
		return nil, err
	}

	return art, nil
}

func (r *RedisArticleCache) GetArticleDetail(ctx context.Context, id uint64) (*domain.Article, error) {
	cmd, err := r.client.Get(ctx, r.keyArticleDetail(id)).Result()
	if err != nil {
		//log
		return nil, err
	}

	var art *domain.Article
	err = r.compression.Decompress([]byte(cmd), art)
	if err != nil {
		//log，机器错误
		return nil, err
	}

	return art, nil
}

func (r *RedisArticleCache) SetArticleDetail(ctx context.Context, art *domain.Article) error {
	compressed, err := r.compression.Compressed(art)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.keyArticleDetail(art.ID), compressed, time.Hour*1).Err()
}

func (r *RedisArticleCache) keyArticleDetail(id uint64) string {
	return "article_article_detail:" + strconv.FormatUint(id, 10)
}

// SetListInfo 保存文章信息
func (r *RedisArticleCache) SetListInfo(ctx context.Context, id uint64, arts []domain.Article) error {
	marshal, err := json.Marshal(r.ToPageCache(arts...))
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.keyListInfo(id), marshal, time.Hour*1).Err()
}

func (r *RedisArticleCache) GetListInfo(ctx context.Context, id uint64) ([]domain.Article, error) {
	var arts []domain.Article
	cmd, err := r.client.Get(ctx, r.keyListInfo(id)).Result()
	if err != nil {
		if err == redis.Nil {
			//log
			return nil, nil
		}
		//log
		return nil, err
	}

	err = json.Unmarshal([]byte(cmd), &arts)
	if err != nil {
		//log
		return nil, err
	}
	return arts, nil
}

func (r *RedisArticleCache) DeleteListInfo(ctx context.Context, id uint64) error {
	return r.client.Del(ctx, r.keyListInfo(id)).Err()
}

func (r *RedisArticleCache) keyListInfo(id uint64) string {
	return "article_list_info:" + strconv.FormatUint(id, 10)
}

func (r *RedisArticleCache) ToPageCache(arts ...domain.Article) []PageCache {
	page := make([]PageCache, 0)
	for _, art := range arts {
		page = append(page, PageCache{
			ID:    art.ID,
			Title: art.Title,
			Ctime: art.Ctime,
			Utime: art.Utime,
		})
	}
	return page
}

func NewRedisArticleCache(client redis.Cmdable, compression compression.Compression) ArticleCache {
	return &RedisArticleCache{
		client:      client,
		compression: compression,
	}
}
