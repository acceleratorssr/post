package job

import (
	"context"
	"post/interactive/repository/cache"
	"post/interactive/service"
	"strconv"
)

type BatchUpdateDBJob struct {
	svc   service.LikeService
	key   string
	cache cache.ArticleLikeCache
}

func NewBatchUpdateDBJobJob(svc service.LikeService, cache cache.ArticleLikeCache) *BatchUpdateDBJob {
	return &BatchUpdateDBJob{
		svc:   svc,
		cache: cache,
		key:   "cron_job:batch_update_db",
	}
}

func (u *BatchUpdateDBJob) Name() string {
	return "batch_update_db"
}

func (u *BatchUpdateDBJob) Run() error {
	ctx := context.Background()
	hmap, err := u.cache.GetCountByPrefix(ctx, "incr_read_count:article:")
	if err != nil {
		return err
	}

	merged := make(map[uint64]int64)
	objType := "article"
	var objID uint64 = 0
	for key := range hmap {
		v := hmap[key]
		objID, err = parseKey(key)
		if err != nil {
			// log
			continue
		}
		merged[objID] = v
	}

	err = u.svc.UpdateReadCountMany(ctx, objType, merged)
	if err != nil {
		// log
	}

	return nil
}

func parseKey(key string) (uint64, error) {
	s := key[len("incr_read_count:article:"):]
	parseUint, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return parseUint, nil
}
