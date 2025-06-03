package repository

import (
	"context"

	"github.com/codepzj/stellux/server/internal/document/internal/domain"
	"github.com/codepzj/stellux/server/internal/document/internal/repository/dao"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IDocumentRepository interface {
	Create(ctx context.Context, doc *domain.Document) error
	CreateRoot(ctx context.Context, doc *domain.DocumentRoot) error
	FindByKeyword(ctx context.Context, keyword string, documentID bson.ObjectID) ([]*domain.Document, error)
	FindDocumentIsPublic(ctx context.Context, documentID bson.ObjectID) (bool, error)
	FindAllRootDocument(ctx context.Context) ([]*domain.DocumentRoot, error)
	FindAllPublicRootDocument(ctx context.Context) ([]*domain.DocumentRoot, error)
    FindAllByTypeAndDocumentID(ctx context.Context,document_type string,documentID bson.ObjectID) ([]*domain.Document, error)
	FindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error)
	FindAllByDocumentIDList(ctx context.Context, documentIDList []bson.ObjectID) ([]*domain.Document, error)
	FindAllPublicByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error)
	UpdateDocumentByID(ctx context.Context, id bson.ObjectID, title string, content string) error
	DeleteByID(ctx context.Context, id bson.ObjectID) error
	GetDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.Document, error)	
	GetRootDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.DocumentRoot, error)
	DeleteByIDList(ctx context.Context, idList []bson.ObjectID) error
	RenameDocumentByID(ctx context.Context, id bson.ObjectID, title string) error
	UpdateRootDocumentByID(ctx context.Context, id bson.ObjectID, doc *domain.DocumentRoot) error
	DeleteRootDocumentByID(ctx context.Context, id bson.ObjectID) error
}

var _ IDocumentRepository = (*DocumentRepository)(nil)

func NewDocumentRepository(dao dao.IDocumentDao) *DocumentRepository {
	return &DocumentRepository{dao: dao}
}

type DocumentRepository struct {
	dao dao.IDocumentDao
}

// 新增文档
func (r *DocumentRepository) Create(ctx context.Context, doc *domain.Document) error {
	return r.dao.Create(ctx, r.DomainToDao(doc))
}

// 新增根文档
func (r *DocumentRepository) CreateRoot(ctx context.Context, doc *domain.DocumentRoot) error {
	return r.dao.CreateRoot(ctx, r.DomainToRootDao(doc))
}

// 更新根文档
func (r *DocumentRepository) UpdateRootDocumentByID(ctx context.Context, id bson.ObjectID, doc *domain.DocumentRoot) error {
	return r.dao.UpdateRootDocumentByID(ctx, id, r.DomainToRootDao(doc))
}

// 删除根文档
func (r *DocumentRepository) DeleteRootDocumentByID(ctx context.Context, id bson.ObjectID) error {
	return r.dao.DeleteRootDocumentByID(ctx, id)
}

// 根据关键词查询文档
func (r *DocumentRepository) FindByKeyword(ctx context.Context, keyword string, documentID bson.ObjectID) ([]*domain.Document, error) {
	documentList, err := r.dao.FindByKeyword(ctx, keyword, documentID)
	if err != nil {
		return nil, err
	}
	return r.DaoToDomainList(documentList), nil
}

// 根据id查询文档是否公开
func (r *DocumentRepository) FindDocumentIsPublic(ctx context.Context, documentID bson.ObjectID) (bool, error) {
	return r.dao.FindDocumentIsPublic(ctx, documentID)
}

// 查询所有根文档
func (r *DocumentRepository) FindAllRootDocument(ctx context.Context) ([]*domain.DocumentRoot, error) {
	documentList, err := r.dao.FindAllByType(ctx, "root")
	if err != nil {
		return nil, err
	}
	return r.DaoToRootDomainList(documentList), nil
}

// 查询所有公开根文档
func (r *DocumentRepository) FindAllPublicRootDocument(ctx context.Context) ([]*domain.DocumentRoot, error) {
	documentList, err := r.dao.FindAllPublicByType(ctx, "root")
	if err != nil {
		return nil, err
	}
	return r.DaoToRootDomainList(documentList), nil
}

// 根据文档类型和文档id查询文档中的子文档
func (r *DocumentRepository) FindAllByTypeAndDocumentID(ctx context.Context, document_type string, documentID bson.ObjectID) ([]*domain.Document, error) {
	documentList, err := r.dao.FindAllByTypeAndDocumentID(ctx, document_type, documentID)
	if err != nil {
		return nil, err
	}
	return r.DaoToDomainList(documentList), nil
}

