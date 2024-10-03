//go:build ignore

package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"post/domain"
	"post/repository"
	articleRepoMock "post/repository/mock"
	"testing"
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

const (
	userId   = 99
	timeNow  = 99
	makingID = 99
)

func TestArticleService_Publish(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (
			repository.ArticleAuthorRepository, repository.ArticleReaderRepository)
		art      domain.Article
		before   func()
		after    func()
		wantBody Result
		wantErr  error
		wantId   int64
	}{
		{
			name: "创建&发表成功",
			mock: func(ctrl *gomock.Controller) (
				repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
				author := articleRepoMock.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "title",
					Content: "content",
					Author: domain.Author{
						Id:   userId,
						Name: "test",
					},
					Status: domain.TypePublished,
				}).Return(int64(makingID), nil)
				reader := articleRepoMock.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					ID:      makingID,
					Title:   "title",
					Content: "content",
					Author: domain.Author{
						Id:   userId,
						Name: "test",
					},
					Status: domain.TypePublished,
				}).Return(int64(makingID), nil)
				return author, reader
			},
			art: domain.Article{
				Title:   "title",
				Content: "content",
				Author: domain.Author{
					Id:   userId,
					Name: "test",
				},
			},
			wantErr: nil,
			wantId:  makingID,
		},
		{
			name: "修改&发表成功",
			mock: func(ctrl *gomock.Controller) (
				repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
				author := articleRepoMock.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					ID:      makingID,
					Title:   "title",
					Content: "content",
					Author: domain.Author{
						Id:   userId,
						Name: "test",
					},
					Status: domain.TypePublished,
				}).Return(nil)
				reader := articleRepoMock.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					ID:      makingID,
					Title:   "title",
					Content: "content",
					Author: domain.Author{
						Id:   userId,
						Name: "test",
					},
					Status: domain.TypePublished,
				}).Return(int64(makingID), nil)
				return author, reader
			},
			art: domain.Article{
				ID:      makingID,
				Title:   "title",
				Content: "content",
				Author: domain.Author{
					Id:   userId,
					Name: "test",
				},
			},
			wantErr: nil,
			wantId:  makingID,
		},
		{
			name: "保存到草稿库失败", // todo 新建&修改
			mock: func(ctrl *gomock.Controller) (
				repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
				author := articleRepoMock.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					ID:      makingID,
					Title:   "title",
					Content: "content",
					Author: domain.Author{
						Id:   userId,
						Name: "test",
					},
					Status: domain.TypePublished,
				}).Return(errors.New("save to making failed"))

				reader := articleRepoMock.NewMockArticleReaderRepository(ctrl)
				return author, reader
			},
			art: domain.Article{
				ID:      makingID,
				Title:   "title",
				Content: "content",
				Author: domain.Author{
					Id:   userId,
					Name: "test",
				},
			},
			wantErr: errors.New("save to making failed"),
			wantId:  0,
		},
		{
			name: "保存到草稿库成功，线上库失败，即发表失败",
			mock: func(ctrl *gomock.Controller) (
				repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
				author := articleRepoMock.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					ID:      makingID,
					Title:   "title",
					Content: "content",
					Author: domain.Author{
						Id:   userId,
						Name: "test",
					},
					Status: domain.TypePublished,
				}).Return(nil)
				reader := articleRepoMock.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					ID:      makingID,
					Title:   "title",
					Content: "content",
					Author: domain.Author{
						Id:   userId,
						Name: "test",
					},
					Status: domain.TypePublished,
				}).Return(int64(0), errors.New("save to online failed")).AnyTimes()

				return author, reader
			},
			art: domain.Article{
				ID:      makingID,
				Title:   "title",
				Content: "content",
				Author: domain.Author{
					Id:   userId,
					Name: "test",
				},
			},
			wantErr: errors.New("save to online failed"),
			wantId:  0,
		},
		{
			name: "保存到草稿库成功，线上库失败，即发表失败",
			mock: func(ctrl *gomock.Controller) (
				repository.ArticleAuthorRepository, repository.ArticleReaderRepository) {
				author := articleRepoMock.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					ID:      makingID,
					Title:   "title",
					Content: "content",
					Author: domain.Author{
						Id:   userId,
						Name: "test",
					},
					Status: domain.TypePublished,
				}).Return(nil)
				reader := articleRepoMock.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					ID:      makingID,
					Title:   "title",
					Content: "content",
					Author: domain.Author{
						Id:   userId,
						Name: "test",
					},
					Status: domain.TypePublished,
				}).Return(int64(0), errors.New("save to online failed"))

				// 第二次调用，以此类推
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					ID:      makingID,
					Title:   "title",
					Content: "content",
					Author: domain.Author{
						Id:   userId,
						Name: "test",
					},
					Status: domain.TypePublished,
				}).Return(int64(makingID), nil)
				return author, reader
			},

			art: domain.Article{
				ID:      makingID,
				Title:   "title",
				Content: "content",
				Author: domain.Author{
					Id:   userId,
					Name: "test",
				},
			},
			wantErr: nil,
			wantId:  userId,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//tc.before()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			author, reader := tc.mock(ctrl)
			svc := NewArticleService(author, reader)
			id, err := svc.Publish(context.Background(), tc.art)
			require.Equal(t, tc.wantErr, err)
			require.Equal(t, tc.wantId, id)

			//tc.after()
		})
	}
}
