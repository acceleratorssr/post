package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"net/http/httptest"
	"post/domain"
	"post/ioc"
	"post/repository"
	"post/repository/dao"
	"post/service"
	"post/user"
	"post/web"
	"testing"
	"time"
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
	server     *gin.Engine
	db         *mongo.Database
	coll       *mongo.Collection
	readerColl *mongo.Collection
	node       *snowflake.Node
}

const (
	userId  = 99
	timeNow = 99
)

// 在测试执行前，初始化需要的数据
// setupTestSuite, ok := suite.(SetupTestSuite);
func (s *ArticleTestSuite) SetupTest() {
	s.server = gin.Default()
	// 此处手动设置jwt的信息保存到ctx里

	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("userClaims", &user.ClaimsUser{
			Id:   userId,
			Name: "test",
		})
	})

	s.node = dao.NewSnowflakeNode()
	s.db = ioc.InitMongoDB()
	//defer func() {
	//	if err = client.Disconnect(ctx); err != nil {
	//		panic(err)
	//	}
	//}()
	s.coll = s.db.Collection("post")
	s.readerColl = s.db.Collection("reader")
	d := dao.NewMongoDB(s.db, s.node)

	artHandler := web.NewArticleHandler(service.NewArticleService(
		repository.NewArticleAuthorRepository(d), repository.NewArticleReaderRepository(d)))

	artHandler.RegisterRoutes(s.server)

	//web.NewArticleHandler(service.NewArticleService(repository.NewArticleRepository(1.NewGORMArticleDao(dbInit())))).RegisterRoutes(s.server)
}

func (s *ArticleTestSuite) TestArticle() {
	t := s.T()
	snowID := s.node.Generate().Int64()
	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		article Article // 期望输入

		wantCode int
		wantData Result[int64] // 期望http响应带上帖子的id
	}{
		{
			name: "T1_新建草稿保存成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art dao.ArticleAuthor
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				err := s.coll.FindOne(ctx, bson.M{"title": "测试标题"}).Decode(&art)

				s.Require().NoError(err)
				s.Require().True(art.Ctime > 0)
				s.Require().True(art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				if art.Id > 1 {
					art.Id = 0
				} else {
					s.T().Error("雪花id不正常")
				}

				s.Require().Equal(dao.ArticleAuthor{
					Id:       0,
					Title:    "测试标题",
					Content:  "测试内容",
					Authorid: userId,
					Ctime:    0,
					Utime:    0,
					Status:   domain.TypeSaved.ToUint8(),
				}, art)

				// 清空所有数据
				fmt.Println("Drop all data")
				err = s.coll.Drop(ctx)
				if err != nil {
					log.Fatalf("Failed to drop collection 'article_authors': %v\n", err)
				}

				err = s.readerColl.Drop(ctx)
				if err != nil {
					log.Fatalf("Failed to drop collection 'article_readers': %v\n", err)
				}
			},
			article: Article{
				Title:   "测试标题",
				Content: "测试内容",
			},
			wantCode: http.StatusOK,
			wantData: Result[int64]{
				Code: 200,
				Data: 0,
				Msg:  "success",
			},
		},
		{
			name: "T2_修改草稿保存成功",
			before: func(t *testing.T) {
				_, err := s.coll.InsertOne(context.Background(), &dao.ArticleAuthor{
					Id:       snowID,
					Title:    "测试标题2",
					Content:  "测试内容2",
					Authorid: userId,
					Ctime:    timeNow,
					Utime:    timeNow,
					Status:   domain.TypeSaved.ToUint8(),
				})
				s.Require().NoError(err)
			},
			after: func(t *testing.T) {
				var art dao.ArticleAuthor
				err := s.coll.FindOne(context.Background(), bson.M{"id": snowID}).Decode(&art)
				s.Require().NoError(err)
				s.Require().True(art.Utime > 0)
				art.Utime = 0

				s.Require().Equal(dao.ArticleAuthor{
					Id:       snowID,
					Title:    "测试标题",
					Content:  "测试内容",
					Authorid: userId,
					Ctime:    timeNow,
					Utime:    0,
					Status:   domain.TypeSaved.ToUint8(),
				}, art)

				ctx := context.TODO()

				// 清空所有数据
				fmt.Println("Drop all data")
				err = s.coll.Drop(ctx)
				if err != nil {
					log.Fatalf("Failed to drop collection 'article_authors': %v\n", err)
				}

				err = s.readerColl.Drop(ctx)
				if err != nil {
					log.Fatalf("Failed to drop collection 'article_readers': %v\n", err)
				}
			},
			article: Article{
				ID:      snowID,
				Title:   "测试标题",
				Content: "测试内容",
			},
			wantCode: http.StatusOK,
			wantData: Result[int64]{
				Code: 200,
				Data: 0,
				Msg:  "success",
			},
		},
		{
			// 更新时带上author_id，更新失败即不是作者
			name: "T3_修改其他用户的草稿",
			before: func(t *testing.T) {
				_, err := s.coll.InsertOne(context.Background(), &dao.ArticleAuthor{
					Id:       snowID,
					Title:    "测试标题3",
					Content:  "测试内容3",
					Authorid: 965,
					Ctime:    timeNow,
					Utime:    timeNow,
					Status:   domain.TypeSaved.ToUint8(),
				})
				s.Require().NoError(err)
			},
			after: func(t *testing.T) {
				var art dao.ArticleAuthor
				err := s.coll.FindOne(context.Background(), bson.M{"id": snowID}).Decode(&art)
				s.Require().NoError(err)

				s.Require().Equal(dao.ArticleAuthor{
					Id:       snowID,
					Title:    "测试标题3",
					Content:  "测试内容3",
					Authorid: 965,
					Ctime:    timeNow,
					Utime:    timeNow,
					Status:   domain.TypeSaved.ToUint8(),
				}, art)

				ctx := context.TODO()

				// 清空所有数据
				fmt.Println("Drop all data")
				err = s.coll.Drop(ctx)
				if err != nil {
					log.Fatalf("Failed to drop collection 'article_authors': %v\n", err)
				}

				err = s.readerColl.Drop(ctx)
				if err != nil {
					log.Fatalf("Failed to drop collection 'article_readers': %v\n", err)
				}
			},
			article: Article{
				ID:      snowID,
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
			req.Header.Set("Content- Type", "application/json")

			//// 进行集成测试或端到端测试，这样可以确保测试覆盖实际的网络请求和响应处理
			//httpResp, err := http.DefaultClient.Do(req)

			// 进行单元测试，这样可以更快地测试处理程序的逻辑，并且便于控制测试环境和模拟各种请求场景
			resp := httptest.NewRecorder()
			s.server.ServeHTTP(resp, req) //为了获取resp

			//s.Require().Equal(tc.wantCode, resp.Code)
			var data Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&data)
			if data.Data > 1 {
				data.Data = 0
			}

			s.Require().NoError(err)

			fmt.Println(data)
			s.Require().Equal(tc.wantData, data)

			tc.after(t)
		})
	}
}

// 钩子函数
func TestArticle(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite))
}
