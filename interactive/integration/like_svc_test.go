package integration

import (
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	intrv1 "post/api/proto/gen/intr/v1"
	"post/interactive/grpc"
	"post/interactive/integration/startup"
	"post/interactive/repository/dao"
	"testing"
	"time"
)

type InteractiveTestSuite struct {
	suite.Suite
	db     *gorm.DB
	rdb    redis.Cmdable
	server *grpc.LikeServiceServer
}

func (s *InteractiveTestSuite) SetupSuite() {
	s.db = startup.InitDB()
	s.rdb = startup.InitRedis()
	s.server = startup.InitGRPCServer()
}

func (s *InteractiveTestSuite) TearDownTest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := s.db.Exec("TRUNCATE TABLE `likes`").Error
	assert.NoError(s.T(), err)
	err = s.db.Exec("TRUNCATE TABLE `user_give_likes`").Error
	assert.NoError(s.T(), err)
	err = s.db.Exec("TRUNCATE TABLE `user_give_collects`").Error
	assert.NoError(s.T(), err)
	// 清空 Redis
	err = s.rdb.FlushDB(ctx).Err()
	assert.NoError(s.T(), err)
}

func (s *InteractiveTestSuite) TestIncrReadCount() {
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		objType string
		objID   int64

		wantErr  error
		wantResp *intrv1.IncrReadCountResponse
	}{
		{
			// DB 和缓存都有数据
			name: "增加成功,db和redis",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.Create(dao.Like{
					ID:           1,
					ObjType:      "test",
					ObjID:        2,
					ViewCount:    3,
					CollectCount: 4,
					LikeCount:    5,
					Ctime:        6,
					Utime:        7,
				}).Error
				assert.NoError(t, err)
				err = s.rdb.HSet(ctx, "article_incr_read_count:test:2",
					"read_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Like
				err := s.db.Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 7)
				data.Utime = 0
				assert.Equal(t, dao.Like{
					ID:           1,
					ObjType:      "test",
					ObjID:        2,
					ViewCount:    4,
					CollectCount: 4,
					LikeCount:    5,
					Ctime:        6,
				}, data)
				cnt, err := s.rdb.HGet(ctx, "article_incr_read_count:test:2", "read_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 4, cnt)
				err = s.rdb.Del(ctx, "article_incr_read_count:test:2").Err()
				assert.NoError(t, err)
			},
			objType: "test",
			objID:   2,
			wantResp: &intrv1.IncrReadCountResponse{
				Code:    200,
				Message: "success",
			},
		},
		{
			// DB 有数据，缓存没有数据
			name: "增加成功,db有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.WithContext(ctx).Create(dao.Like{
					ID:           3,
					ObjType:      "test",
					ObjID:        3,
					ViewCount:    3,
					CollectCount: 4,
					LikeCount:    5,
					Ctime:        6,
					Utime:        7,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Like
				err := s.db.Where("id = ?", 3).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 7)
				data.Utime = 0
				assert.Equal(t, dao.Like{
					ID:      3,
					ObjType: "test",
					ObjID:   3,
					// +1 之后
					ViewCount:    4,
					CollectCount: 4,
					LikeCount:    5,
					Ctime:        6,
				}, data)
				cnt, err := s.rdb.Exists(ctx, "article_incr_read_count:test:3").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
			},
			objType: "test",
			objID:   3,
			wantResp: &intrv1.IncrReadCountResponse{
				Code:    200,
				Message: "success",
			},
		},
		{
			name:   "增加成功-都没有",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Like
				err := s.db.Where("obj_type = ? AND obj_id = ?", "test", 4).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 0)
				assert.True(t, data.Ctime > 0)
				assert.True(t, data.ID > 0)
				data.ID = 0
				data.Utime = 0
				data.Ctime = 0
				assert.Equal(t, dao.Like{
					ObjType:   "test",
					ObjID:     4,
					ViewCount: 1,
				}, data)
				cnt, err := s.rdb.Exists(ctx, "article_incr_read_count:test:4").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
			},
			objType: "test",
			objID:   4,
			wantResp: &intrv1.IncrReadCountResponse{
				Code:    200,
				Message: "success",
			},
		},
	}

	// 不同于 AsyncSms 服务，我们不需要 mock，所以创建一个就可以
	// 不需要每个测试都创建
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			resp, err := s.server.IncrReadCount(context.Background(), &intrv1.IncrReadCountRequest{
				ObjType: tc.objType, ObjID: tc.objID,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)
			tc.after(t)
		})
	}
}

