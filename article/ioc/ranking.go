package ioc

import (
	"github.com/robfig/cron/v3"
	"post/article/job"
	"post/article/service"
	distLock "post/pkg/redis-extra/distributed_lock"
)

func InitRankingJob(svc service.RankService, client *distLock.Client) *job.RankingJob {
	return job.NewRankingJob(svc, client)
}

// InitJobs 定时任务
func InitJobs(rankingJob *job.RankingJob) *cron.Cron {
	res := cron.New(cron.WithSeconds())
	_, err := res.AddJob("@every 1h", job.NewCronJobBuilder().Build(rankingJob))
	if err != nil {
		panic(err)
	}

	return res
}
