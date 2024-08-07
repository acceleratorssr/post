package job

import (
	"context"
	"errors"
	"golang.org/x/sync/semaphore"
	"post/internal/domain"
	"post/internal/service"
	"time"
)

type Executor interface {
	Name() string
	Exec(ctx context.Context, job domain.Job) error
}

// todo Executor还可以调用http/grpc接口

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
		limiter: semaphore.NewWeighted(100), //信号量控制goroutine
	}
}

func (s *Scheduler) RegisterExecutor(exec Executor) {
	s.execs[exec.Name()] = exec
}

// Schedule 调度
func (s *Scheduler) Schedule(ctx context.Context) error {
	for { // 循环从mysql取出任务，如果自己用对应任务的Executor，则可正常运行
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			//log
			return err
		}
		dbctx, cancel := context.WithTimeout(ctx, time.Second)

		job, err := s.svc.Preempt(dbctx) // 获取任务的执行权
		cancel()
		if err != nil {
			//log
		}

		exec, ok := s.execs[job.Executor] // 查看当前执行器有没有对应的local任务执行
		if !ok {
			//log
			s.limiter.Release(1)
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

			e := exec.Exec(ctx, job) // 执行该任务的job
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
