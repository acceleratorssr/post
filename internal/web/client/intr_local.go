package client

import (
	"context"
	"google.golang.org/grpc"
	intrv1 "post/api/proto/gen/intr/v1"
	"post/interactive/domain"
	"post/interactive/service"
)

// LikeServiceAdapter 将本地实现伪装为grpc客户端，即实际调用本地方法
type LikeServiceAdapter struct {
	svc service.LikeService
}

func NewLikeServiceAdapter(svc service.LikeService) *LikeServiceAdapter {
	return &LikeServiceAdapter{
		svc: svc,
	}
}

func (l *LikeServiceAdapter) IncrReadCount(ctx context.Context, in *intrv1.IncrReadCountRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCountResponse, error) {
	err := l.svc.IncrReadCount(ctx, in.GetObjType(), in.GetObjID())
	if err != nil {
		return nil, err
	}
	return &intrv1.IncrReadCountResponse{
		Code:    200,
		Message: "success",
	}, nil
}

func (l *LikeServiceAdapter) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	err := l.svc.Like(ctx, in.GetObjType(), in.GetObjID(), in.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.LikeResponse{
		Code:    200,
		Message: "success",
	}, nil
}

func (l *LikeServiceAdapter) UnLike(ctx context.Context, in *intrv1.UnLikeRequest, opts ...grpc.CallOption) (*intrv1.UnLikeResponse, error) {
	err := l.svc.UnLike(ctx, in.GetObjType(), in.GetObjID(), in.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.UnLikeResponse{
		Code:    200,
		Message: "success",
	}, nil
}

func (l *LikeServiceAdapter) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	err := l.svc.Collect(ctx, in.GetObjType(), in.GetObjID(), in.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.CollectResponse{
		Code:    200,
		Message: "success",
	}, nil
}

func (l *LikeServiceAdapter) GetListBatchOfLikes(ctx context.Context, in *intrv1.GetListBatchOfLikesRequest, opts ...grpc.CallOption) (*intrv1.GetListBatchOfLikesResponse, error) {
	data, err := l.svc.GetListBatchOfLikes(ctx, in.GetObjType(), int(in.GetOffset()), int(in.GetLimit()), in.GetNow())
	if err != nil {
		return nil, err
	}
	return &intrv1.GetListBatchOfLikesResponse{
		Code:    200,
		Message: "success",
		Data:    l.domain2grpc(data),
	}, nil
}

func (l *LikeServiceAdapter) domain2grpc(like []domain.Like) []*intrv1.Like {
	intr := make([]*intrv1.Like, 0, len(like))
	for _, li := range like {
		intr = append(intr, &intrv1.Like{
			ID:        li.ID,
			Ctime:     li.Ctime,
			LikeCount: li.LikeCount,
		})
	}
	return intr
}
