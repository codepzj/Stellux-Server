package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DocumentContent struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	CreatedAt   time.Time     `bson:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"`
	DeletedAt   *time.Time    `bson:"deleted_at,omitempty"`
	DocumentId  bson.ObjectID `bson:"document_id"` // 根文档Id
	Title       string        `bson:"title"`       // 文档标题
	Content     string        `bson:"content"`     // 文档内容
	Description string        `bson:"description"` // 文档描述
	Alias       string        `bson:"alias"`       // 文档别名
	ParentId    bson.ObjectID `bson:"parent_id"`   // 父级Id
	IsDir       bool          `bson:"is_dir"`      // 是否是目录
	Sort        int           `bson:"sort"`        // 排序
	IsDeleted   bool          `bson:"is_deleted"`  // 是否删除
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
	GetDocumentContentList(ctx context.Context, page *apiwrap.Page, documentId bson.ObjectID) ([]*DocumentContent, int64, error)
	GetPublicDocumentContentListByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]*DocumentContent, error)
	SearchDocumentContent(ctx context.Context, keyword string) ([]DocumentContent, error)
	SearchPublicDocumentContent(ctx context.Context, keyword string) ([]DocumentContent, error)
	FindPublicDocumentContentById(ctx context.Context, id bson.ObjectID) (DocumentContent, error)
	FindPublicDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]DocumentContent, error)
	FindPublicDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]DocumentContent, error)
	FindPublicDocumentContentByRootIdAndAlias(ctx context.Context, documentId bson.ObjectID, alias string) (DocumentContent, error)
	DeleteDocumentContentList(ctx context.Context, ids []string) error
	GetDocumentContentListByAlias(ctx context.Context, alias string, documentId bson.ObjectID) ([]*DocumentContent, error)
}

type DocumentContentDao struct {
	coll *mongo.Collection
}

func NewDocumentContentDao(db *mongo.Database) *DocumentContentDao {
	return &DocumentContentDao{coll: db.Collection("document_content")}
}

// CreateDocumentContent 创建文档内容
func (d *DocumentContentDao) CreateDocumentContent(ctx context.Context, doc DocumentContent) (bson.ObjectID, error) {
	doc.ID = bson.NewObjectID()
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()
	result, err := d.coll.InsertOne(ctx, &doc)
	if err != nil {
		return bson.ObjectID{}, err
	}
	return result.InsertedID.(bson.ObjectID), nil
}

// FindDocumentContentById 根据id查询文档内容
func (d *DocumentContentDao) FindDocumentContentById(ctx context.Context, id bson.ObjectID) (DocumentContent, error) {
	var doc DocumentContent
	err := d.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return DocumentContent{}, errors.New("document_content中不存在文档")
	}
	if err != nil {
		return DocumentContent{}, err
	}
	return doc, nil
}

// DeleteDocumentContentById 根据id删除文档内容
func (d *DocumentContentDao) DeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	result, err := d.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("文档不存在")
	}
	return nil
}

// SoftDeleteDocumentContentById 根据id软删除文档内容
func (d *DocumentContentDao) SoftDeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"is_deleted": true,
			"updated_at": now,
		},
	}
	result, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
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
	update := bson.M{
		"$set": bson.M{
			"is_deleted": false,
			"updated_at": time.Now(),
		},
		"$unset": bson.M{
			"deleted_at": "",
		},
	}
	result, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
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
	cursor, err := d.coll.Find(ctx, bson.M{"parent_id": parentId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []DocumentContent
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// FindDocumentContentByDocumentId 根据文档Id查询文档内容
func (d *DocumentContentDao) FindDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]DocumentContent, error) {
	cursor, err := d.coll.Find(ctx, bson.M{"document_id": documentId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []DocumentContent
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// UpdateDocumentContentById 根据id更新文档内容
func (d *DocumentContentDao) UpdateDocumentContentById(ctx context.Context, id bson.ObjectID, doc DocumentContent) error {
	update := bson.M{
		"$set": bson.M{
			"document_id": doc.DocumentId,
			"title":       doc.Title,
			"content":     doc.Content,
			"description": doc.Description,
			"alias":       doc.Alias,
			"parent_id":   doc.ParentId,
			"is_dir":      doc.IsDir,
			"sort":        doc.Sort,
			"is_deleted":  doc.IsDeleted,
			"updated_at":  time.Now(),
		},
	}
	result, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("文档不存在")
	}
	return nil
}

// GetDocumentContentList 获取文档内容列表
func (d *DocumentContentDao) GetDocumentContentList(ctx context.Context, page *apiwrap.Page, documentId bson.ObjectID) ([]*DocumentContent, int64, error) {
	skip := (page.PageNo - 1) * page.PageSize
	filter := bson.M{
		"document_id": documentId,
		"is_deleted":  false,
	}

	count, err := d.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "sort", Value: 1}, {Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(page.PageSize)

	cursor, err := d.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*DocumentContent
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

// GetPublicDocumentContentList 获取公开文档内容列表
func (d *DocumentContentDao) GetPublicDocumentContentListByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]*DocumentContent, error) {
	filter := bson.M{
		"document_id": documentId,
		"is_deleted":  false,
	}

	cursor, err := d.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []*DocumentContent
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, err
	}

	return documents, nil
}

