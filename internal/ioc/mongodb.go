package ioc

import (
	"context"
	"fmt"

	"github.com/chenmingyong0423/go-mongox/v2"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewMongoDB() *mongox.Database {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%d/?authSource=admin", viper.GetString("mongodb.MONGO_INITDB_ROOT_USERNAME"), viper.GetString("mongodb.MONGO_INITDB_ROOT_PASSWORD"), viper.GetString("mongodb.HOST"), viper.GetInt("mongodb.PORT"))
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(url))
	if err != nil {
		panic(errors.Wrap(err, "数据库连接失败"))
	}
	if err := mongoClient.Ping(context.Background(), readpref.Primary()); err != nil {
		panic(errors.Wrap(err, "数据库无法ping通"))
	}
	return mongox.NewClient(mongoClient, &mongox.Config{}).NewDatabase(viper.GetString("mongodb.MONGO_INITDB_DATABASE"))
}
