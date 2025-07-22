package dao

import (
	"context"
	"time"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Document struct {
	mongox.Model `bson:",inline"`
	Title        string `bson:"title"`
	Description  string `bson:"description"`
	Thumbnail    string `bson:"thumbnail"`
	Alias        string `bson:"alias"`
	Sort         int    `bson:"sort"`
	IsPublic     bool   `bson:"isPublic"`
	IsDeleted    bool   `bson:"isDeleted"`
}

type IDocumentDao interface {
	CreateDocument(ctx context.Context, doc Document) (bson.ObjectID, error)
	FindDocumentById(ctx context.Context, id bson.ObjectID) (Document, error)
	UpdateDocumentById(ctx context.Context, id bson.ObjectID, doc Document) error
	DeleteDocumentById(ctx context.Context, id bson.ObjectID) error
	SoftDeleteDocumentById(ctx context.Context, id bson.ObjectID) error
	RestoreDocumentById(ctx context.Context, id bson.ObjectID) error
	FindDocumentByAlias(ctx context.Context, alias string, filter bson.D) (Document, error)
	GetDocumentList(ctx context.Context, page *Page) ([]Document, int64, error)
	GetPublicDocumentList(ctx context.Context, page *Page) ([]Document, int64, error)
	GetAllPublicDocuments(ctx context.Context) ([]Document, error)
}

// Page 分页查询参数
type Page struct {
	PageNo   int64 `bson:"pageNo"`   // 页码
	PageSize int64 `bson:"pageSize"` // 每页大小
}

var _ IDocumentDao = (*DocumentDao)(nil)

func NewDocumentDao(db *mongox.Database) *DocumentDao {
	return &DocumentDao{coll: mongox.NewCollection[Document](db, "document")}
}

type DocumentDao struct {
	coll *mongox.Collection[Document]
}

func (d *DocumentDao) CreateDocument(ctx context.Context, doc Document) (bson.ObjectID, error) {
	result, err := d.coll.Creator().InsertOne(ctx, &doc)
	if err != nil {
		return bson.ObjectID{}, err
	}
	return result.InsertedID.(bson.ObjectID), nil
}

func (d *DocumentDao) FindDocumentById(ctx context.Context, id bson.ObjectID) (Document, error) {
	document, err := d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
	if err != nil {
		return Document{}, err
	}
	return *document, nil
}

func (d *DocumentDao) UpdateDocumentById(ctx context.Context, id bson.ObjectID, doc Document) error {
	result, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.SetFields(bson.D{
		{Key: "title", Value: doc.Title},
		{Key: "description", Value: doc.Description},
		{Key: "thumbnail", Value: doc.Thumbnail},
		{Key: "alias", Value: doc.Alias},
		{Key: "sort", Value: doc.Sort},
		{Key: "isPublic", Value: doc.IsPublic},
	})).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("文档不存在, 更新失败")
	}
	return nil
}

func (d *DocumentDao) DeleteDocumentById(ctx context.Context, id bson.ObjectID) error {
	result, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("文档不存在, 删除失败")
	}
	return nil
}

func (d *DocumentDao) SoftDeleteDocumentById(ctx context.Context, id bson.ObjectID) error {
	result, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.SetFields(bson.D{
		{Key: "deleted_at", Value: time.Now()},
		{Key: "isDeleted", Value: true},
	})).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("文档不存在, 软删除失败")
	}
	return nil
}

func (d *DocumentDao) RestoreDocumentById(ctx context.Context, id bson.ObjectID) error {
	result, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.SetFields(bson.D{
		{Key: "deleted_at", Value: nil},
		{Key: "isDeleted", Value: false},
	})).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("文档不存在, 恢复失败")
	}
	return nil
}

func (d *DocumentDao) FindDocumentByAlias(ctx context.Context, alias string, filter bson.D) (Document, error) {
	document, err := d.coll.Finder().Filter(query.Eq("alias", alias)).Filter(filter).FindOne(ctx)
	if err != nil {
		return Document{}, err
	}
	if document == nil {
		return Document{}, errors.New("文档不存在")
	}
	return *document, nil
}

func (d *DocumentDao) GetDocumentList(ctx context.Context, page *Page) ([]Document, int64, error) {
	skip := (page.PageNo - 1) * page.PageSize

	// 获取总数
	count, err := d.coll.Finder().Filter(query.Eq("isDeleted", false)).Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 获取列表
	documents, err := d.coll.Finder().
		Filter(query.Eq("isDeleted", false)).
		Sort(bson.D{{Key: "sort", Value: 1}, {Key: "createdAt", Value: -1}}).
		Skip(skip).
		Limit(page.PageSize).
		Find(ctx)
	if err != nil {
		return nil, 0, err
	}

	results := make([]Document, len(documents))
	for i, doc := range documents {
		results[i] = *doc
	}

	return results, count, nil
}

func (d *DocumentDao) GetPublicDocumentList(ctx context.Context, page *Page) ([]Document, int64, error) {
	skip := (page.PageNo - 1) * page.PageSize

	// 获取总数
	count, err := d.coll.Finder().Filter(query.And(query.Eq("isPublic", true), query.Eq("isDeleted", false))).Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 获取列表
	documents, err := d.coll.Finder().
		Filter(query.And(query.Eq("isPublic", true), query.Eq("isDeleted", false))).
		Sort(bson.D{{Key: "sort", Value: 1}, {Key: "createdAt", Value: -1}}).
		Skip(skip).
		Limit(page.PageSize).
		Find(ctx)
	if err != nil {
		return nil, 0, err
	}

	results := make([]Document, len(documents))
	for i, doc := range documents {
		results[i] = *doc
	}

	return results, count, nil
}

// GetAllPublicDocuments 获取所有公开文档
func (d *DocumentDao) GetAllPublicDocuments(ctx context.Context) ([]Document, error) {
	documents, err := d.coll.Finder().
		Filter(query.And(query.Eq("isPublic", true), query.Eq("isDeleted", false))).
		Sort(bson.D{{Key: "sort", Value: 1}, {Key: "createdAt", Value: -1}}).
		Find(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]Document, len(documents))
	for i, doc := range documents {
		results[i] = *doc
	}

	return results, nil
}
