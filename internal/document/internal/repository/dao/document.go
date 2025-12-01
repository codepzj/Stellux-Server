package dao

import (
	"context"
	"time"

	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Document struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	CreatedAt   time.Time     `bson:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"`
	DeletedAt   *time.Time    `bson:"deleted_at,omitempty"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	Thumbnail   string        `bson:"thumbnail"`
	Alias       string        `bson:"alias"`
	Sort        int           `bson:"sort"`
	IsPublic    bool          `bson:"is_public"`
	IsDeleted   bool          `bson:"is_deleted"`
}

type IDocumentDao interface {
	CreateDocument(ctx context.Context, doc *Document) (bson.ObjectID, error)
	FindDocumentById(ctx context.Context, id bson.ObjectID) (*Document, error)
	UpdateDocumentById(ctx context.Context, id bson.ObjectID, doc *Document) error
	DeleteDocumentById(ctx context.Context, id bson.ObjectID) error
	SoftDeleteDocumentById(ctx context.Context, id bson.ObjectID) error
	RestoreDocumentById(ctx context.Context, id bson.ObjectID) error
	FindDocumentByAlias(ctx context.Context, alias string) (*Document, error)
	GetDocumentListByFilter(ctx context.Context, filter bson.D, page *apiwrap.Page) ([]*Document, int64, error)
	GetAllPublicDocuments(ctx context.Context) ([]*Document, error)
}

var _ IDocumentDao = (*DocumentDao)(nil)

func NewDocumentDao(db *mongo.Database) *DocumentDao {
	return &DocumentDao{coll: db.Collection("document")}
}

type DocumentDao struct {
	coll *mongo.Collection
}

// CreateDocument 创建文档
func (d *DocumentDao) CreateDocument(ctx context.Context, doc *Document) (bson.ObjectID, error) {
	doc.ID = bson.NewObjectID()
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()
	result, err := d.coll.InsertOne(ctx, doc)
	if err != nil {
		return bson.ObjectID{}, err
	}
	return result.InsertedID.(bson.ObjectID), nil
}

// FindDocumentById 根据ID查询文档
func (d *DocumentDao) FindDocumentById(ctx context.Context, id bson.ObjectID) (*Document, error) {
	var document Document
	err := d.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&document)
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// UpdateDocumentById 更新文档
func (d *DocumentDao) UpdateDocumentById(ctx context.Context, id bson.ObjectID, doc *Document) error {
	update := bson.M{
		"$set": bson.M{
			"title":       doc.Title,
			"description": doc.Description,
			"thumbnail":   doc.Thumbnail,
			"alias":       doc.Alias,
			"sort":        doc.Sort,
			"is_public":   doc.IsPublic,
			"updated_at":  time.Now(),
		},
	}
	result, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("文档不存在, 更新失败")
	}
	return nil
}

// DeleteDocumentById 根据ID删除文档
func (d *DocumentDao) DeleteDocumentById(ctx context.Context, id bson.ObjectID) error {
	result, err := d.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("文档不存在, 删除失败")
	}
	return nil
}

// SoftDeleteDocumentById 根据ID软删除文档
func (d *DocumentDao) SoftDeleteDocumentById(ctx context.Context, id bson.ObjectID) error {
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

// RestoreDocumentById 根据ID恢复文档
func (d *DocumentDao) RestoreDocumentById(ctx context.Context, id bson.ObjectID) error {
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

// FindDocumentByAlias 根据别名查询文档
func (d *DocumentDao) FindDocumentByAlias(ctx context.Context, alias string) (*Document, error) {
	var document Document
	err := d.coll.FindOne(ctx, bson.M{"alias": alias}).Decode(&document)
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// GetDocumentListByFilter 根据过滤条件获取文档列表
func (d *DocumentDao) GetDocumentListByFilter(ctx context.Context, filter bson.D, page *apiwrap.Page) ([]*Document, int64, error) {
	skip := (page.PageNo - 1) * page.PageSize

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

	var documents []*Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, 0, err
	}

	return documents, count, nil
}

// GetAllPublicDocuments 获取所有公开文档
func (d *DocumentDao) GetAllPublicDocuments(ctx context.Context) ([]*Document, error) {
	filter := bson.M{
		"is_public":  true,
		"is_deleted": false,
	}
	opts := options.Find().SetSort(bson.D{{Key: "sort", Value: 1}, {Key: "created_at", Value: -1}})

	cursor, err := d.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documents []*Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, err
	}
	return documents, nil
}
