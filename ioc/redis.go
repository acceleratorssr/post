package ioc

import (
	_ "embed"
	"github.com/redis/go-redis/v9"
)

//go:embed redis.yaml
var addr_redis string

func InitRedis() redis.Cmdable {
	cmd := redis.NewClient(&redis.Options{
		Addr: addr_redis,
	})
	return cmd
}
