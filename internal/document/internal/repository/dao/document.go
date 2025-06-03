package dao

import (
	"context"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Document struct {
	mongox.Model `bson:",inline"`
	Title        string
	Content      string
	Alias        string
	Description  string
	Thumbnail    string
	DocumentType string        `bson:"document_type"`
	IsPublic     bool          `bson:"is_public"`
	ParentID     bson.ObjectID `bson:"parent_id,omitempty"`   // 根节点不需要parent_id
	DocumentID   bson.ObjectID `bson:"document_id,omitempty"` // 根节点不需要document_id
}

type IDocumentDao interface {
	Create(ctx context.Context, doc *Document) error
	CreateRoot(ctx context.Context, doc *Document) error
	FindDocumentIsPublic(ctx context.Context, documentID bson.ObjectID) (bool, error)
	FindByKeyword(ctx context.Context, keyword string, documentID bson.ObjectID) ([]*Document, error)
	FindAllByType(ctx context.Context, document_type string) ([]*Document, error)
	FindAllPublicByType(ctx context.Context, document_type string) ([]*Document, error)
	FindAllByTypeAndDocumentID(ctx context.Context, document_type string, documentID bson.ObjectID) ([]*Document, error)
	FindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*Document, error)
	FindAllPublicByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*Document, error)
	FindAllByDocumentIDList(ctx context.Context, documentIDList []bson.ObjectID) ([]*Document, error)
	GetDocumentByID(ctx context.Context, id bson.ObjectID) (*Document, error)
	UpdateRootDocumentByID(ctx context.Context, id bson.ObjectID, doc *Document) error
	UpdateDocumentByID(ctx context.Context, id bson.ObjectID, title string, content string) error
	RenameDocumentByID(ctx context.Context, id bson.ObjectID, title string) error
	DeleteByID(ctx context.Context, id bson.ObjectID) error
	DeleteByIDList(ctx context.Context, idList []bson.ObjectID) error
	DeleteRootDocumentByID(ctx context.Context, id bson.ObjectID) error
}

var _ IDocumentDao = (*DocumentDao)(nil)

func NewDocumentDao(db *mongox.Database) *DocumentDao {
	return &DocumentDao{coll: mongox.NewCollection[Document](db, "document")}
}

type DocumentDao struct {
	coll *mongox.Collection[Document]
}

// 新增根文档
func (d *DocumentDao) CreateRoot(ctx context.Context, doc *Document) error {
	insertResult, err := d.coll.Creator().InsertOne(ctx, &Document{
		Title:        doc.Title,
		Alias:        doc.Alias,
		Description:  doc.Description,
		Thumbnail:    doc.Thumbnail,
		IsPublic:     doc.IsPublic,
		DocumentType: "root",
	})
	if err != nil {
		return err
	}
	if insertResult.InsertedID == nil {
		return errors.New("新增根文档失败")
	}
	return nil
}

// 新增文档
func (d *DocumentDao) Create(ctx context.Context, doc *Document) error {
	insertResult, err := d.coll.Creator().InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	if insertResult.InsertedID == nil {
		return errors.New("新增文档失败")
	}
	return nil
}

// 更新根文档
func (d *DocumentDao) UpdateRootDocumentByID(ctx context.Context, id bson.ObjectID, doc *Document) error {
	updateResult, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.NewBuilder().Set("title", doc.Title).Set("alias", doc.Alias).Set("description", doc.Description).Set("thumbnail", doc.Thumbnail).Set("is_public", doc.IsPublic).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrap(err, "更新根文档失败")
	}
	if updateResult.MatchedCount == 0 {
		return errors.New("更新根文档失败")
	}
	return nil
}

// 删除根文档
func (d *DocumentDao) DeleteRootDocumentByID(ctx context.Context, id bson.ObjectID) error {
	_, err := d.coll.Deleter().Filter(query.Or(query.Id(id), query.Eq("document_id", id))).DeleteMany(ctx)
	if err != nil {
		return errors.Wrap(err, "删除根文档失败")
	}
	return nil
}

// 根据关键词查询文档

