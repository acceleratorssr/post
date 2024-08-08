package main

import (
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	intrv1 "post/api/proto/gen/intr/v1"
	"testing"
	"time"
)

func TestGRPCClient(t *testing.T) {
	c, err := grpc.NewClient("localhost:9200",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := intrv1.NewLikeServiceClient(c)
	resp, err := client.GetListBatchOfLikes(context.Background(), &intrv1.GetListBatchOfLikesRequest{
		ObjType: "article",
		Offset:  0,
		Limit:   10,
		Now:     time.Now().UnixMilli(),
	})
	require.NoError(t, err)
	t.Log(resp)
}
