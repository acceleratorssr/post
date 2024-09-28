package ioc

import (
	"github.com/redis/go-redis/v9"
	"post/sso/config"
)

func InitRedis(info *config.Info) redis.Cmdable {
	cmd := redis.NewClient(&redis.Options{
		Addr: info.Config.Redis.Host + info.Config.Redis.Port,
	})
	return cmd
}
