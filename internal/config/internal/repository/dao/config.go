package dao

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Config struct {
	ID        bson.ObjectID          `bson:"_id,omitempty"`
	CreatedAt time.Time              `bson:"created_at"`
	UpdatedAt time.Time              `bson:"updated_at"`
	Type      string                 `bson:"type"`
	Content   map[string]interface{} `bson:"content"`
}

type ConfigUpdate struct {
	Type      string                 `bson:"type"`
	Content   map[string]interface{} `bson:"content"`
	UpdatedAt time.Time              `bson:"updated_at"`
}

type IConfigDao interface {
	Create(ctx context.Context, config *Config) error
	Update(ctx context.Context, id bson.ObjectID, config *ConfigUpdate) error
	GetByID(ctx context.Context, id bson.ObjectID) (*Config, error)
	GetByType(ctx context.Context, configType string) (*Config, error)
	List(ctx context.Context) ([]*Config, error)
	Delete(ctx context.Context, id bson.ObjectID) error
}

var _ IConfigDao = (*ConfigDao)(nil)

func NewConfigDao(db *mongo.Database) *ConfigDao {
	return &ConfigDao{coll: db.Collection("config")}
}

type ConfigDao struct {
	coll *mongo.Collection
}

// Create 创建网站配置
func (d *ConfigDao) Create(ctx context.Context, config *Config) error {
	config.ID = bson.NewObjectID()
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()
	_, err := d.coll.InsertOne(ctx, config)
	return errors.Wrap(err, "failed to create config")
}

// Update 更新网站配置
func (d *ConfigDao) Update(ctx context.Context, id bson.ObjectID, config *ConfigUpdate) error {
	config.UpdatedAt = time.Now()
	update := bson.M{"$set": config}
	_, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	return errors.Wrap(err, "failed to update config")
}

// GetByID 根据ID获取网站配置
func (d *ConfigDao) GetByID(ctx context.Context, id bson.ObjectID) (*Config, error) {
	var config Config
	err := d.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config by id")
	}
	return &config, nil
}

// GetByType 根据类型获取网站配置
func (d *ConfigDao) GetByType(ctx context.Context, configType string) (*Config, error) {
	var config Config
	err := d.coll.FindOne(ctx, bson.M{"type": configType}).Decode(&config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config by type")
	}
	return &config, nil
}

// List 获取所有网站配置
func (d *ConfigDao) List(ctx context.Context) ([]*Config, error) {
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cursor, err := d.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list configs")
	}
	defer cursor.Close(ctx)

	var configs []*Config
	if err = cursor.All(ctx, &configs); err != nil {
		return nil, errors.Wrap(err, "failed to decode configs")
	}
	return configs, nil
}

// Delete 删除网站配置
func (d *ConfigDao) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := d.coll.DeleteOne(ctx, bson.M{"_id": id})
	return errors.Wrap(err, "failed to delete config")
}
