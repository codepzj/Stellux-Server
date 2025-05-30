package dao

import (
	"context"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Setting struct {
	Key   string `bson:"key"`
	Value any    `bson:"value"`
}

type ISettingDao interface {
	Upsert(ctx context.Context, setting *Setting) error
	GetSetting(ctx context.Context, key string) (*Setting, error)
}

var _ ISettingDao = (*SettingDao)(nil)

func NewSettingDao(db *mongox.Database) *SettingDao {
	return &SettingDao{coll: mongox.NewCollection[Setting](db, "setting")}
}

type SettingDao struct {
	coll *mongox.Collection[Setting]
}

func (s *SettingDao) Upsert(ctx context.Context, setting *Setting) error {
	updateResult, err := s.coll.Updater().Filter(query.Eq("key", setting.Key)).Updates(update.NewBuilder().Set("key", setting.Key).Set("value", setting.Value).Build()).Upsert(ctx)
	if err != nil {
		return err
	}
	if updateResult.ModifiedCount == 0 && updateResult.UpsertedID == nil {
		return errors.New("设置新增或更新失败")
	}
	return nil
}

func (s *SettingDao) GetSetting(ctx context.Context, key string) (*Setting, error) {
	setting, err := s.coll.Finder().Filter(query.Eq("key", key)).FindOne(ctx)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &Setting{
			Key:   key,
			Value: nil,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &Setting{
		Key:   setting.Key,
		Value: setting.Value,
	}, nil
}
