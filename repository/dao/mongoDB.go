package dao

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/mongo"
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
	return nil
}

func (m *MongoDB) SyncStatus(ctx context.Context, art ArticleAuthor) error {
	return nil
}
