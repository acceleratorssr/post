package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	intrv1 "post/api/proto/gen/intr/v1"
	"post/interactive/domain"
	"post/interactive/service"
)

// LikeServiceServer 将service包装为grpc暴露出去
// 即此处不考虑调用方
type LikeServiceServer struct {
	intrv1.UnimplementedLikeServiceServer
	svc service.LikeService
}

func NewLikeServiceServer(svc service.LikeService) *LikeServiceServer {
	return &LikeServiceServer{
		svc: svc,
	}
}

func (l *LikeServiceServer) Register(server *grpc.Server) {
	intrv1.RegisterLikeServiceServer(server, l)
}

// IncrReadCount GetObjType 防止req为空
func (l *LikeServiceServer) IncrReadCount(ctx context.Context, request *intrv1.IncrReadCountRequest) (*intrv1.IncrReadCountResponse, error) {
	err := l.svc.IncrReadCount(ctx, request.GetObjType(), request.GetObjID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "interactive 增长阅读数失败: %s", err)
	}
	return &intrv1.IncrReadCountResponse{}, nil
}

func (l *LikeServiceServer) Like(ctx context.Context, request *intrv1.LikeRequest) (*intrv1.LikeResponse, error) {
	// 参数校验，有自动生成校验的grpc插件
	if request.Uid <= 0 {
		return nil, status.Error(codes.InvalidArgument, "interactive 用户ID非法")
	}
	err := l.svc.Like(ctx, request.GetObjType(), request.GetObjID(), request.GetUid())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "interactive 点赞失败: %s", err)
	}
	return &intrv1.LikeResponse{}, nil
}

func (l *LikeServiceServer) UnLike(ctx context.Context, request *intrv1.UnLikeRequest) (*intrv1.UnLikeResponse, error) {
	err := l.svc.UnLike(ctx, request.GetObjType(), request.GetObjID(), request.GetUid())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "interactive 取消点赞失败: %s", err)
	}
	return &intrv1.UnLikeResponse{}, nil
}

func (l *LikeServiceServer) Collect(ctx context.Context, request *intrv1.CollectRequest) (*intrv1.CollectResponse, error) {
	err := l.svc.Collect(ctx, request.GetObjType(), request.GetObjID(), request.GetUid())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "interactive 收藏失败: %s", err)
	}
	return &intrv1.CollectResponse{}, nil
}

func (l *LikeServiceServer) GetListBatchOfLikes(ctx context.Context, request *intrv1.GetListBatchOfLikesRequest) (*intrv1.GetListBatchOfLikesResponse, error) {
	data, err := l.svc.GetListBatchOfLikes(ctx, request.GetObjType(), &domain.List{
		Desc:      request.GetDesc(),
		LastValue: request.GetLastValue(),
		Limit:     int(request.GetLimit()),
		OrderBy:   request.GetOrderBy(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "interactive 获取指定类型的一批数据的点赞数失败: %s", err)
	}
	return &intrv1.GetListBatchOfLikesResponse{
		Data: l.toDTO(data...),
	}, nil
}

// data transfer object 数据传输对象
func (l *LikeServiceServer) toDTO(intr ...domain.Like) []*intrv1.Like {
	data := make([]*intrv1.Like, 0) // 声明并初始化为空切片
	for _, v := range intr {
		data = append(data, &intrv1.Like{
			ID:        v.ID,
			LikeCount: v.LikeCount,
			Ctime:     v.Ctime,
		})
	}
	return data
}
