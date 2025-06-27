package dao

import (
	"context"
	"errors"

	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/chenmingyong0423/go-mongox/v2"
)

type Friend struct {
	mongox.Model `bson:",inline"`
	Name         string
	Description  string
	SiteUrl      string `bson:"site_url"`
	AvatarUrl    string `bson:"avatar_url"`
	WebsiteType  string `bson:"website_type"`
	IsActive     bool   `bson:"is_active"`
}

type IFriendDao interface {
	Create(ctx context.Context, friend *Friend) error
	FindAll(ctx context.Context) ([]*Friend, error)
	Update(ctx context.Context, id bson.ObjectID, friend *Friend) error
	Delete(ctx context.Context, id bson.ObjectID) error
}

var _ IFriendDao = (*FriendDao)(nil)

func NewFriendDao(db *mongox.Database) *FriendDao {
	return &FriendDao{coll: mongox.NewCollection[Friend](db, "friend")}
}

type FriendDao struct {
	coll *mongox.Collection[Friend]
}

func (d *FriendDao) Create(ctx context.Context, friend *Friend) error {
	insertResult, err := d.coll.Creator().InsertOne(ctx, friend)
	if err != nil {
		return err
	}
	if insertResult.InsertedID == nil {
		return errors.New("插入朋友记录失败")
	}
	return nil
}

func (d *FriendDao) FindAll(ctx context.Context) ([]*Friend, error) {
	return d.coll.Finder().Find(ctx)
}

func (d *FriendDao) Update(ctx context.Context, id bson.ObjectID, friend *Friend) error {
	updateResult, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.NewBuilder().Set("name", friend.Name).Set("description", friend.Description).Set("site_url", friend.SiteUrl).Set("avatar_url", friend.AvatarUrl).Set("website_type", friend.WebsiteType).Set("is_active", friend.IsActive).Build()).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if updateResult.MatchedCount == 0 {
		return errors.New("更新朋友记录失败")
	}
	return nil
}

func (d *FriendDao) Delete(ctx context.Context, id bson.ObjectID) error {
	deleteResult, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return errors.New("删除朋友记录失败")
	}
	return nil
}