package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	domain2 "post/interactive/domain"
	"post/internal/domain"
	svcmocks "post/internal/service/mock"
	"testing"
	"time"
)

func TestRankTopN(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (ArticleService, LikeService)

		wantErr  error
		wantArts []domain.Article
	}{
		{
			name: "成功",
			mock: func(ctrl *gomock.Controller) (ArticleService, LikeService) {
				artSvc := svcmocks.NewMockArticleService(ctrl)
				likeSvc := svcmocks.NewMockLikeService(ctrl)

				likeSvc.EXPECT().GetListAllOfLikes(gomock.Any(), gomock.Any(), 0, 100).Return([]domain2.Like{
					{
						ID:        1,
						LikeCount: 1,
					},
					{
						ID:        2,
						LikeCount: 2,
					},
					{
						ID:        3,
						LikeCount: 3,
					},
				}, nil)

				artSvc.EXPECT().GetPublishedByIDS(gomock.Any(), []int64{3, 2, 1}).Return([]domain.Article{
					{
						ID:      3,
						Content: "3",
						Ctime:   now,
					},
					{
						ID:      2,
						Content: "2",
						Ctime:   now,
					},
					{
						ID:      1,
						Content: "1",
						Ctime:   now,
					},
				}, nil)

				return artSvc, likeSvc
			},
			wantArts: []domain.Article{
				{
					ID:      3,
					Content: "3",
					Ctime:   now,
				},
				{
					ID:      2,
					Content: "2",
					Ctime:   now,
				},
				{
					ID:      1,
					Content: "1",
					Ctime:   now,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			artSvc, likeSvc := tc.mock(ctrl)
			svc := NewBatchRankService(artSvc, likeSvc)

			n, err := svc.GetRankTopN(context.Background(), 100)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantArts, n)
		})
	}
}
