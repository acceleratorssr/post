package job

import (
	"context"
	"post/service"
	"time"
)

type RankingJob struct {
	svc service.RankService
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return r.svc.GetRankTopN(ctx, 100) // top100
}
