package ioc

import (
	"github.com/robfig/cron/v3"
	"post/interactive/job"
	"post/interactive/repository/cache"
	"post/interactive/service"
)

func InitBatchUpdateDBJob(svc service.LikeService, cache cache.ArticleLikeCache) *job.BatchUpdateDBJob {
	return job.NewBatchUpdateDBJobJob(svc, cache)
}

func InitJobs(rankingJob *job.BatchUpdateDBJob) *cron.Cron {
	res := cron.New(cron.WithSeconds())
	_, err := res.AddJob("@every 1m", job.NewCronJobBuilder().Build(rankingJob))
	if err != nil {
		panic(err)
	}

	return res
}
