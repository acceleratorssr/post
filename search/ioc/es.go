package ioc

import (
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"post/search/repository/dao"
	"time"
)

// InitESClient 读取配置文件，进行初始化ES客户端
func InitESClient() *elastic.Client {
	type Config struct {
		Urls  string `yaml:"urls"`
		Sniff bool   `yaml:"sniff"`
	}
	var cfg Config
	err := viper.UnmarshalKey("es", &cfg)
	if err != nil {
		panic(err)
	}

	const timeout = 10 * time.Second
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(cfg.Urls),
		elastic.SetSniff(cfg.Sniff),
		elastic.SetHealthcheckTimeoutStartup(timeout),
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		panic(err)
	}

	err = dao.InitES(client)
	if err != nil {
		panic(err)
	}
	return client
}
