package repository

import (
	"context"
	"post/internal/domain"
	"post/internal/repository/dao"
	"time"
)

type JobRepository interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Done(ctx context.Context, id, version int64) error
	Refresh(id int64) error
	UpdateExecTime(ctx context.Context, id int64, t time.Time) error
	Stop(ctx context.Context, id int64) error
}

type PreemptJobRepository struct {
	dao dao.JobDAO
}

func NewPreemptJobRepository(dao dao.JobDAO) JobRepository {
	return &PreemptJobRepository{
		dao: dao,
	}
}

func (p *PreemptJobRepository) Stop(ctx context.Context, id int64) error {
	return p.dao.Stop(ctx, id)
}

func (p *PreemptJobRepository) UpdateExecTime(ctx context.Context, id int64, t time.Time) error {
	return p.dao.UpdateExecTime(ctx, id, t.Unix())
}

func (p *PreemptJobRepository) Refresh(id int64) error {
	return p.dao.UpdateUtime(id)
}

func (p *PreemptJobRepository) Done(ctx context.Context, id, version int64) error {
	return p.dao.Done(ctx, id, version)
}

func (p *PreemptJobRepository) Preempt(ctx context.Context) (domain.Job, error) {
	job, err := p.dao.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	return domain.Job{
		ID:       job.ID,
		Name:     job.Name,
		Executor: job.Executor,
		Cfg:      job.Cfg,
	}, nil
}
