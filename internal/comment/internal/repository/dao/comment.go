package dao

import (
	"context"
	"errors"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Comment struct {
	mongox.Model `bson:",inline"`
	Nickname     string `bson:"nickname"`
	Email        string `bson:"email"`
	Avatar       string `bson:"avatar"`
	Url          string `bson:"url"`
	Content      string `bson:"content"`
	PostId       bson.ObjectID `bson:"post_id"`
	CommentId    bson.ObjectID `bson:"comment_id"`
}

type ICommentDao interface {
	Create(ctx context.Context, comment *Comment) error
	Edit(ctx context.Context, comment bson.D) error
	Delete(ctx context.Context, id bson.ObjectID) error
	GetListByPostId(ctx context.Context, postId bson.ObjectID) ([]*Comment, error)
}

var _ ICommentDao = (*CommentDao)(nil)

func NewCommentDao(db *mongox.Database) *CommentDao {
	return &CommentDao{coll: mongox.NewCollection[Comment](db, "comment")}
}

type CommentDao struct {
	coll *mongox.Collection[Comment]
}

func (d *CommentDao) Create(ctx context.Context, comment *Comment) error {
	result, err := d.coll.Creator().InsertOne(ctx, comment)
	if err != nil {
		return err
	}
	if result.InsertedID == nil {
		return errors.New("insert failed")
	}
	return nil
}

func (d *CommentDao) Edit(ctx context.Context, comment bson.D) error {
	result, err := d.coll.Updater().Updates(update.SetFields(comment)).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("comment not found")
	}
	return err
}

func (d *CommentDao) Delete(ctx context.Context, id bson.ObjectID) error {
	result, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("comment not found")
	}
	return nil
}

func (d *CommentDao) GetListByPostId(ctx context.Context, postId bson.ObjectID) ([]*Comment, error) {
	comments, err := d.coll.Finder().Filter(bson.M{"post_id": postId}).Find(ctx)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