func (s *InteractiveTestSuite) TestLike() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		objType string
		objID   int64
		uid     int64

		wantErr  error
		wantResp *intrv1.LikeResponse
	}{
		{
			name: "点赞-DB和cache都有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.Create(dao.Like{
					ID:           1,
					ObjType:      "test",
					ObjID:        2,
					ViewCount:    3,
					CollectCount: 4,
					LikeCount:    5,
					Ctime:        6,
					Utime:        7,
				}).Error
				assert.NoError(t, err)
				err = s.rdb.HSet(ctx, "article_incr_Like_count:test:2",
					"like_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Like
				err := s.db.Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 7)
				data.Utime = 0
				assert.Equal(t, dao.Like{
					ID:           1,
					ObjType:      "test",
					ObjID:        2,
					ViewCount:    3,
					CollectCount: 4,
					LikeCount:    6,
					Ctime:        6,
				}, data)

				var likeBiz dao.UserGiveLike
				err = s.db.Where("obj_type = ? AND obj_id = ? AND uid = ?",
					"test", 2, 123).First(&likeBiz).Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.ID > 0)
				assert.True(t, likeBiz.Ctime > 0)
				assert.True(t, likeBiz.Utime > 0)
				likeBiz.ID = 0
				likeBiz.Ctime = 0
				likeBiz.Utime = 0
				assert.Equal(t, dao.UserGiveLike{
					ObjType: "test",
					ObjID:   2,
					Uid:     123,
					Status:  0, // 0点赞 1取消
				}, likeBiz)

				cnt, err := s.rdb.HGet(ctx, "article_incr_Like_count:test:2", "like_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 4, cnt)
				err = s.rdb.Del(ctx, "article_incr_Like_count:test:2").Err()
				assert.NoError(t, err)
			},
			objType: "test",
			objID:   2,
			uid:     123,
			wantResp: &intrv1.LikeResponse{
				Code:    200,
				Message: "success",
			},
		},
		{
			name:   "点赞-都没有",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Like
				err := s.db.Where("obj_type = ? AND obj_id = ?", "test", 3).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 0)
				assert.True(t, data.Ctime > 0)
				assert.True(t, data.ID > 0)
				data.Utime = 0
				data.Ctime = 0
				data.ID = 0
				assert.Equal(t, dao.Like{
					ObjType:   "test",
					ObjID:     3,
					LikeCount: 1,
				}, data)

				var likeBiz dao.UserGiveLike
				err = s.db.Where("obj_type = ? AND obj_id = ? AND uid = ?",
					"test", 3, 123).First(&likeBiz).Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.ID > 0)
				assert.True(t, likeBiz.Ctime > 0)
				assert.True(t, likeBiz.Utime > 0)
				likeBiz.ID = 0
				likeBiz.Ctime = 0
				likeBiz.Utime = 0
				assert.Equal(t, dao.UserGiveLike{
					ObjType: "test",
					ObjID:   3,
					Uid:     123,
					Status:  0,
				}, likeBiz)

				cnt, err := s.rdb.Exists(ctx, "article_incr_Like_count:test:2").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
			},
			objType: "test",
			objID:   3,
			uid:     123,
			wantResp: &intrv1.LikeResponse{
				Code:    200,
				Message: "success",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			resp, err := s.server.Like(context.Background(), &intrv1.LikeRequest{
				ObjType: tc.objType, ObjID: tc.objID, Uid: tc.uid,
			})
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResp, resp)
			tc.after(t)
		})
	}
}