func (d *DocumentDao) FindByKeyword(ctx context.Context, keyword string, documentID bson.ObjectID) ([]*Document, error) {
	documentList, err := d.coll.Finder().Filter(query.NewBuilder().Or(query.RegexOptions("title", keyword, "i"), query.RegexOptions("description", keyword, "i"), query.RegexOptions("content", keyword, "i")).And(query.Eq("document_id", documentID)).Build()).Find(ctx)
	if err != nil {
		return nil, err
	}
	return documentList, nil
}

// 根据id查询文档是否公开
func (d *DocumentDao) FindDocumentIsPublic(ctx context.Context, documentID bson.ObjectID) (bool, error) {
	document, err := d.coll.Finder().Filter(query.Id(documentID)).FindOne(ctx)
	if err != nil {
		return false, err
	}
	return document.IsPublic, nil
}

// 根据文档类型查询文档
func (d *DocumentDao) FindAllByType(ctx context.Context, document_type string) ([]*Document, error) {
	documentList, err := d.coll.Finder().Filter(query.NewBuilder().Eq("document_type", document_type).Build()).Sort(bson.M{"updated_at": -1}).Find(ctx)
	if err != nil {
		return nil, err
	}
	return documentList, nil
}

// 根据文档类型查询公开文档
func (d *DocumentDao) FindAllPublicByType(ctx context.Context, document_type string) ([]*Document, error) {
	documentList, err := d.coll.Finder().Filter(query.NewBuilder().Eq("document_type", document_type).Eq("is_public", true).Build()).Sort(bson.M{"updated_at": -1}).Find(ctx)
	if err != nil {
		return nil, err
	}
	return documentList, nil
}

// 根据文档类型和文档id查询文档
func (d *DocumentDao) FindAllByTypeAndDocumentID(ctx context.Context, document_type string, documentID bson.ObjectID) ([]*Document, error) {
	documentList, err := d.coll.Finder().Filter(query.NewBuilder().Eq("document_type", document_type).Eq("document_id", documentID).Build()).Find(ctx)
	if err != nil {
		return nil, err
	}
	return documentList, nil
}

// 根据文档id查询子文档
func (d *DocumentDao) FindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*Document, error) {
	documentList, err := d.coll.Finder().Filter(query.Eq("document_id", documentID)).Find(ctx)
	if err != nil {
		return nil, err
	}
	return documentList, nil
}

// 根据文档id查询文档中的公开子文档
func (d *DocumentDao) FindAllPublicByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*Document, error) {
	documentList, err := d.coll.Finder().Filter(query.NewBuilder().Eq("document_id", documentID).Eq("is_public", true).Build()).Find(ctx)
	if err != nil {
		return nil, err
	}
	return documentList, nil
}

// 根据document_id列表查询文档
func (d *DocumentDao) FindAllByDocumentIDList(ctx context.Context, documentIDList []bson.ObjectID) ([]*Document, error) {
	documentList, err := d.coll.Finder().Filter(query.In("document_id", documentIDList...)).Find(ctx)
	if err != nil {
		return nil, err
	}
	return documentList, nil
}

// 根据id查询文档
func (d *DocumentDao) GetDocumentByID(ctx context.Context, id bson.ObjectID) (*Document, error) {
	document, err := d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
	if err != nil {
		return nil, err
	}
	return document, nil
}

// 更新文档
func (d *DocumentDao) UpdateDocumentByID(ctx context.Context, id bson.ObjectID, title string, content string) error {
	_, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.NewBuilder().Set("title", title).Set("content", content).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrap(err, "更新文档失败")
	}
	return nil
}

// 重命名文档
func (d *DocumentDao) RenameDocumentByID(ctx context.Context, id bson.ObjectID, title string) error {
	_, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.NewBuilder().Set("title", title).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrap(err, "重命名文档失败")
	}
	return nil
}

// 删除文档
func (d *DocumentDao) DeleteByID(ctx context.Context, id bson.ObjectID) error {
	_, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	if err != nil {
		return errors.Wrap(err, "删除文档失败")
	}
	return nil
}

// 删除多个文档
func (d *DocumentDao) DeleteByIDList(ctx context.Context, idList []bson.ObjectID) error {
	_, err := d.coll.Deleter().Filter(query.In("_id", idList...)).DeleteMany(ctx)
	if err != nil {
		return errors.Wrap(err, "删除多个文档失败")
	}
	return nil
}
