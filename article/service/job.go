package service

import (
	"context"
	"post/article/domain"
	"post/article/repository"
	"time"
)

type JobService interface {
	// Preempt 抢占式的任务调度
	Preempt(ctx context.Context) (domain.Job, error)
	refresh(id int64) error

	// SetNextExecTime 参考https://www.yuque.com/u32495242/qqahru/mlfcpyzyrem61g44#HeQtr
	SetNextExecTime(ctx context.Context, job domain.Job) error
}

type CronJobService struct {
	repo            repository.JobRepository
	refreshInterval time.Duration
}

func NewCronJobService(repo repository.JobRepository) JobService {
	return &CronJobService{
		repo: repo,
	}
}

func (p *CronJobService) SetNextExecTime(ctx context.Context, job domain.Job) error {
	next := job.NextExecTime()
	if next.IsZero() {
		// 停止任务
		return p.repo.Stop(ctx, job.ID)
	}
	return p.repo.UpdateExecTime(ctx, job.ID, next)
}

func (p *CronJobService) refresh(id int64) error {
	return p.repo.Refresh(id)
}

func (p *CronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	job, err := p.repo.Preempt(ctx)

	// 续约，通过utime和status判断节点是否失联
	ticker := time.NewTicker(p.refreshInterval)

	go func() {
		for range ticker.C {
			err := p.refresh(job.ID) // 通过任务的utime判断节点是否失联
			if err != nil {
				return
			}
		}
	}()

	job.CancelFunc = func() error {
		ticker.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		return p.repo.Done(ctx, job.ID, job.Version)
	}
	return job, err
}
