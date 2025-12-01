package dao

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Friend struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	CreatedAt   time.Time     `bson:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"`
	DeletedAt   *time.Time    `bson:"deleted_at,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	SiteUrl     string        `bson:"site_url"`
	AvatarUrl   string        `bson:"avatar_url"`
	WebsiteType int           `bson:"website_type"`
	IsActive    bool          `bson:"is_active"`
}

type IFriendDao interface {
	Create(ctx context.Context, friend *Friend) error
	FindAll(ctx context.Context) ([]*Friend, error)
	FindAllActive(ctx context.Context) ([]*Friend, error)
	ExistsBySiteUrl(ctx context.Context, siteUrl string) (bool, error)
	ExistsBySiteUrlExceptID(ctx context.Context, siteUrl string, excludeID bson.ObjectID) (bool, error)
	Update(ctx context.Context, id bson.ObjectID, friend *Friend) error
	Delete(ctx context.Context, id bson.ObjectID) error
}

var _ IFriendDao = (*FriendDao)(nil)

func NewFriendDao(db *mongo.Database) *FriendDao {
	return &FriendDao{coll: db.Collection("friend")}
}

type FriendDao struct {
	coll *mongo.Collection
}

func (d *FriendDao) Create(ctx context.Context, friend *Friend) error {
	friend.ID = bson.NewObjectID()
	friend.CreatedAt = time.Now()
	friend.UpdatedAt = time.Now()
	insertResult, err := d.coll.InsertOne(ctx, friend)
	if err != nil {
		return err
	}
	if insertResult.InsertedID == nil {
		return errors.New("插入朋友记录失败")
	}
	return nil
}

func (d *FriendDao) FindAll(ctx context.Context) ([]*Friend, error) {
	cursor, err := d.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var friends []*Friend
	if err = cursor.All(ctx, &friends); err != nil {
		return nil, err
	}
	return friends, nil
}

// FindAllActive 查询所有活跃的友链
func (d *FriendDao) FindAllActive(ctx context.Context) ([]*Friend, error) {
	cursor, err := d.coll.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var friends []*Friend
	if err = cursor.All(ctx, &friends); err != nil {
		return nil, err
	}
	return friends, nil
}

// ExistsBySiteUrl 判断指定 site_url 是否已存在
func (d *FriendDao) ExistsBySiteUrl(ctx context.Context, siteUrl string) (bool, error) {
	count, err := d.coll.CountDocuments(ctx, bson.M{"site_url": siteUrl})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsBySiteUrlExceptID 判断指定 site_url 是否已存在（排除指定ID）
func (d *FriendDao) ExistsBySiteUrlExceptID(ctx context.Context, siteUrl string, excludeID bson.ObjectID) (bool, error) {
	filter := bson.M{
		"site_url": siteUrl,
		"_id":      bson.M{"$ne": excludeID},
	}
	count, err := d.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *FriendDao) Update(ctx context.Context, id bson.ObjectID, friend *Friend) error {
	update := bson.M{
		"$set": bson.M{
			"name":         friend.Name,
			"description":  friend.Description,
			"site_url":     friend.SiteUrl,
			"avatar_url":   friend.AvatarUrl,
			"website_type": friend.WebsiteType,
			"is_active":    friend.IsActive,
			"updated_at":   time.Now(),
		},
	}
	updateResult, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if updateResult.MatchedCount == 0 {
		return errors.New("更新朋友记录失败")
	}
	return nil
}

func (d *FriendDao) Delete(ctx context.Context, id bson.ObjectID) error {
	deleteResult, err := d.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return errors.New("删除朋友记录失败")
	}
	return nil
}
