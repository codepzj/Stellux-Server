package dao

import (
	"context"
	"errors"

	"github.com/chenmingyong0423/go-mongox/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type DocumentContent struct {
	mongox.Model `bson:",inline"`
	DocumentId   bson.ObjectID `bson:"documentId"`
	Title        string        `bson:"title"`   // 文档标题
	Content      string        `bson:"content"` // 文档内容
	Version      int           `bson:"version"` // 文档版本
	Alias        string        `bson:"alias"`   // 文档别名
	ParentID     bson.ObjectID `bson:"parentId"` // 父级ID
	IsDir        bool          `bson:"isDir"`  // 是否是目录
	LikeCount    int           `bson:"likeCount"` // 点赞数
	DislikeCount int           `bson:"dislikeCount"` // 反对数
	CommentCount int           `bson:"commentCount"` // 评论数
}

var _ IDocumentContentDao = (*DocumentContentDao)(nil)

type IDocumentContentDao interface {
	CreateDocumentContent(ctx context.Context, doc *DocumentContent) error
	FindDocumentContentByDocumentID(ctx context.Context, documentId bson.ObjectID) ([]*DocumentContent, error)
	
}

type DocumentContentDao struct {
	coll *mongox.Collection[DocumentContent]
}

func NewDocumentContentDao(db *mongox.Database) *DocumentContentDao {
	return &DocumentContentDao{coll: mongox.NewCollection[DocumentContent](db, "document_content")}
}

// CreateDocumentContent 创建文档内容
func (d *DocumentContentDao) CreateDocumentContent(ctx context.Context, doc *DocumentContent) error {
	_, err := d.coll.Creator().InsertOne(ctx, doc)
	return err
}

// FindDocumentContentByDocumentID 根据文档ID查询文档内容
func (d *DocumentContentDao) FindDocumentContentByDocumentID(ctx context.Context, documentId bson.ObjectID) ([]*DocumentContent, error) {
	docs, err := d.coll.Finder().Filter(bson.M{"documentId": documentId}).Find(ctx)

	// 如果文档不存在，返回空数组
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return docs, nil
}
