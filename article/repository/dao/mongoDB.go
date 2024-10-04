//go:build mongo

package dao

import (
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDB struct {
	db         *mongo.Database
	coll       *mongo.Collection
	readerColl *mongo.Collection
	node       *snowflake.Node
}

func (m *MongoDB) GetListByAuthor(ctx context.Context, uid int64, limit int, offset int) ([]ArticleAuthor, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDB) GetByID(ctx context.Context, id int64) (ArticleAuthor, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDB) GetPublishedByID(ctx context.Context, id int64) (ArticleReader, error) {
	//TODO implement me
	panic("implement me")
}

func NewMongoDB(db *mongo.Database, node *snowflake.Node) ArticleDao {
	return &MongoDB{
		coll:       db.Collection("internal"),
		readerColl: db.Collection("reader"),
		node:       node,
	}
}

func (m *MongoDB) Insert(ctx context.Context, art ArticleAuthor) (int64, error) {
	id := m.node.Generate().Int64()
	art.ID = id
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now

	_, err := m.coll.InsertOne(ctx, art)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *MongoDB) UpdateByID(ctx context.Context, art ArticleAuthor) error {
	filter := bson.M{"id": art.ID, "authorid": art.Authorid}
	u := bson.D{bson.E{Key: "$set", Value: bson.M{
		"title":   art.Title,
		"content": art.Content,
		"status":  art.Status,
		"utime":   time.Now().UnixMilli(),
	}}}

	one, err := m.coll.UpdateOne(ctx, filter, u)

	if err != nil {
		return err
	}

	if one.ModifiedCount == 0 {
		return fmt.Errorf("系统错误")
	}

	return nil
}

func (m *MongoDB) SyncStatus(ctx context.Context, art ArticleAuthor) error {
	var err error
	var id int64

	if art.ID > 0 {
		err = m.UpdateByID(ctx, art)
	} else {
		id, err = m.Insert(ctx, art) //new
	}

	if err != nil {
		return err
	}

	art.ID = id
	now := time.Now().UnixMilli()
	art.Utime = now

	filter := bson.M{"id": art.ID}
	update := bson.E{Key: "$set", Value: art}
	upsert := bson.E{Key: "$setOnInsert", Value: bson.D{
		bson.E{Key: "ctime", Value: now}}}
	_, err = m.readerColl.UpdateOne(ctx, filter,
		bson.D{update, upsert},
		options.Update().SetUpsert(true))

	return err
}
