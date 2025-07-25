package dao

import (
	"context"
	"errors"
	"time"

	"fmt"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DocumentContent struct {
	mongox.Model `bson:",inline"`
	DocumentId   bson.ObjectID `bson:"document_id"` // 根文档Id
	Title        string        `bson:"title"`       // 文档标题
	Content      string        `bson:"content"`     // 文档内容
	Description  string        `bson:"description"` // 文档描述
	Alias        string        `bson:"alias"`       // 文档别名
	ParentId     bson.ObjectID `bson:"parent_id"`   // 父级Id
	IsDir        bool          `bson:"is_dir"`      // 是否是目录
	Sort         int           `bson:"sort"`        // 排序
	IsDeleted    bool          `bson:"is_deleted"`  // 是否删除
}

// Page 分页查询参数
type Page struct {
	PageNo   int64 `bson:"page_no"`   // 页码
	PageSize int64 `bson:"page_size"` // 每页大小
}

var _ IDocumentContentDao = (*DocumentContentDao)(nil)

type IDocumentContentDao interface {
	CreateDocumentContent(ctx context.Context, doc DocumentContent) (bson.ObjectID, error)
	FindDocumentContentById(ctx context.Context, id bson.ObjectID) (DocumentContent, error)
	DeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error
	SoftDeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error
	RestoreDocumentContentById(ctx context.Context, id bson.ObjectID) error
	FindDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]DocumentContent, error)
	FindDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]DocumentContent, error)
	UpdateDocumentContentById(ctx context.Context, id bson.ObjectID, doc DocumentContent) error
	GetDocumentContentList(ctx context.Context, page *Page) ([]DocumentContent, int64, error)
	GetPublicDocumentContentListByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]*DocumentContent, error)
	SearchDocumentContent(ctx context.Context, keyword string) ([]DocumentContent, error)
	SearchPublicDocumentContent(ctx context.Context, keyword string) ([]DocumentContent, error)
	FindPublicDocumentContentById(ctx context.Context, id bson.ObjectID) (DocumentContent, error)
	FindPublicDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]DocumentContent, error)
	FindPublicDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]DocumentContent, error)
	FindPublicDocumentContentByRootIdAndAlias(ctx context.Context, documentId bson.ObjectID, alias string) (DocumentContent, error)
	DeleteDocumentContentList(ctx context.Context, ids []string) error
	JudgeDocumentContentAliasUnique(ctx context.Context, alias string, documentId bson.ObjectID) (bool, error)
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

// SoftDeleteDocumentContentById 根据id软删除文档内容
func (d *DocumentContentDao) SoftDeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	result, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.SetFields(bson.D{
		{Key: "deleted_at", Value: time.Now()},
		{Key: "is_deleted", Value: true},
	})).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("文档不存在, 软删除失败")
	}
	return nil
}

// RestoreDocumentContentById 根据id恢复文档内容
func (d *DocumentContentDao) RestoreDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	result, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.SetFields(bson.D{
		{Key: "deleted_at", Value: nil},
		{Key: "is_deleted", Value: false},
	})).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("文档不存在, 恢复失败")
	}
	return nil
}

// FindDocumentContentByParentId 根据父级Id查询文档内容
func (d *DocumentContentDao) FindDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]DocumentContent, error) {
	docs, err := d.coll.Finder().Filter(query.Eq("parent_id", parentId)).Find(ctx)
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
	docs, err := d.coll.Finder().Filter(query.Eq("document_id", documentId)).Find(ctx)
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
		{Key: "document_id", Value: doc.DocumentId},
		{Key: "title", Value: doc.Title},
		{Key: "content", Value: doc.Content},
		{Key: "description", Value: doc.Description},
		{Key: "alias", Value: doc.Alias},
		{Key: "parent_id", Value: doc.ParentId},
		{Key: "is_dir", Value: doc.IsDir},
		{Key: "sort", Value: doc.Sort},
		{Key: "is_deleted", Value: doc.IsDeleted},
	})).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("文档不存在")
	}
	return err
}

// GetDocumentContentList 获取文档内容列表
func (d *DocumentContentDao) GetDocumentContentList(ctx context.Context, page *Page) ([]DocumentContent, int64, error) {
	skip := (page.PageNo - 1) * page.PageSize

	// 获取总数
	count, err := d.coll.Finder().Filter(query.Eq("is_deleted", false)).Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 获取列表
	documents, err := d.coll.Finder().
		Filter(query.Eq("is_deleted", false)).
		Sort(bson.D{{Key: "sort", Value: 1}, {Key: "created_at", Value: -1}}).
		Skip(skip).
		Limit(page.PageSize).
		Find(ctx)
	if err != nil {
		return nil, 0, err
	}

	results := make([]DocumentContent, len(documents))
	for i, doc := range documents {
		results[i] = *doc
	}

	return results, count, nil
}

