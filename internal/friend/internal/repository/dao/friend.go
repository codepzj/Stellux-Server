package dao

import (
	"context"
	"errors"

	"github.com/chenmingyong0423/go-mongox/v2"
)

type Friend struct {
	mongox.Model `bson:",inline"`
	Name        string
	Description string
	SiteUrl     string
	AvatarUrl   string
	WebsiteType string
	IsActive    bool
}

type IFriendDao interface {
	Create(ctx context.Context, friend *Friend) error
	FindAll(ctx context.Context) ([]*Friend, error)
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
		return errors.New("插入朋友")
	}
	return nil
}

func (d *FriendDao) FindAll(ctx context.Context) ([]*Friend, error) {
	return d.coll.Finder().Find(ctx)
}