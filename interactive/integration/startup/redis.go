package startup

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
)

var redisClient redis.Cmdable

//go:embed redis.yaml
var addrRedis string

func InitRedis() redis.Cmdable {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr: addrRedis,
		})

		for err := redisClient.Ping(context.Background()).Err(); err != nil; {
			panic(err)
		}
	}
	return redisClient
}