// GetPublicDocumentContentList 获取公开文档内容列表
func (d *DocumentContentDao) GetPublicDocumentContentListByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]*DocumentContent, error) {
	// 获取列表
	documents, err := d.coll.Finder().
		Filter(query.And(query.Eq("document_id", documentId), query.Eq("is_deleted", false))).
		Find(ctx)
	if err != nil {
		return nil, err
	}

	return documents, nil
}

// SearchDocumentContent 搜索文档内容
func (d *DocumentContentDao) SearchDocumentContent(ctx context.Context, keyword string) ([]DocumentContent, error) {
	docs, err := d.coll.Finder().
		Filter(query.And(
			query.Eq("is_deleted", false),
			query.Or(
				query.Regex("title", keyword),
				query.Regex("description", keyword),
			),
		)).
		Sort(bson.D{{Key: "sort", Value: 1}, {Key: "created_at", Value: -1}}).
		Find(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = *doc
	}
	return results, nil
}

// SearchPublicDocumentContent 搜索公开文档内容
func (d *DocumentContentDao) SearchPublicDocumentContent(ctx context.Context, keyword string) ([]DocumentContent, error) {
	docs, err := d.coll.Finder().
		Filter(query.And(
			query.Eq("is_deleted", false),
			query.Or(
				query.Regex("title", keyword),
				query.Regex("description", keyword),
			),
		)).
		Sort(bson.D{{Key: "sort", Value: 1}, {Key: "created_at", Value: -1}}).
		Find(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = *doc
	}
	return results, nil
}

// FindPublicDocumentContentById 根据id查询公开文档内容
func (d *DocumentContentDao) FindPublicDocumentContentById(ctx context.Context, id bson.ObjectID) (DocumentContent, error) {
	doc, err := d.coll.Finder().
		Filter(query.And(query.Id(id), query.Eq("is_deleted", false))).
		FindOne(ctx)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return DocumentContent{}, errors.New("document_content中不存在文档")
	}
	if err != nil {
		return DocumentContent{}, err
	}
	return *doc, nil
}

// FindPublicDocumentContentByParentId 根据父级Id查询公开文档内容
func (d *DocumentContentDao) FindPublicDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]DocumentContent, error) {
	docs, err := d.coll.Finder().
		Filter(query.And(query.Eq("parent_id", parentId), query.Eq("is_deleted", false))).
		Sort(bson.D{{Key: "sort", Value: 1}, {Key: "created_at", Value: -1}}).
		Find(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = *doc
	}
	return results, nil
}

// FindPublicDocumentContentByDocumentId 根据文档Id查询公开文档内容
func (d *DocumentContentDao) FindPublicDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]DocumentContent, error) {
	docs, err := d.coll.Finder().
		Filter(query.And(query.Eq("document_id", documentId), query.Eq("is_deleted", false))).
		Sort(bson.D{{Key: "sort", Value: 1}, {Key: "created_at", Value: -1}}).
		Find(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = *doc
	}
	return results, nil
}

func (d *DocumentContentDao) DeleteDocumentContentList(ctx context.Context, ids []string) error {
	objIDs := make([]bson.ObjectID, 0, len(ids))
	for _, id := range ids {
		oid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("id格式错误: %s", id)
		}
		objIDs = append(objIDs, oid)
	}
	filter := bson.M{"_id": bson.M{"$in": objIDs}}
	_, err := d.coll.Deleter().Filter(filter).DeleteMany(ctx)
	return err
}

// FindPublicDocumentContentByRootIdAndAlias 根据根文档ID和别名查询公开文档内容
func (d *DocumentContentDao) FindPublicDocumentContentByRootIdAndAlias(ctx context.Context, documentId bson.ObjectID, alias string) (DocumentContent, error) {
	doc, err := d.coll.Finder().
		Filter(query.And(
			query.Eq("document_id", documentId),
			query.Eq("alias", alias),
			query.Eq("is_deleted", false),
		)).
		FindOne(ctx)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return DocumentContent{}, errors.New("document_content中不存在该别名的文档")
	}
	if err != nil {
		return DocumentContent{}, err
	}
	return *doc, nil
}

// JudgeDocumentContentAliasUnique 判断文档内容别名是否唯一
func (d *DocumentContentDao) JudgeDocumentContentAliasUnique(ctx context.Context, alias string, documentId bson.ObjectID) (bool, error) {
	count, err := d.coll.Finder().
		Filter(query.And(
			query.Eq("alias", alias),
			query.Eq("document_id", documentId),
		)).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}