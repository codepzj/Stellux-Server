package dao

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Comment struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty"`
	Path      string        `bson:"path"`
	Content   string        `bson:"content"`
	RootId    bson.ObjectID `bson:"root_id"`   // 根评论Id
	ParentId  bson.ObjectID `bson:"parent_id"` // 父级评论Id
	Nickname  string        `bson:"nickname"`  // 昵称
	Avatar    string        `bson:"avatar"`    // 头像
	Email     string        `bson:"email"`     // 邮箱
	SiteUrl   string        `bson:"site_url"`  // 链接
	IsAdmin   bool          `bson:"is_admin"`  // 是否管理员
}

type ICommentDao interface {
	Create(ctx context.Context, comment *Comment) error
	GetList(ctx context.Context, filter bson.D) ([]*Comment, error)
	Update(ctx context.Context, id bson.ObjectID, comment *Comment) error
	Delete(ctx context.Context, id bson.ObjectID) error
}

var _ ICommentDao = (*CommentDao)(nil)

func NewCommentDao(db *mongo.Database) *CommentDao {
	return &CommentDao{coll: db.Collection("comment")}
}

type CommentDao struct {
	coll *mongo.Collection
}

// Create 创建评论
func (d *CommentDao) Create(ctx context.Context, comment *Comment) error {
	comment.ID = bson.NewObjectID()
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	res, err := d.coll.InsertOne(ctx, comment)
	if err != nil {
		return err
	}
	if res.InsertedID == nil {
		return errors.New("创建评论失败")
	}
	return nil
}

// GetList 获取评论列表
func (d *CommentDao) GetList(ctx context.Context, filter bson.D) ([]*Comment, error) {
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cursor, err := d.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

// Update 更新评论
func (d *CommentDao) Update(ctx context.Context, id bson.ObjectID, comment *Comment) error {
	update := bson.M{
		"$set": bson.M{
			"content":    comment.Content,
			"nickname":   comment.Nickname,
			"avatar":     comment.Avatar,
			"email":      comment.Email,
			"site_url":   comment.SiteUrl,
			"updated_at": time.Now(),
		},
	}
	res, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("更新评论失败")
	}
	return nil
}

// Delete 删除评论
func (d *CommentDao) Delete(ctx context.Context, id bson.ObjectID) error {
	res, err := d.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("删除评论失败")
	}
	return nil
}