// SearchDocumentContent 搜索文档内容
func (d *DocumentContentDao) SearchDocumentContent(ctx context.Context, keyword string) ([]DocumentContent, error) {
	filter := bson.M{
		"is_deleted": false,
		"$or": []bson.M{
			{"title": bson.M{"$regex": keyword, "$options": "i"}},
			{"description": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "sort", Value: 1}, {Key: "created_at", Value: -1}})

	cursor, err := d.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []DocumentContent
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// SearchPublicDocumentContent 搜索公开文档内容
func (d *DocumentContentDao) SearchPublicDocumentContent(ctx context.Context, keyword string) ([]DocumentContent, error) {
	filter := bson.M{
		"is_deleted": false,
		"$or": []bson.M{
			{"title": bson.M{"$regex": keyword, "$options": "i"}},
			{"description": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "sort", Value: 1}, {Key: "created_at", Value: -1}})

	cursor, err := d.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []DocumentContent
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// FindPublicDocumentContentById 根据id查询公开文档内容
func (d *DocumentContentDao) FindPublicDocumentContentById(ctx context.Context, id bson.ObjectID) (DocumentContent, error) {
	filter := bson.M{
		"_id":        id,
		"is_deleted": false,
	}
	var doc DocumentContent
	err := d.coll.FindOne(ctx, filter).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return DocumentContent{}, errors.New("document_content中不存在文档")
	}
	if err != nil {
		return DocumentContent{}, err
	}
	return doc, nil
}

// FindPublicDocumentContentByParentId 根据父级Id查询公开文档内容
func (d *DocumentContentDao) FindPublicDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]DocumentContent, error) {
	filter := bson.M{
		"parent_id":  parentId,
		"is_deleted": false,
	}
	opts := options.Find().SetSort(bson.D{{Key: "sort", Value: 1}, {Key: "created_at", Value: -1}})

	cursor, err := d.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []DocumentContent
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// FindPublicDocumentContentByDocumentId 根据文档Id查询公开文档内容
func (d *DocumentContentDao) FindPublicDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]DocumentContent, error) {
	filter := bson.M{
		"document_id": documentId,
		"is_deleted":  false,
	}
	opts := options.Find().SetSort(bson.D{{Key: "sort", Value: 1}, {Key: "created_at", Value: -1}})

	cursor, err := d.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []DocumentContent
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
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
	_, err := d.coll.DeleteMany(ctx, filter)
	return err
}

// FindPublicDocumentContentByRootIdAndAlias 根据根文档ID和别名查询公开文档内容
func (d *DocumentContentDao) FindPublicDocumentContentByRootIdAndAlias(ctx context.Context, documentId bson.ObjectID, alias string) (DocumentContent, error) {
	filter := bson.M{
		"document_id": documentId,
		"alias":       alias,
		"is_deleted":  false,
	}
	var doc DocumentContent
	err := d.coll.FindOne(ctx, filter).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return DocumentContent{}, errors.New("document_content中不存在该别名的文档")
	}
	if err != nil {
		return DocumentContent{}, err
	}
	return doc, nil
}

// GetDocumentContentListByAlias 获取文档内容别名数量
func (d *DocumentContentDao) GetDocumentContentListByAlias(ctx context.Context, alias string, documentId bson.ObjectID) ([]*DocumentContent, error) {
	filter := bson.M{
		"alias":       alias,
		"document_id": documentId,
	}

	cursor, err := d.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docContentList []*DocumentContent
	if err = cursor.All(ctx, &docContentList); err != nil {
		return nil, err
	}
	return docContentList, nil
}
