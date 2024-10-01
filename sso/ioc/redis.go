package ioc

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"post/sso/config"
)

func InitRedis(info *config.Info) redis.Cmdable {
	cmd := redis.NewClient(&redis.Options{
		Addr: info.Config.Redis.Host + ":" + info.Config.Redis.Port,
	})

	_, err := cmd.Ping(context.Background()).Result()
	if err != nil {
		// 打印错误并返回
		panic(fmt.Errorf("failed to connect to Redis: %v", err))
	}

	return cmd
}
