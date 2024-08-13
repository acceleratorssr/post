package ioc

import (
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"post/interactive/repository/dao"
	"post/migrator/events"
	"post/migrator/events/fixer"
	"post/migrator/scheduler"
	"post/pkg/gin_ex"
	"post/pkg/gorm_ex/connpool"
)

func InitMigratorServer(base BaseDB, target TargetDB,
	pool *connpool.DoubleWritePool, producer events.InconsistentProducer) *gin_ex.Server {
	intrScheduler := scheduler.NewScheduler[dao.Like](base, target, pool, producer)

	engine := gin.Default()
	intrScheduler.RegisterRoutes(engine.Group("/migrator/like"))

	return &gin_ex.Server{
		Addr:   "9300",
		Engine: engine,
	}
}

func InitFixConsumer(base BaseDB, target TargetDB, client sarama.Client) *fixer.Consumer[dao.Like] {
	res := fixer.NewConsumer[dao.Like](client, base, target, "migrator_like")
	return res
}

func InitMigratorProducer(p sarama.SyncProducer) events.InconsistentProducer {
	return events.NewSaramaProducer(p, "migrator_like")
}

func InitDoubleWritePool(base BaseDB, target TargetDB) *connpool.DoubleWritePool {
	return connpool.NewDoubleWritePool(base.ConnPool, target.ConnPool)
}