func (s *InteractiveTestSuite) TestUnLike() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		objType string
		objID   int64
		uid     int64

		wantErr  error
		wantResp *intrv1.UnLikeResponse
	}{
		{
			name: "取消点赞-DB和cache都有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.Create(dao.Like{
					ID:           1,
					ObjType:      "test",
					ObjID:        2,
					ViewCount:    3,
					CollectCount: 4,
					LikeCount:    5,
					Ctime:        6,
					Utime:        7,
				}).Error
				assert.NoError(t, err)
				err = s.db.Create(dao.UserGiveLike{
					ID:      1,
					ObjType: "test",
					ObjID:   2,
					Uid:     123,
					Ctime:   6,
					Utime:   7,
					Status:  0,
				}).Error
				assert.NoError(t, err)
				err = s.rdb.HSet(ctx, "article_incr_Like_count:test:2",
					"like_cnt", 5).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Like
				err := s.db.Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 7)
				data.Utime = 0
				assert.Equal(t, dao.Like{
					ID:           1,
					ObjType:      "test",
					ObjID:        2,
					ViewCount:    3,
					CollectCount: 4,
					LikeCount:    4,
					Ctime:        6,
				}, data)

				var likeType dao.UserGiveLike
				err = s.db.Where("id = ?", 1).First(&likeType).Error
				assert.NoError(t, err)
				assert.True(t, likeType.Utime > 7)
				likeType.Utime = 0
				assert.Equal(t, dao.UserGiveLike{
					ID:      1,
					ObjType: "test",
					ObjID:   2,
					Uid:     123,
					Ctime:   6,
					Status:  1,
				}, likeType)

				cnt, err := s.rdb.HGet(ctx, "article_incr_Like_count:test:2", "like_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 4, cnt)
				err = s.rdb.Del(ctx, "article_incr_Like_count:test:2").Err()
				assert.NoError(t, err)
			},
			objType: "test",
			objID:   2,
			uid:     123,
			wantResp: &intrv1.UnLikeResponse{
				Code:    200,
				Message: "success",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			resp, err := s.server.UnLike(context.Background(), &intrv1.UnLikeRequest{
				ObjType: tc.objType, ObjID: tc.objID, Uid: tc.uid,
			})
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResp, resp)
			tc.after(t)
		})
	}
}

func (s *InteractiveTestSuite) TestCollect() {
	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		objType string
		objID   int64
		cid     int64
		uid     int64

		wantErr  error
		wantResp *intrv1.CollectResponse
	}{
		{
			name:   "收藏成功",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				var intr dao.Like
				err := s.db.Where("obj_type = ? AND obj_id = ?", "test", 1).First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime > 0)
				intr.Ctime = 0
				assert.True(t, intr.Utime > 0)
				intr.Utime = 0
				assert.True(t, intr.ID > 0)
				intr.ID = 0
				assert.Equal(t, dao.Like{
					ObjType:      "test",
					ObjID:        1,
					CollectCount: 1,
				}, intr)

				// 收藏记录
				var cbiz dao.UserGiveCollect
				err = s.db.WithContext(ctx).
					Where("uid = ? AND obj_type = ? AND obj_id = ?", 1, "test", 1).
					First(&cbiz).Error
				assert.NoError(t, err)
				assert.True(t, cbiz.Ctime > 0)
				cbiz.Ctime = 0
				assert.True(t, cbiz.Utime > 0)
				cbiz.Utime = 0
				assert.True(t, cbiz.ID > 0)
				cbiz.ID = 0
				assert.Equal(t, dao.UserGiveCollect{
					ObjType: "test",
					ObjID:   1,
					Uid:     1,
				}, cbiz)
			},
			objID:   1,
			objType: "test",
			cid:     1,
			uid:     1,
			wantResp: &intrv1.CollectResponse{
				Code:    200,
				Message: "success",
			},
		},
		{
			name: "收藏成功,db有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := s.db.WithContext(ctx).Create(&dao.Like{
					ObjType:      "test",
					ObjID:        2,
					CollectCount: 10,
					Ctime:        123,
					Utime:        234,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				var intr dao.Like
				err := s.db.WithContext(ctx).
					Where("obj_type = ? AND obj_id = ?", "test", 2).First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime > 0)
				intr.Ctime = 0
				assert.True(t, intr.Utime > 0)
				intr.Utime = 0
				assert.True(t, intr.ID > 0)
				intr.ID = 0
				assert.Equal(t, dao.Like{
					ObjType:      "test",
					ObjID:        2,
					CollectCount: 11,
				}, intr)

				var cbiz dao.UserGiveCollect
				err = s.db.WithContext(ctx).
					Where("uid = ? AND obj_type = ? AND obj_id = ?", 1, "test", 2).
					First(&cbiz).Error
				assert.NoError(t, err)
				assert.True(t, cbiz.Ctime > 0)
				cbiz.Ctime = 0
				assert.True(t, cbiz.Utime > 0)
				cbiz.Utime = 0
				assert.True(t, cbiz.ID > 0)
				cbiz.ID = 0
				assert.Equal(t, dao.UserGiveCollect{
					ObjType: "test",
					ObjID:   2,
					Uid:     1,
				}, cbiz)
			},
			objID:   2,
			objType: "test",
			cid:     1,
			uid:     1,
			wantResp: &intrv1.CollectResponse{
				Code:    200,
				Message: "success",
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			resp, err := s.server.Collect(context.Background(), &intrv1.CollectRequest{
				ObjType: tc.objType, ObjID: tc.objID, Uid: tc.uid,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)
			tc.after(t)
		})
	}
}

