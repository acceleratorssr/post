package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"post/article/domain"
	"strconv"
	"time"
)

type ArticleCache interface {
	SetListInfo(ctx context.Context, id uint64, arts []domain.Article) error
	GetListInfo(ctx context.Context, id uint64) ([]domain.Article, error)
	DeleteListInfo(ctx context.Context, id uint64) error

	SetListDetailByHashKey(ctx context.Context, arts []domain.Article, key string)
	GetListDetailByHashKey(ctx context.Context, id uint64, key string) *domain.Article

	GetArticleDetail(ctx context.Context, id uint64) (*domain.Article, error)
	SetArticleDetail(ctx context.Context, id uint64, art *domain.Article) error
}

type PageCache struct {
	ID     uint64
	Title  string
	Status uint8
	Ctime  int64
	Utime  int64
}

type RedisArticleCache struct {
	client redis.Cmdable
}

// todo 压缩+pipline优化性能

// SetListDetailByHashKey 以用户为hash表单位，过期时间为30min，保存文章；
// 可以是个性化推荐的列表，个人文章列表等
// 因为全错误不会造成数据不一致，所以有问题打日志就好
// key: 业务信息+uid
func (r *RedisArticleCache) SetListDetailByHashKey(ctx context.Context, arts []domain.Article, key string) {
	errs := make([]error, 0, len(arts))
	for _, art := range arts {
		val, err := json.Marshal(r.ToPageCache(art)[0])
		err = r.client.HSet(ctx, key, strconv.FormatUint(art.ID, 10), val).Err()
		errs = append(errs, err)
	}
	expireErr := r.client.Expire(ctx, key, 30*time.Minute).Err()
	if expireErr != nil {
		// log
	}

	go func() {
		for _, err := range errs {
			if err != nil {
				//log
			}
		}
	}()
}

func (r *RedisArticleCache) GetListDetailByHashKey(ctx context.Context, id uint64, key string) *domain.Article {
	result, err := r.client.HGet(ctx, key, strconv.FormatUint(id, 10)).Result()
	if err != nil {
		return nil
	}
	var art domain.Article
	err = json.Unmarshal([]byte(result), &art)
	if err != nil {
		// log
		return nil
	}

	return &art
}

func (r *RedisArticleCache) GetArticleDetail(ctx context.Context, id uint64) (*domain.Article, error) {
	var art domain.Article
	cmd, err := r.client.Get(ctx, r.keyArticleDetail(id)).Result()
	if err != nil {
		//log
		return nil, err
	}

	err = json.Unmarshal([]byte(cmd), &art)
	if err != nil {
		//log
		return nil, err
	}
	return &art, nil
}

func (r *RedisArticleCache) SetArticleDetail(ctx context.Context, id uint64, art *domain.Article) error {
	val, err := json.Marshal(r.ToPageCache(*art)[0])
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.keyArticleDetail(id), val, time.Hour*1).Err()
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

func NewRedisArticleCache(client redis.Cmdable) ArticleCache {
	return &RedisArticleCache{
		client: client,
	}
}
