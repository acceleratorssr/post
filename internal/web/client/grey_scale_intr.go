package client

import (
	"context"
	"google.golang.org/grpc"
	"math/rand"
	intrv1 "post/api/proto/gen/intr/v1"
	"sync/atomic"
)

// GreyScaleServiceAdapter
// 通过装饰器，0~100设置阈值，控制流量比例，为100则全走本地
type GreyScaleServiceAdapter struct {
	remote intrv1.LikeServiceClient
	local  intrv1.LikeServiceClient

	threshold *atomic.Int32
}

func NewGreyScaleServiceAdapter(local, remote intrv1.LikeServiceClient) *GreyScaleServiceAdapter {
	return &GreyScaleServiceAdapter{
		local:     local,
		remote:    remote,
		threshold: new(atomic.Int32),
	}
}

func (g *GreyScaleServiceAdapter) IncrReadCount(ctx context.Context, in *intrv1.IncrReadCountRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCountResponse, error) {
	return g.client().IncrReadCount(ctx, in, opts...)
}

func (g *GreyScaleServiceAdapter) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	return g.client().Like(ctx, in, opts...)
}

func (g *GreyScaleServiceAdapter) UnLike(ctx context.Context, in *intrv1.UnLikeRequest, opts ...grpc.CallOption) (*intrv1.UnLikeResponse, error) {
	return g.client().UnLike(ctx, in, opts...)
}

func (g *GreyScaleServiceAdapter) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	return g.client().Collect(ctx, in, opts...)
}

func (g *GreyScaleServiceAdapter) GetListBatchOfLikes(ctx context.Context, in *intrv1.GetListBatchOfLikesRequest, opts ...grpc.CallOption) (*intrv1.GetListBatchOfLikesResponse, error) {
	return g.client().GetListBatchOfLikes(ctx, in, opts...)
}

func (g *GreyScaleServiceAdapter) UpdateThreshold(newValue int32) {
	if newValue > 100 {
		newValue = 100
	} else if newValue < 0 {
		newValue = 0
	}
	g.threshold.Store(newValue)
}

func (g *GreyScaleServiceAdapter) client() intrv1.LikeServiceClient {
	threshold := g.threshold.Load()
	num := rand.Int31n(100)
	if num < threshold {
		return g.local
	} else {
		return g.remote
	}
}
