package ioc

import (
	_ "embed"
	"github.com/spf13/viper"
	"github.com/zhenghaoz/gorse/client"
)

func InitGorse() *client.GorseClient {
	type Config struct {
		Addrs string `yaml:"addr"`
		Key   string `yaml:"key"`
	}

	var cfg Config
	err := viper.UnmarshalKey("gorse", &cfg)
	if err != nil {
		panic(err)
	}
	return client.NewGorseClient(cfg.Addrs, cfg.Key)
}
