package grpc

import (
	"context"
	"github.com/zhenghaoz/gorse/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	recommendv1 "post/api/proto/gen/recommend/v1"
	"post/recommend/service"
)

type RecommendServiceServer struct {
	recommendv1.UnimplementedRecommendServiceServer
	svc service.RecommendService
}

func (r *RecommendServiceServer) GetItemByID(ctx context.Context, request *recommendv1.GetItemByIDRequest) (*recommendv1.GetItemByIDResponse, error) {
	item, err := r.svc.GetItemByID(ctx, request.GetItemId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "文章尚未被添加过推荐系统")
	}

	return &recommendv1.GetItemByIDResponse{
		Item: &recommendv1.Item{
			ItemId:     item.ItemId,
			IsHidden:   item.IsHidden,
			Labels:     item.Labels,
			Categories: item.Categories,
			Timestamp:  item.Timestamp,
			Comment:    item.Comment,
		},
	}, nil
}

func (r *RecommendServiceServer) GetNeighbors(ctx context.Context, request *recommendv1.GetNeighborsRequest) (*recommendv1.GetNeighborsResponse, error) {
	neighbors, err := r.svc.GetNeighbors(ctx, request.GetItemId(), request.GetUserId(), int(request.GetN()), int(request.GetOffset()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "文章暂无相似数据")
	}
	return &recommendv1.GetNeighborsResponse{
		Scores: r.toScoreDTO(neighbors),
	}, nil
}

// GetRecommend 返回的是 item IDs
func (r *RecommendServiceServer) GetRecommend(ctx context.Context, request *recommendv1.GetRecommendRequest) (*recommendv1.GetRecommendResponse, error) {
	recommend, err := r.svc.GetRecommend(ctx, request.GetUserId(),
		request.GetWriteBackType(), request.GetWriteBackDelay(),
		int(request.GetN()), int(request.GetOffset()), request.GetCategory()...)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "文章暂无相关推荐")
	}

	return &recommendv1.GetRecommendResponse{
		Items: recommend,
	}, nil
}

func (r *RecommendServiceServer) GetUser(ctx context.Context, request *recommendv1.GetUserRequest) (*recommendv1.GetUserResponse, error) {
	user, err := r.svc.GetUser(ctx, request.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "文章暂无相关推荐")
	}

	return &recommendv1.GetUserResponse{
		User: r.toUserDTO([]client.User{user}...)[0],
	}, nil
}

func (r *RecommendServiceServer) GetUsers(ctx context.Context, request *recommendv1.GetUsersRequest) (*recommendv1.GetUsersResponse, error) {
	user, err := r.svc.GetUsers(ctx, request.GetCursor(), int(request.GetN()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "文章暂无相关推荐")
	}

	return &recommendv1.GetUsersResponse{
		Users: &recommendv1.Users{
			Users:  r.toUserDTO(user.Users...),
			Cursor: user.Cursor,
		},
	}, nil
}

func (r *RecommendServiceServer) toScoreDTO(neighbors []client.Score) []*recommendv1.Score {
	var scores []*recommendv1.Score
	for _, neighbor := range neighbors {
		scores = append(scores, &recommendv1.Score{
			ID:    neighbor.Id,
			Score: int32(neighbor.Score),
		})
	}
	return scores
}

func (r *RecommendServiceServer) toUserDTO(user ...client.User) []*recommendv1.User {
	var users []*recommendv1.User
	for _, label := range user {
		users = append(users, &recommendv1.User{
			UserId:    label.UserId,
			Labels:    label.Labels,
			Subscribe: label.Subscribe,
			Comment:   label.Comment,
		})
	}

	return users
}

func (r *RecommendServiceServer) RegisterServer(server *grpc.Server) {
	recommendv1.RegisterRecommendServiceServer(server, r)
}

func NewRecommendServiceServer(svc service.RecommendService) *RecommendServiceServer {
	return &RecommendServiceServer{
		svc: svc,
	}
}
