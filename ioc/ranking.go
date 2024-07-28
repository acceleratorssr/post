package ioc

import (
	"github.com/robfig/cron/v3"
	"post/job"
	distLock "post/redis_distributed_lock"
	"post/service"
)

func InitRankingJob(svc service.RankService, client *distLock.Client) *job.RankingJob {
	return job.NewRankingJob(svc, client)
}

func InitJobs(rankingJob *job.RankingJob) *cron.Cron {
	res := cron.New(cron.WithSeconds())
	_, err := res.AddJob("@every 1h", job.NewCronJobBuilder().Build(rankingJob))
	if err != nil {
		panic(err)
	}

	return res
}