// 根据文档id查询文档中的子文档
func (r *DocumentRepository) FindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error) {
	documentList, err := r.dao.FindAllByDocumentID(ctx, documentID)
	if err != nil {
		return nil, err
	}
	return r.DaoToDomainList(documentList), nil
}

// 根据文档id查询文档中的公开子文档
func (r *DocumentRepository) FindAllPublicByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error) {
	documentList, err := r.dao.FindAllPublicByDocumentID(ctx, documentID)
	if err != nil {
		return nil, err
	}
	return r.DaoToDomainList(documentList), nil
}

// 根据id列表查询文档
func (r *DocumentRepository) FindAllByDocumentIDList(ctx context.Context, documentIDList []bson.ObjectID) ([]*domain.Document, error) {
	documentList, err := r.dao.FindAllByDocumentIDList(ctx, documentIDList)
	if err != nil {
		return nil, err
	}
	return r.DaoToDomainList(documentList), nil
}

// 根据id查询文档
func (r *DocumentRepository) GetDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.Document, error) {
	document, err := r.dao.GetDocumentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.DaoToDomain(document), nil
}

// 根据id查询根文档
func (r *DocumentRepository) GetRootDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.DocumentRoot, error) {
	document, err := r.dao.GetDocumentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.DaoToRootDomain(document), nil
}

// 更新文档
func (r *DocumentRepository) UpdateDocumentByID(ctx context.Context, id bson.ObjectID, title string, content string) error {
	return r.dao.UpdateDocumentByID(ctx, id, title, content)
}

// 重命名文档
func (r *DocumentRepository) RenameDocumentByID(ctx context.Context, id bson.ObjectID, title string) error {
	return r.dao.RenameDocumentByID(ctx, id, title)
}

func (r *DocumentRepository) DeleteByID(ctx context.Context, id bson.ObjectID) error {
	return r.dao.DeleteByID(ctx, id)
}

func (r *DocumentRepository) DeleteByIDList(ctx context.Context, idList []bson.ObjectID) error {
	return r.dao.DeleteByIDList(ctx, idList)
}

// 将文章节点的domain转换为dao
func (r *DocumentRepository) DomainToDao(doc *domain.Document) *dao.Document {
	return &dao.Document{
		Title:      doc.Title,
		Content:    doc.Content,
		DocumentType: doc.DocumentType,
		ParentID:   doc.ParentID,
		DocumentID: doc.DocumentID,
	}
}

// 将文章根节点的domain转换为dao
func (r *DocumentRepository) DomainToRootDao(doc *domain.DocumentRoot) *dao.Document {
	return &dao.Document{
		Title:      doc.Title,
		Alias:      doc.Alias,
		Description: doc.Description,
		Thumbnail:   doc.Thumbnail,
		DocumentType: doc.DocumentType,
		IsPublic:     doc.IsPublic,
	}
}

func (r *DocumentRepository) DaoToDomain(doc *dao.Document) *domain.Document {
	return &domain.Document{
		ID:         doc.ID,
		CreatedAt:  doc.CreatedAt,
		UpdatedAt:  doc.UpdatedAt,
		Title:      doc.Title,
		Content:    doc.Content,
		DocumentType: doc.DocumentType,
		ParentID:   doc.ParentID,
		DocumentID: doc.DocumentID,
	}
}

func (r *DocumentRepository) DaoToRootDomain(doc *dao.Document) *domain.DocumentRoot {
	return &domain.DocumentRoot{
		ID:         doc.ID,
		CreatedAt:  doc.CreatedAt,
		UpdatedAt:  doc.UpdatedAt,
		Title:      doc.Title,
		Alias:      doc.Alias,
		Description: doc.Description,
		Thumbnail: doc.Thumbnail,
		IsPublic: doc.IsPublic,
		DocumentType: doc.DocumentType,
	}
}

func (r *DocumentRepository) DaoToDomainList(docList []*dao.Document) []*domain.Document {
	return lo.Map(docList, func(doc *dao.Document, _ int) *domain.Document {
		return r.DaoToDomain(doc)
	})
}

func (r *DocumentRepository) DaoToRootDomainList(docList []*dao.Document) []*domain.DocumentRoot {
	return lo.Map(docList, func(doc *dao.Document, _ int) *domain.DocumentRoot {
		return r.DaoToRootDomain(doc)
	})
}