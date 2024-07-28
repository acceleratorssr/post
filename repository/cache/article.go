package cache

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"post/domain"
	"strconv"
	"time"
)

type RedisArticleCache struct {
	client redis.Cmdable
}

type PageCache struct {
	ID     int64
	Title  string
	Status uint8
	Ctime  time.Time
	Utime  time.Time
}

//go:embed cnt.lua
var luaIncrCnt string

type ArticleCache interface {
	GetFirstPage(ctx context.Context, id int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, id int64, arts []domain.Article) error
	DeleteFirstPage(ctx context.Context, id int64) error

	GetArticleDetail(ctx context.Context, id int64) (domain.Article, error)
	SetArticleDetail(ctx context.Context, id int64, art domain.Article) error

	IncrReadCount(ctx context.Context, ObjType string, ObjID int64) error
	IncrLikeCount(ctx context.Context, ObjType string, ObjID int64) error
	DecrLikeCount(ctx context.Context, ObjType string, ObjID int64) error
}

func NewRedisArticleCache(client redis.Cmdable) ArticleCache {
	return &RedisArticleCache{
		client: client,
	}
}

func (r *RedisArticleCache) DecrLikeCount(ctx context.Context, ObjType string, ObjID int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.keyIncrReadCount(ObjType, ObjID)}, "like_cnt", -1).Err()
}

func (r *RedisArticleCache) IncrLikeCount(ctx context.Context, ObjType string, ObjID int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.keyIncrReadCount(ObjType, ObjID)}, "like_cnt", 1).Err()
}

func (r *RedisArticleCache) IncrReadCount(ctx context.Context, ObjType string, ObjID int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.keyIncrReadCount(ObjType, ObjID)}, "read_cnt", 1).Err()
}

func (r *RedisArticleCache) keyIncrReadCount(ObjType string, ObjID int64) string {
	return "article_incr_read_count:" + ObjType + ":" + strconv.FormatInt(ObjID, 10)
}

func (r *RedisArticleCache) GetArticleDetail(ctx context.Context, id int64) (domain.Article, error) {
	var art domain.Article
	cmd, err := r.client.Get(ctx, r.keyArticleDetail(id)).Result()
	if err != nil {
		if err == redis.Nil {
			//log
			return domain.Article{}, nil
		}
		//log
		return domain.Article{}, err
	}

	err = json.Unmarshal([]byte(cmd), &art)
	if err != nil {
		//log
		return domain.Article{}, err
	}
	return art, nil
}

func (r *RedisArticleCache) SetArticleDetail(ctx context.Context, id int64, art domain.Article) error {
	val, err := json.Marshal(r.ToPageCache(art))
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.keyArticleDetail(id), val, time.Hour*1).Err()
}

func (r *RedisArticleCache) keyArticleDetail(id int64) string {
	return "article_article_detail:" + strconv.FormatInt(id, 10)
}

func (r *RedisArticleCache) GetFirstPage(ctx context.Context, id int64) ([]domain.Article, error) {
	var arts []domain.Article
	cmd, err := r.client.Get(ctx, r.keyFirstPage(id)).Result()
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
func (r *RedisArticleCache) SetFirstPage(ctx context.Context, id int64, arts []domain.Article) error {
	marshal, err := json.Marshal(r.ToPageCacheMany(arts))
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.keyFirstPage(id), marshal, time.Hour*1).Err()
}

func (r *RedisArticleCache) keyFirstPage(id int64) string {
	return "article_first_page:" + strconv.FormatInt(id, 10)
}
func (r *RedisArticleCache) DeleteFirstPage(ctx context.Context, id int64) error {
	return r.client.Del(ctx, r.keyFirstPage(id)).Err()
}

func (r *RedisArticleCache) ToPageCacheMany(arts []domain.Article) []PageCache {
	page := make([]PageCache, 0)
	for _, art := range arts {
		page = append(page, r.ToPageCache(art))
	}
	return page
}

func (r *RedisArticleCache) ToPageCache(art domain.Article) PageCache {
	return PageCache{
		ID:     art.ID,
		Title:  art.Title,
		Status: art.Status.ToUint8(),
		Ctime:  art.Ctime,
		Utime:  art.Utime,
	}
}
