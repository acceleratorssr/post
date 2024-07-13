package dao

import (
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type MongoDB struct {
	db         *mongo.Database
	coll       *mongo.Collection
	readerColl *mongo.Collection
	node       *snowflake.Node
}

func NewMongoDB(db *mongo.Database, node *snowflake.Node) ArticleDao {
	return &MongoDB{
		coll:       db.Collection("post"),
		readerColl: db.Collection("reader"),
		node:       node,
	}
}

func (m *MongoDB) Insert(ctx context.Context, art ArticleAuthor) (int64, error) {
	id := m.node.Generate().Int64()
	art.Id = id
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
	filter := bson.M{"id": art.Id, "authorid": art.Authorid}
	u := bson.D{bson.E{Key: "$set", Value: bson.M{
		"title":   art.Title,
		"content": art.Content,
		"status":  art.Status,
		"utime":   time.Now().UnixMilli(),
	}}}

	var existing ArticleAuthor
	err := m.coll.FindOne(ctx, filter).Decode(&existing)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No matching documents found")
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Found matching document: %+v\n", existing)
	}

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
	return nil
}
