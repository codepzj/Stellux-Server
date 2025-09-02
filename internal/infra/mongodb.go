package infra

import (
	"context"
	"fmt"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/codepzj/Stellux-Server/conf"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewMongoDB(cfg *conf.Config) *mongox.Database {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%d/?authSource=admin", cfg.MongoDB.MongoInitdbRootUsername, cfg.MongoDB.MongoInitdbRootPassword, cfg.MongoDB.Host, cfg.MongoDB.Port)
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(url))
	if err != nil {
		panic(errors.Wrap(err, "数据库连接失败"))
	}
	if err := mongoClient.Ping(context.Background(), readpref.Primary()); err != nil {
		panic(errors.Wrap(err, "数据库无法ping通"))
	}
	return mongox.NewClient(mongoClient, &mongox.Config{}).NewDatabase(cfg.MongoDB.MongoInitdbDatabase)
}
