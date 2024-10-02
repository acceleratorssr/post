package job

import (
	"context"
	"post/article/service"
	"post/pkg/redis_ex/distributed_lock"
	"time"
)

type RankingJob struct {
	svc    service.RankService
	client *distributed_lock.Client
	key    string
	lock   *distributed_lock.Lock
}

func NewRankingJob(svc service.RankService, client *distributed_lock.Client) *RankingJob {
	return &RankingJob{
		svc:    svc,
		client: client,
		key:    "distLock:cron_job:rank",
	}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

// Run
// 可以考虑一个节点拿到锁后，一直负责计算排行榜，直到该节点关闭；
// todo 定时任务的基础上，做负载均衡
func (r *RankingJob) Run() error {
	if r.lock == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		lock, err := r.client.Lock(ctx, r.key, 6*time.Second,
			&distributed_lock.FixIntervalRetry{}, time.Second*2) // ctx对于整个函数的超时，timeout是对于redis的超时
		if err != nil {
			return err
		}
		r.lock = lock

		// 此处只能启动一个goroutine执行，故也可用mutex保证，
		// 但此处业务本身就互斥，所以也可以不用本地锁
		go func() {
			err := lock.AutoRefresh(time.Second*3, time.Second*3)
			if err != nil {
				// log
			}
			r.lock = nil
		}()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return r.svc.SetRankTopN(ctx, 100) // top100
}
