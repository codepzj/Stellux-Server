package dao

import (
	"context"
	"errors"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DocumentContent struct {
	mongox.Model `bson:",inline"`
	DocumentId   bson.ObjectID `bson:"documentId"`   // 根文档Id
	Title        string        `bson:"title"`        // 文档标题
	Content      string        `bson:"content"`      // 文档内容
	Description  string        `bson:"description"`  // 文档描述
	Version      string        `bson:"version"`      // 文档版本
	Alias        string        `bson:"alias"`        // 文档别名
	ParentId     bson.ObjectID `bson:"parentId"`     // 父级Id
	IsDir        bool          `bson:"isDir"`        // 是否是目录
	Sort         int           `bson:"sort"`         // 排序
	LikeCount    int           `bson:"likeCount"`    // 点赞数
	DislikeCount int           `bson:"dislikeCount"` // 反对数
	CommentCount int           `bson:"commentCount"` // 评论数
}

var _ IDocumentContentDao = (*DocumentContentDao)(nil)

type IDocumentContentDao interface {
	CreateDocumentContent(ctx context.Context, doc DocumentContent) (bson.ObjectID, error)
	FindDocumentContentById(ctx context.Context, id bson.ObjectID) (DocumentContent, error)
	DeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error
	FindDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]DocumentContent, error)
	FindDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]DocumentContent, error)
	UpdateDocumentContentById(ctx context.Context, id bson.ObjectID, doc DocumentContent) error
}

type DocumentContentDao struct {
	coll *mongox.Collection[DocumentContent]
}

func NewDocumentContentDao(db *mongox.Database) *DocumentContentDao {
	return &DocumentContentDao{coll: mongox.NewCollection[DocumentContent](db, "document_content")}
}

// CreateDocumentContent 创建文档内容
func (d *DocumentContentDao) CreateDocumentContent(ctx context.Context, doc DocumentContent) (bson.ObjectID, error) {
	result, err := d.coll.Creator().InsertOne(ctx, &doc)
	if err != nil {
		return bson.ObjectID{}, err
	}
	return result.InsertedID.(bson.ObjectID), nil
}

// FindDocumentContentById 根据id查询文档内容
func (d *DocumentContentDao) FindDocumentContentById(ctx context.Context, id bson.ObjectID) (DocumentContent, error) {
	doc, err := d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)

	// 如果文档不存在，返回空
	if errors.Is(err, mongo.ErrNoDocuments) {
		return DocumentContent{}, errors.New("document_content中不存在文档")
	}
	if err != nil {
		return DocumentContent{}, err
	}
	return *doc, nil
}

// DeleteDocumentContentById 根据id删除文档内容
func (d *DocumentContentDao) DeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	result, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("文档不存在")
	}
	return err
}

// FindDocumentContentByParentId 根据父级Id查询文档内容
func (d *DocumentContentDao) FindDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]DocumentContent, error) {
	docs, err := d.coll.Finder().Filter(query.Eq("parentId", parentId)).Find(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = *doc
	}
	return results, nil
}

// FindDocumentContentByDocumentId 根据文档Id查询文档内容
func (d *DocumentContentDao) FindDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]DocumentContent, error) {
	docs, err := d.coll.Finder().Filter(query.Eq("documentId", documentId)).Find(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = *doc
	}
	return results, nil
}

// UpdateDocumentContentById 根据id更新文档内容
func (d *DocumentContentDao) UpdateDocumentContentById(ctx context.Context, id bson.ObjectID, doc DocumentContent) error {
	result, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.SetFields(bson.D{
		{Key: "documentId", Value: doc.DocumentId},
		{Key: "title", Value: doc.Title},
		{Key: "content", Value: doc.Content},
		{Key: "description", Value: doc.Description},
		{Key: "version", Value: doc.Version},
		{Key: "alias", Value: doc.Alias},
		{Key: "parentId", Value: doc.ParentId},
		{Key: "isDir", Value: doc.IsDir},
		{Key: "sort", Value: doc.Sort},
	})).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("文档不存在")
	}
	return err
}
