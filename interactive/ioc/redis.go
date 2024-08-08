package ioc

import (
	_ "embed"
	"github.com/redis/go-redis/v9"
)

//go:embed redis.yaml
var addrRedis string

func InitRedis() redis.Cmdable {
	cmd := redis.NewClient(&redis.Options{
		Addr: addrRedis,
	})
	return cmd
}
