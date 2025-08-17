package repository

import (
	"context"

	"github.com/codepzj/stellux/server/internal/document/internal/domain"
	"github.com/codepzj/stellux/server/internal/document/internal/repository/dao"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IDocumentRepository interface {
	CreateDocument(ctx context.Context, doc *domain.Document) (bson.ObjectID, error)
	FindDocumentById(ctx context.Context, id bson.ObjectID) (*domain.Document, error)
	UpdateDocumentById(ctx context.Context, id bson.ObjectID, doc *domain.Document) error
	DeleteDocumentById(ctx context.Context, id bson.ObjectID) error
	SoftDeleteDocumentById(ctx context.Context, id bson.ObjectID) error
	RestoreDocumentById(ctx context.Context, id bson.ObjectID) error
	FindDocumentByAlias(ctx context.Context, alias string) (*domain.Document, error)
	GetDocumentList(ctx context.Context, page *domain.Page) ([]*domain.Document, int64, error)
	GetDocumentBinList(ctx context.Context, page *domain.Page) ([]*domain.Document, int64, error)
	GetPublicDocumentList(ctx context.Context, page *domain.Page) ([]*domain.Document, int64, error)
	GetAllPublicDocuments(ctx context.Context) ([]*domain.Document, error)
}

var _ IDocumentRepository = (*DocumentRepository)(nil)

func NewDocumentRepository(dao dao.IDocumentDao) *DocumentRepository {
	return &DocumentRepository{dao: dao}
}

type DocumentRepository struct {
	dao dao.IDocumentDao
}

func (r *DocumentRepository) CreateDocument(ctx context.Context, doc *domain.Document) (bson.ObjectID, error) {
	return r.dao.CreateDocument(ctx, &dao.Document{
		Title:       doc.Title,
		Description: doc.Description,
		Thumbnail:   doc.Thumbnail,
		Alias:       doc.Alias,
		Sort:        doc.Sort,
		IsPublic:    doc.IsPublic,
		IsDeleted:   false,
	})
}

// FindDocumentById 根据条件查询文档
func (r *DocumentRepository) FindDocumentById(ctx context.Context, id bson.ObjectID) (*domain.Document, error) {
	doc, err := r.dao.FindDocumentById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &domain.Document{
		Id:          doc.ID,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
		DeletedAt:   doc.DeletedAt,
		Title:       doc.Title,
		Description: doc.Description,
		Thumbnail:   doc.Thumbnail,
		Alias:       doc.Alias,
		Sort:        doc.Sort,
		IsPublic:    doc.IsPublic,
		IsDeleted:   doc.IsDeleted,
	}, nil
}

// UpdateDocumentById 根据id更新文档
func (r *DocumentRepository) UpdateDocumentById(ctx context.Context, id bson.ObjectID, doc *domain.Document) error {
	return r.dao.UpdateDocumentById(ctx, id, &dao.Document{
		Title:       doc.Title,
		Description: doc.Description,
		Thumbnail:   doc.Thumbnail,
		Alias:       doc.Alias,
		Sort:        doc.Sort,
		IsPublic:    doc.IsPublic,
		IsDeleted:   doc.IsDeleted,
	})
}

// DeleteDocumentById 根据id删除文档
func (r *DocumentRepository) DeleteDocumentById(ctx context.Context, id bson.ObjectID) error {
	return r.dao.DeleteDocumentById(ctx, id)
}

// SoftDeleteDocumentById 根据id软删除文档
func (r *DocumentRepository) SoftDeleteDocumentById(ctx context.Context, id bson.ObjectID) error {
	return r.dao.SoftDeleteDocumentById(ctx, id)
}

// RestoreDocumentById 根据id恢复文档
func (r *DocumentRepository) RestoreDocumentById(ctx context.Context, id bson.ObjectID) error {
	return r.dao.RestoreDocumentById(ctx, id)
}

// FindDocumentByAlias 根据别名查询文档
func (r *DocumentRepository) FindDocumentByAlias(ctx context.Context, alias string) (*domain.Document, error) {
	doc, err := r.dao.FindDocumentByAlias(ctx, alias)
	if err != nil {
		return nil, err
	}
	return &domain.Document{
		Id:          doc.ID,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
		DeletedAt:   doc.DeletedAt,
		Title:       doc.Title,
		Description: doc.Description,
		Thumbnail:   doc.Thumbnail,
		Alias:       doc.Alias,
		Sort:        doc.Sort,
		IsPublic:    doc.IsPublic,
		IsDeleted:   doc.IsDeleted,
	}, nil
}

// GetDocumentList 获取文档列表
func (r *DocumentRepository) GetDocumentList(ctx context.Context, page *domain.Page) ([]*domain.Document, int64, error) {
	docs, count, err := r.dao.GetDocumentListByFilter(ctx, bson.D{{Key: "is_deleted", Value: false}}, &dao.Page{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	})
	if err != nil {
		return nil, 0, err
	}

	results := make([]*domain.Document, len(docs))
	for i, doc := range docs {
		results[i] = &domain.Document{
			Id:          doc.ID,
			CreatedAt:   doc.CreatedAt,
			UpdatedAt:   doc.UpdatedAt,
			DeletedAt:   doc.DeletedAt,
			Title:       doc.Title,
			Description: doc.Description,
			Thumbnail:   doc.Thumbnail,
			Alias:       doc.Alias,
			Sort:        doc.Sort,
			IsPublic:    doc.IsPublic,
			IsDeleted:   doc.IsDeleted,
		}
	}
	return results, count, nil
}

// GetDocumentBinList 获取文档回收箱列表
func (r *DocumentRepository) GetDocumentBinList(ctx context.Context, page *domain.Page) ([]*domain.Document, int64, error) {
	docs, count, err := r.dao.GetDocumentListByFilter(ctx, bson.D{{Key: "is_deleted", Value: true}}, &dao.Page{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	})
	if err != nil {
		return nil, 0, err
	}

	results := make([]*domain.Document, len(docs))
	for i, doc := range docs {
		results[i] = &domain.Document{
			Id:          doc.ID,
			CreatedAt:   doc.CreatedAt,
			UpdatedAt:   doc.UpdatedAt,
			DeletedAt:   doc.DeletedAt,
			Title:       doc.Title,
			Description: doc.Description,
			Thumbnail:   doc.Thumbnail,
			Alias:       doc.Alias,
			Sort:        doc.Sort,
			IsPublic:    doc.IsPublic,
			IsDeleted:   doc.IsDeleted,
		}
	}
	return results, count, nil
}

// GetPublicDocumentList 获取公开文档列表
func (r *DocumentRepository) GetPublicDocumentList(ctx context.Context, page *domain.Page) ([]*domain.Document, int64, error) {
	docs, count, err := r.dao.GetDocumentListByFilter(ctx, bson.D{{Key: "is_public", Value: true}, {Key: "is_deleted", Value: false}}, &dao.Page{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	})
	if err != nil {
		return nil, 0, err
	}

	results := make([]*domain.Document, len(docs))
	for i, doc := range docs {
		results[i] = &domain.Document{
			Id:          doc.ID,
			CreatedAt:   doc.CreatedAt,
			UpdatedAt:   doc.UpdatedAt,
			DeletedAt:   doc.DeletedAt,
			Title:       doc.Title,
			Description: doc.Description,
			Thumbnail:   doc.Thumbnail,
			Alias:       doc.Alias,
			Sort:        doc.Sort,
			IsPublic:    doc.IsPublic,
			IsDeleted:   doc.IsDeleted,
		}
	}
	return results, count, nil
}

// GetAllPublicDocuments 获取所有公开文档
func (r *DocumentRepository) GetAllPublicDocuments(ctx context.Context) ([]*domain.Document, error) {
	docs, err := r.dao.GetAllPublicDocuments(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*domain.Document, len(docs))
	for i, doc := range docs {
		results[i] = &domain.Document{
			Id:          doc.ID,
			CreatedAt:   doc.CreatedAt,
			UpdatedAt:   doc.UpdatedAt,
			DeletedAt:   doc.DeletedAt,
			Title:       doc.Title,
			Description: doc.Description,
			Thumbnail:   doc.Thumbnail,
			Alias:       doc.Alias,
			Sort:        doc.Sort,
			IsPublic:    doc.IsPublic,
			IsDeleted:   doc.IsDeleted,
		}
	}
	return results, nil
}
