package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"post/domain"
	"post/repository/dao"
	"testing"
)

type Article struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// 使用泛型，防止反序列化时出现问题
type Result[T any] struct {
	Code int    `json:"code"`
	Data T      `json:"data"`
	Msg  string `json:"msg"`
}

// ArticleTestSuite
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

const (
	userId  = 99
	timeNow = 99
)

// 在测试执行前，初始化需要的数据
// setupTestSuite, ok := suite.(SetupTestSuite);
func (s *ArticleTestSuite) SetupTest() {
	app := InitApp()
	s.server = app.server
	s.db = app.db
	// 此处手动设置jwt的信息保存到ctx里

	//s.server.Use(func(ctx *gin.Context) {
	//	ctx.Set("userClaims", &user.ClaimsUser{
	//		Id:   userId,
	//		Name: "test",
	//	})
	//})

	//app := InitApp()
	//
	//app.articleHandler.RegisterRoutes(s.server)

	//web.NewArticleHandler(service.NewArticleService(repository.NewArticleRepository(dao.NewGORMArticleDao(dbInit())))).RegisterRoutes(s.server)
}

func (s *ArticleTestSuite) TestArticle() {
	t := s.T()
	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		article  Article // 期望输入
		wantCode int
		wantData Result[int64] // 期望http响应带上帖子的id
	}{
		{
			name: "T1_新建草稿保存成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art dao.ArticleAuthor
				err := s.db.Where("id = ?", 1).First(&art).Error //预期获取id1
				s.Require().NoError(err)
				s.Require().True(art.Ctime > 0)
				s.Require().True(art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0

				s.Require().Equal(dao.ArticleAuthor{
					Id:       1,
					Title:    "测试标题",
					Content:  "测试内容",
					Authorid: userId,
					Ctime:    0,
					Utime:    0,
					Status:   domain.TypeSaved.ToUint8(),
				}, art)
			},
			article: Article{
				Title:   "测试标题",
				Content: "测试内容",
			},
			wantCode: http.StatusOK,
			wantData: Result[int64]{
				Code: 200,
				Data: 1,
				Msg:  "success",
			},
		},
		{
			name: "T2_修改草稿保存成功",
			before: func(t *testing.T) {
				err := s.db.Create(&dao.ArticleAuthor{
					Id:       2,
					Title:    "测试标题2",
					Content:  "测试内容2",
					Authorid: userId,
					Ctime:    timeNow, //时间可以自己手写一个固定的，不用now
					Utime:    timeNow,
					Status:   domain.TypeSaved.ToUint8(),
				}).Error
				s.Require().NoError(err)
			},
			after: func(t *testing.T) {
				var art dao.ArticleAuthor
				err := s.db.Where("id = ?", 2).First(&art).Error //预期获取id1
				s.Require().NoError(err)
				s.Require().True(art.Utime > 0)
				art.Utime = 0

				s.Require().Equal(dao.ArticleAuthor{
					Id:       2,
					Title:    "测试标题",
					Content:  "测试内容",
					Authorid: userId,
					Ctime:    timeNow,
					Utime:    0,
					Status:   domain.TypeSaved.ToUint8(),
				}, art)
			},
			article: Article{
				ID:      2,
				Title:   "测试标题",
				Content: "测试内容",
			},
			wantCode: http.StatusOK,
			wantData: Result[int64]{
				Code: 200,
				Data: 2,
				Msg:  "success",
			},
		},
		{
			// 更新时带上author_id，更新失败即不是作者
			name: "T2_修改其他用户的草稿",
			before: func(t *testing.T) {
				err := s.db.Create(&dao.ArticleAuthor{
					Id:       3,
					Title:    "测试标题3",
					Content:  "测试内容3",
					Authorid: 965,
					Ctime:    timeNow,
					Utime:    timeNow,
					Status:   domain.TypeSaved.ToUint8(),
				}).Error
				s.Require().NoError(err)
			},
			after: func(t *testing.T) {
				var art dao.ArticleAuthor
				err := s.db.Where("id = ?", 3).First(&art).Error //预期获取id1
				s.Require().NoError(err)

				s.Require().Equal(dao.ArticleAuthor{
					Id:       3,
					Title:    "测试标题3",
					Content:  "测试内容3",
					Authorid: 965,
					Ctime:    timeNow,
					Utime:    timeNow,
					Status:   domain.TypeSaved.ToUint8(),
				}, art)
			},
			article: Article{
				ID:      3,
				Title:   "测试标题",
				Content: "测试内容",
			},
			wantCode: http.StatusOK,
			wantData: Result[int64]{
				Code: 555,
				Data: 0,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)

			reqBody, err := json.Marshal(tc.article)
			s.Require().NoError(err)
			req, err := http.NewRequest(http.MethodPost,
				"/articles/save", bytes.NewBuffer(reqBody))
			s.Require().NoError(err)
			req.Header.Set("Content-Type", "application/json")

			//// 进行集成测试或端到端测试，这样可以确保测试覆盖实际的网络请求和响应处理
			//httpResp, err := http.DefaultClient.Do(req)

			// 进行单元测试，这样可以更快地测试处理程序的逻辑，并且便于控制测试环境和模拟各种请求场景
			resp := httptest.NewRecorder()
			s.server.ServeHTTP(resp, req) //为了获取resp

			//s.Require().Equal(tc.wantCode, resp.Code)
			var data Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&data)
			s.Require().NoError(err)
			s.Require().Equal(tc.wantData, data)

			tc.after(t)
		})
	}
}

// 钩子函数
func (s *ArticleTestSuite) TearDownTest() {
	// 清空所有数据，主键重置回0
	s.db.Exec("truncate table article_authors")
	s.db.Exec("truncate table article_readers")
}

func TestArticle(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite))
}

func dbInit() *gorm.DB {
	dsn := "root:20031214pzw!@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
