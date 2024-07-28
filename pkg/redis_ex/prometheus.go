package redis_ex

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"net"
	"strconv"
	"time"
)

type PrometheusRedisHook struct {
	vector *prometheus.SummaryVec
}

// NewPrometheusRedisHook 创建redis的hook，调用.AddHook(hook)即可注册
func NewPrometheusRedisHook(opt prometheus.SummaryOpts) *PrometheusRedisHook {
	vector := prometheus.NewSummaryVec(opt,
		[]string{"cmd", "key_exist"}) // 命令、是否命中缓存
	prometheus.MustRegister(vector)

	return &PrometheusRedisHook{
		vector: vector,
	}
}

func (p *PrometheusRedisHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr) //不做操作
	}
}

// ProcessHook
// todo hook可以用于拦截命令后，先查询本地缓存，如果命中则直接返回，否则才发送命令到redis上执行
func (p *PrometheusRedisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		var err error
		defer func() {
			// 统计redis命令执行时间
			duration := time.Since(start).Milliseconds()
			keyExist := err == redis.Nil
			p.vector.WithLabelValues(cmd.Name(), strconv.FormatBool(keyExist)).Observe(float64(duration))
		}()
		// next后即命令被发送到redis上执行
		err = next(ctx, cmd)
		return err
	}
}

func (p *PrometheusRedisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds) //不做操作
	}
}
