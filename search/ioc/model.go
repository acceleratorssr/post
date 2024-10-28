package ioc

import (
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"os"
)

func InitMarsCode() *arkruntime.Client {
	return arkruntime.NewClientWithApiKey(
		os.Getenv("ARK_API_KEY"),
		arkruntime.WithBaseUrl("https://ark.cn-beijing.volces.com/api/v3"),
		arkruntime.WithRegion("cn-beijing"),
	)
}
