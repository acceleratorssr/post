package scheduler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"post/migrator"
	"post/migrator/events"
	"post/migrator/validator"
	"post/pkg/gin_ex"
	"post/pkg/gorm_ex/connpool"
	"sync"
	"time"
)

type StartIncrRequest struct {
	Utime int64 `json:"utime"`
}

// Scheduler 用来统一管理整个迁移过程
type Scheduler[T migrator.Entity] struct {
	lock       sync.Mutex
	src        *gorm.DB
	dst        *gorm.DB
	pool       *connpool.DoubleWritePool
	pattern    string
	cancelFull func()
	cancelIncr func()
	producer   events.InconsistentProducer
}

func NewScheduler[T migrator.Entity](src *gorm.DB, dst *gorm.DB,
	pool *connpool.DoubleWritePool, producer events.InconsistentProducer) *Scheduler[T] {
	return &Scheduler[T]{
		src:     src,
		dst:     dst,
		pattern: connpool.PattenOnlyBase,
		cancelFull: func() { // 防nil
		},
		cancelIncr: func() {
		},
		pool:     pool,
		producer: producer,
	}
}

func (s *Scheduler[T]) RegisterRoutes(server *gin.RouterGroup) {
	server.POST("/base_only", s.OnlyBase)
	server.POST("/base_first", s.DependBase)
	server.POST("/target_first", s.DependTarget)
	server.POST("/target_only", s.OnlyTarget)
	server.POST("/full/start", s.StartFullValidation)
	server.POST("/full/stop", s.StopFullValidation)
	server.POST("/incr/stop", s.StopIncrementValidation)
	server.POST("/incr/start", gin_ex.WrapWithReq[StartIncrRequest](s.StartIncrementValidation))
}

func (s *Scheduler[T]) OnlyBase(c *gin.Context) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PattenOnlyBase
	s.pool.UpdatePatten(connpool.PattenOnlyBase)
	gin_ex.OKWithMessage("OK", c)
}

func (s *Scheduler[T]) DependBase(c *gin.Context) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PattenDependBase
	s.pool.UpdatePatten(connpool.PattenDependBase)
	gin_ex.OKWithMessage("OK", c)
}

func (s *Scheduler[T]) DependTarget(c *gin.Context) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PattenDependTarget
	s.pool.UpdatePatten(connpool.PattenDependTarget)
	gin_ex.OKWithMessage("OK", c)
}

func (s *Scheduler[T]) OnlyTarget(c *gin.Context) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PattenOnlyTarget
	s.pool.UpdatePatten(connpool.PattenOnlyTarget)
	gin_ex.OKWithMessage("OK", c)
}

func (s *Scheduler[T]) StopIncrementValidation(c *gin.Context) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cancelIncr()
	gin_ex.OKWithMessage("OK", c)
}

func (s *Scheduler[T]) StartIncrementValidation(c context.Context,
	req StartIncrRequest) (gin_ex.Response, error) {
	// 指定时间校验
	s.lock.Lock()
	defer s.lock.Unlock()

	cancel := s.cancelIncr
	v, err := s.newValidator(req.Utime)
	if err != nil {
		return gin_ex.Response{
			Code: 55,
			Msg:  "系统异常",
		}, err
	}

	var ctx context.Context
	ctx, s.cancelIncr = context.WithCancel(context.Background())

	go func() {
		cancel()
		v.Validate(ctx, 1000, 100)
	}()

	return gin_ex.Response{
		Msg: "启动增量校验成功",
	}, nil
}

func (s *Scheduler[T]) StopFullValidation(c *gin.Context) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cancelFull()
	gin_ex.OKWithMessage("OK", c)
}

// StartFullValidation 全量校验
func (s *Scheduler[T]) StartFullValidation(c *gin.Context) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// 取消上一次的
	cancel := s.cancelFull
	v, err := s.newValidator(time.Now().UnixMilli())
	if err != nil {
		return
	}
	var ctx context.Context
	ctx, s.cancelFull = context.WithCancel(context.Background())

	go func() {
		// 先取消上一次的
		cancel()
		v.Validate(ctx, 1000, 100)
	}()
	gin_ex.OKWithMessage("OK", c)
}

func (s *Scheduler[T]) newValidator(utime int64) (*validator.Validator[T], error) {
	switch s.pattern {
	case connpool.PattenDependBase, connpool.PattenOnlyBase:
		return validator.NewValidator[T](s.src, s.dst, s.producer, "base", utime), nil
	case connpool.PattenDependTarget, connpool.PattenOnlyTarget:
		return validator.NewValidator[T](s.dst, s.src, s.producer, "target", utime), nil
	default:
		return nil, fmt.Errorf("未知的 pattern %s", s.pattern)
	}
}
