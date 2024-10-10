package job

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"post/interactive/repository"
	"post/interactive/repository/cache"
	cacheMocks "post/interactive/repository/cache/mock"
	"post/interactive/repository/dao"
	"post/interactive/service"
	"testing"
)

//go:embed mysql.yaml
var mysqlDSN string

func TestBatchUpdateDBJob_Run(t *testing.T) {
	tests := []struct {
		name          string
		mockCacheData map[string]int64
		mockUpdateErr error
		expectedCalls int
		mock          func(ctrl *gomock.Controller) cache.ArticleLikeCache
		before        func()
		after         func()
	}{
		{
			name: "更新成功",
			mock: func(ctrl *gomock.Controller) cache.ArticleLikeCache {
				c := cacheMocks.NewMockArticleLikeCache(ctrl)
				c.EXPECT().GetCountByPrefix(gomock.Any(),
					"incr_read_count:article:").Return(map[string]int64{
					"incr_read_count:article:1": 5,
					"incr_read_count:article:2": 3,
				}, nil)
				return c
			},
			mockUpdateErr: nil,
			expectedCalls: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := tt.mock(ctrl)
			db, _ := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
			svc := service.NewLikeService(repository.NewLikeRepository(dao.NewGORMArticleLikeDao(db), c))
			job := NewBatchUpdateDBJobJob(svc, c)

			err := job.Run()

			if tt.mockUpdateErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