//func (s *InteractiveTestSuite) TestGet() {
//	testCases := []struct {
//		name string
//
//		before func(t *testing.T)
//
//		objID int64
//		objType   string
//		uid   int64
//
//		wantErr error
//		wantRes *intrv1.GetResponse
//	}{
//		{
//			name:  "全部取出来了-无缓存",
//			objType:   "test",
//			objID: 12,
//			uid:   123,
//			before: func(t *testing.T) {
//				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
//				defer cancel()
//				err := s.db.WithContext(ctx).Create(&dao.Like{
//					ObjType:      "test",
//					ObjID:        12,
//					ViewCount:    100,
//					CollectCount: 200,
//					LikeCount:    300,
//					Ctime:        123,
//					Utime:        234,
//				}).Error
//				assert.NoError(t, err)
//			},
//			wantRes: &intrv1.GetResponse{
//				Intr: &intrv1.Interactive{
//					Biz:        "test",
//					BizId:      12,
//					ReadCnt:    100,
//					CollectCnt: 200,
//					LikeCnt:    300,
//				},
//			},
//		},
//		{
//			name:  "全部取出来了-命中缓存-用户已点赞收藏",
//			objType:   "test",
//			objID: 3,
//			uid:   123,
//			before: func(t *testing.T) {
//				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
//				defer cancel()
//				err := s.db.WithContext(ctx).
//					Create(&dao.UserGiveCollect{
//						ObjType: "test",
//						ObjID:   3,
//						Uid:     123,
//						Ctime:   123,
//						Utime:   124,
//					}).Error
//				assert.NoError(t, err)
//				err = s.db.WithContext(ctx).
//					Create(&dao.UserGiveLike{
//						ObjType: "test",
//						ObjID:   3,
//						Uid:     123,
//						Ctime:   123,
//						Utime:   124,
//						Status:  1,
//					}).Error
//				assert.NoError(t, err)
//				err = s.rdb.HSet(ctx, "interactive:test:3",
//					"read_cnt", 0, "collect_cnt", 1).Err()
//				assert.NoError(t, err)
//			},
//			wantRes: &intrv1.GetResponse{
//				Intr: &intrv1.Interactive{
//					BizId:      3,
//					CollectCnt: 1,
//					Collected:  true,
//					Liked:      true,
//				},
//			},
//		},
//	}
//	for _, tc := range testCases {
//		s.T().Run(tc.name, func(t *testing.T) {
//			tc.before(t)
//			res, err := s.server.Get(context.Background(), &intrv1.GetRequest{
//				Biz: tc.objType, BizId: tc.objID, Uid: tc.uid,
//			})
//			assert.Equal(t, tc.wantErr, err)
//			assert.Equal(t, tc.wantRes, res)
//		})
//	}
//}

func (s *InteractiveTestSuite) TestGetListBatchOfLikes() {
	preCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// 准备数据
	for i := 1; i < 5; i++ {
		i := int64(i)
		err := s.db.WithContext(preCtx).
			Create(&dao.Like{
				ID:           i,
				ObjType:      "test",
				ObjID:        i,
				ViewCount:    i,
				CollectCount: i + 1,
				LikeCount:    i + 2,
			}).Error
		assert.NoError(s.T(), err)
	}

	now := time.Now().UnixMilli()

	testCases := []struct {
		name string

		before  func(t *testing.T)
		objType string
		limit   int32
		offset  int32
		now     int64

		wantErr error
		wantRes *intrv1.GetListBatchOfLikesResponse
	}{
		{
			name:    "查找成功",
			objType: "test",
			limit:   2,
			offset:  0,
			now:     now,
			// todo 时间也可以加入考虑
			wantRes: &intrv1.GetListBatchOfLikesResponse{
				Code:    200,
				Message: "success",
				Data: []*intrv1.Like{
					{
						ID:        1,
						LikeCount: 3,
					},
					{
						ID:        2,
						LikeCount: 4,
					},
				},
			},
		},
		{
			name:    "没有对应的数据",
			objType: "test",
			wantRes: &intrv1.GetListBatchOfLikesResponse{
				Code:    200,
				Message: "success",
				Data:    []*intrv1.Like{},
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			res, err := s.server.GetListBatchOfLikes(context.Background(), &intrv1.GetListBatchOfLikesRequest{
				ObjType: tc.objType, Offset: tc.offset, Limit: tc.limit, Now: tc.now,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestInteractiveService(t *testing.T) {
	suite.Run(t, &InteractiveTestSuite{})
}
