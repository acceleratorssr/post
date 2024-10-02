package ioc

import (
	"context"
	"post/article/domain"
	"post/article/job"
	"post/article/service"
	"time"
)

func InitScheduler(svc service.JobService, local *job.LocalFuncExecutor) *job.Scheduler {
	res := job.NewScheduler(svc)
	res.RegisterExecutor(local)
	return res
}

func InitLocalFuncExecutor(svc service.RankService) *job.LocalFuncExecutor {
	res := job.NewLocalFuncExecutor()
	// 在mysql插入任务记录
	res.RegisterFuncs("rank", func(ctx context.Context, job domain.Job) error {
		ctx, cancel := context.WithTimeout(ctx, 70*time.Minute)
		defer cancel()
		return svc.SetRankTopN(ctx, 100)
	})
	return res
}
