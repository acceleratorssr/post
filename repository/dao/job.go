package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

const (
	jobStatusWaiting = iota
	jobStatusRunning
	jobStatusFinished
)

type Job struct {
	ID       int64  `gorm:"primaryKey,autoIncrement"`
	Name     string `gorm:"unique"`
	Cfg      string
	Executor string

	ExecuteTime int64 `gorm:"column:execute_time,index"`
	Status      int
	Cron        string
	Version     int64

	Ctime int64
	Utime int64
}

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	Done(ctx context.Context, id, version int64) error
	UpdateUtime(id int64) error
	UpdateExecTime(ctx context.Context, id int64, unix int64) error
	Stop(ctx context.Context, id int64) error
}

type GORMJobDAO struct {
	db *gorm.DB
}

func (g *GORMJobDAO) Stop(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Model(&Job{}).Where("id = ?", id).
		Updates(map[string]any{
			"status": jobStatusFinished,
			"utime":  time.Now().UnixMilli(),
		}).Error
}

func (g *GORMJobDAO) UpdateExecTime(ctx context.Context, id int64, unix int64) error {
	return g.db.WithContext(ctx).Model(&Job{}).Where("id = ?", id).
		Update("execute_time", unix).Error
}

func (g *GORMJobDAO) UpdateUtime(id int64) error {
	return g.db.Model(&Job{}).Where("id = ?", id).
		Update("utime", time.Now().UnixMilli()).Error
}

func (g *GORMJobDAO) Done(ctx context.Context, id, version int64) error {
	return g.db.WithContext(ctx).Model(&Job{}).Where("id = ? and version = ?", id, version).Updates(map[string]any{
		"status":  jobStatusFinished,
		"version": gorm.Expr("version + 1"),
		"utime":   time.Now().UnixMilli(),
	}).Error
}

func (g *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	for {
		now := time.Now()
		var job Job
		// todo 此处只取一个任务，高并发环境下容易冲突，所以可以优化为一次性取出多个任务，随机第一个执行任务的ID（防止多个goroutine发生冲突），提高抢占到任务执行的概率
		// 使用乐观锁代替FOR UPDATE
		err := g.db.WithContext(ctx).Model(&Job{}).Where("status = ? and execute_time <= ?", jobStatusWaiting, now).First(&job).Error
		if err != nil {
			return Job{}, err // 此时没有等待执行的Job，直接返回
		}

		res := g.db.WithContext(ctx).Model(&Job{}).Where("id = ? and version = ?", job.ID, job.Version).Updates(map[string]any{
			"status":  jobStatusRunning,
			"version": job.Version + 1,
			"utime":   now,
		})
		if res.Error != nil {
			return Job{}, res.Error
		}
		// 注:如果新数据==旧数据,则此处也为true
		if res.RowsAffected == 0 {
			// todo ！此处可以个人认为可以优化为类似ZK的公平锁机制，即本地维护抢占失败的goroutine进入队列，入队后goroutine睡眠，只有队首才去抢占下一个任务，成功后出队，并唤醒下一个goroutine
			continue // 抢占失败
		}

		return job, nil
	}
}
