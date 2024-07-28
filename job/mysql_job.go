package job

import (
	"context"
	"errors"
	"golang.org/x/sync/semaphore"
	"post/domain"
	"post/service"
	"time"
)

type Executor interface {
	Name() string
	Exec(ctx context.Context, job domain.Job) error
}

type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, job domain.Job) error
}

func NewLocalFuncExecutor() *LocalFuncExecutor {
	return &LocalFuncExecutor{
		funcs: make(map[string]func(ctx context.Context, job domain.Job) error),
	}
}

func (l *LocalFuncExecutor) Name() string {
	return "local"
}

func (l *LocalFuncExecutor) Exec(ctx context.Context, job domain.Job) error {
	fn, ok := l.funcs[job.Name]
	if !ok {
		return errors.New("not found")
	}
	return fn(ctx, job)
}

func (l *LocalFuncExecutor) RegisterFuncs(name string, fn func(ctx context.Context, job domain.Job) error) {
	l.funcs[name] = fn
}

type Scheduler struct {
	execs   map[string]Executor
	svc     service.JobService
	limiter *semaphore.Weighted //信号量
}

// NewScheduler 还可以提供任务编排能力，
func NewScheduler(svc service.JobService) *Scheduler {
	return &Scheduler{
		execs:   make(map[string]Executor),
		svc:     svc,
		limiter: semaphore.NewWeighted(100),
	}
}

func (s *Scheduler) RegisterExecutor(exec Executor) {
	//s.execs[]
}

func (s *Scheduler) Schedule(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			//log
			return err
		}
		dbctx, cancel := context.WithTimeout(ctx, time.Second)

		job, err := s.svc.Preempt(dbctx)
		cancel()
		if err != nil {
			//log
		}

		exec, ok := s.execs[job.Executor]
		if !ok {
			//log\
			continue
		}

		go func() {
			defer func() {
				s.limiter.Release(1)
				// done
				e := job.CancelFunc()
				if e != nil {
					//log
				}
			}()

			e := exec.Exec(ctx, job)
			if e != nil {
				//log 任务失败
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			e = s.svc.SetNextExecTime(ctx, job)
			if e != nil {
				// log
			}
		}()

	}
}
