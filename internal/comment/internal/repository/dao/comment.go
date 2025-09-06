package dao

import (
	"context"
	"errors"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Comment struct {
	mongox.Model `bson:",inline"`
	Path         string        `bson:"path"`
	Content      string        `bson:"content"`
	RootId       bson.ObjectID `bson:"root_id"`    // 根评论Id
	ParentId     bson.ObjectID `bson:"parent_id"`  // 父级评论Id
	Nickname     string        `bson:"nickname"`   // 昵称
	Avatar       string        `bson:"avatar"`     // 头像
	Email        string        `bson:"email"`      // 邮箱
	SiteUrl      string        `bson:"site_url"`   // 链接
	IsAdmin      bool          `bson:"is_admin"`   // 是否管理员
}

type ICommentDao interface {
	Create(ctx context.Context, comment *Comment) error
	GetList(ctx context.Context, filter bson.D) ([]*Comment, error)
	Update(ctx context.Context, id bson.ObjectID, comment *Comment) error
	Delete(ctx context.Context, id bson.ObjectID) error
}

var _ ICommentDao = (*CommentDao)(nil)

func NewCommentDao(db *mongox.Database) *CommentDao {
	return &CommentDao{coll: mongox.NewCollection[Comment](db, "comment")}
}

type CommentDao struct {
	coll *mongox.Collection[Comment]
}

// Create 创建评论
func (d *CommentDao) Create(ctx context.Context, comment *Comment) error {
	res, err := d.coll.Creator().InsertOne(ctx, comment)
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
	return d.coll.Finder().Filter(filter).Sort(bson.M{"created_at": -1}).Find(ctx)
}

// Update 更新评论
func (d *CommentDao) Update(ctx context.Context, id bson.ObjectID, comment *Comment) error {
	res, err := d.coll.Updater().Filter(query.Id(id)).Updates(map[string]any{
		"content":  comment.Content,
		"nickname": comment.Nickname,
		"avatar":   comment.Avatar,
		"email":    comment.Email,
		"site_url": comment.SiteUrl,
	}).UpdateOne(ctx)
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
	res, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("删除评论失败")
	}
	return nil
}
