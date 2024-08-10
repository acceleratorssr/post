package ioc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	intrv1 "post/api/proto/gen/intr/v1"
	"post/interactive/service"
	"post/internal/web/client"
)

func InitIntrGRPCClient(svc service.LikeService) intrv1.LikeServiceClient {
	local := client.NewLikeServiceAdapter(svc)
	c, err := grpc.NewClient("localhost:9200",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	remote := intrv1.NewLikeServiceClient(c)

	g := client.NewGreyScaleServiceAdapter(local, remote)
	g.UpdateThreshold(50) // 调整流量比例
	return g
}
