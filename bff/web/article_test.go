package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"post/article/domain"
	"post/article/service"
	svcmocks "post/article/service/mock"
	"post/article/user"
	"testing"
)

const (
	userId  = 99
	timeNow = 99
)

type Result struct {
	Code int
	Data any
	Msg  string
}

// 奇怪，goland运行就显示未运行任何测试，命令行执行就正常。第二天又正常了，，
func TestArticleHandler_PublishArticle(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.ArticleService
		reqBody  string
		wantCode int
		wantBody Result
	}{
		{
			name: "新建&发表",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "title",
					Content: "content",
					Author: domain.Author{
						Id:   userId,
						Name: "test",
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody: `
			{
				"title":"title",
				"content":"content"
			}`,
			wantCode: http.StatusOK,
			wantBody: Result{
				Code: http.StatusOK,
				Data: float64(1), //json转any类型的数字，会变为float64类型
				Msg:  "success",
			},
		},
		{
			name: "发表失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "title",
					Content: "content",
					Author: domain.Author{
						Id:   userId,
						Name: "test",
					},
				}).Return(int64(0), errors.New("publish failed"))
				return svc
			},
			reqBody: `
			{
				"title":"title",
				"content":"content"
			}`,
			wantCode: 555,
			wantBody: Result{
				Code: 555,
				Data: "publish failed",
				Msg:  "系统错误",
			},
		},
		// todo 修改已有的帖子，发表成功
		// todo bind错误
		// todo 找不到user
		// todo publish返回错误
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := tc.mock(ctrl)

			r := gin.Default()
			r.Use(func(ctx *gin.Context) {
				ctx.Set("userClaims", &user.ClaimsUser{
					Id:   userId,
					Name: "test",
				})
			})

			h := NewArticleHandler(svc)
			h.RegisterRoutes(r)

			req, err := http.NewRequest(http.MethodPost,
				"/articles/publish", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var got Result
			err = json.NewDecoder(w.Body).Decode(&got)
			assert.NoError(t, err)
			require.Equal(t, tc.wantBody, got)
		})
	}
}
