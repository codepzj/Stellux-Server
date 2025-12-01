package infra

import (
	"context"
	"fmt"

	"github.com/codepzj/Stellux-Server/conf"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewMongoDB(cfg *conf.Config) *mongo.Database {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%d/?authSource=admin", cfg.MongoDB.MongoInitdbRootUsername, cfg.MongoDB.MongoInitdbRootPassword, cfg.MongoDB.Host, cfg.MongoDB.Port)
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(url))
	if err != nil {
		panic(errors.Wrap(err, "数据库连接失败"))
	}
	if err := mongoClient.Ping(context.Background(), readpref.Primary()); err != nil {
		panic(errors.Wrap(err, "数据库无法ping通"))
	}
	return mongoClient.Database(cfg.MongoDB.MongoInitdbDatabase)
}
