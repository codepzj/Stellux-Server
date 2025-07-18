package repository

import (
	"context"

	"github.com/codepzj/stellux/server/internal/document_content/internal/domain"
	"github.com/codepzj/stellux/server/internal/document_content/internal/repository/dao"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IDocumentContentRepository interface {
	CreateDocumentContent(ctx context.Context, doc domain.DocumentContent) (bson.ObjectID, error)
	FindDocumentContentById(ctx context.Context, id bson.ObjectID) (domain.DocumentContent, error)
	DeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error
	FindDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]domain.DocumentContent, error)
	FindDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error)
	UpdateDocumentContentById(ctx context.Context, id bson.ObjectID, doc domain.DocumentContent) error
}

var _ IDocumentContentRepository = (*DocumentContentRepository)(nil)

func NewDocumentContentRepository(dao dao.IDocumentContentDao) *DocumentContentRepository {
	return &DocumentContentRepository{dao: dao}
}

type DocumentContentRepository struct {
	dao dao.IDocumentContentDao
}

// CreateDocumentContent 创建文档内容
func (r *DocumentContentRepository) CreateDocumentContent(ctx context.Context, doc domain.DocumentContent) (bson.ObjectID, error) {
	return r.dao.CreateDocumentContent(ctx, dao.DocumentContent{
		DocumentId:   doc.DocumentId,
		Title:        doc.Title,
		Content:      doc.Content,
		Description:  doc.Description,
		Version:      doc.Version,
		Alias:        doc.Alias,
		ParentId:     doc.ParentId,
		IsDir:        doc.IsDir,
		Sort:         doc.Sort,
	})
}

// FindDocumentContentById 根据id查询文档内容
func (r *DocumentContentRepository) FindDocumentContentById(ctx context.Context, id bson.ObjectID) (domain.DocumentContent, error) {
	doc, err := r.dao.FindDocumentContentById(ctx, id)
	if err != nil {
		return domain.DocumentContent{}, err
	}

	return domain.DocumentContent{
		Id:           doc.ID,
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
		DocumentId:   doc.DocumentId,
		Title:        doc.Title,
		Content:      doc.Content,
		Description:  doc.Description,
		Version:      doc.Version,
		Alias:        doc.Alias,
		ParentId:     doc.ParentId,
		IsDir:        doc.IsDir,
		Sort:         doc.Sort,
		LikeCount:    doc.LikeCount,
		DislikeCount: doc.DislikeCount,
		CommentCount: doc.CommentCount,
	}, nil
}

// DeleteDocumentContentById 根据id删除文档内容
func (r *DocumentContentRepository) DeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	return r.dao.DeleteDocumentContentById(ctx, id)
}

// FindDocumentContentByParentId 根据父级ID查询文档内容
func (r *DocumentContentRepository) FindDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]domain.DocumentContent, error) {
	docs, err := r.dao.FindDocumentContentByParentId(ctx, parentId)
	if err != nil {
		return nil, err
	}
	results := make([]domain.DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = domain.DocumentContent{
			Id:           doc.ID,
			CreatedAt:    doc.CreatedAt,
			UpdatedAt:    doc.UpdatedAt,
			DocumentId:   doc.DocumentId,
			Title:        doc.Title,
			Content:      doc.Content,
			Description:  doc.Description,
			Version:      doc.Version,
			Alias:        doc.Alias,
			ParentId:     doc.ParentId,
			IsDir:        doc.IsDir,
			Sort:         doc.Sort,
			LikeCount:    doc.LikeCount,
			DislikeCount: doc.DislikeCount,
			CommentCount: doc.CommentCount,
		}
	}
	return results, nil
}

func (r *DocumentContentRepository) FindDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error) {
	docs, err := r.dao.FindDocumentContentByDocumentId(ctx, documentId)
	if err != nil {
		return nil, err
	}
	results := make([]domain.DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = domain.DocumentContent{
			Id:           doc.ID,
			CreatedAt:    doc.CreatedAt,
			UpdatedAt:    doc.UpdatedAt,
			DocumentId:   doc.DocumentId,
			Title:        doc.Title,
			Content:      doc.Content,
			Description:  doc.Description,
			Version:      doc.Version,
			Alias:        doc.Alias,
			ParentId:     doc.ParentId,
			IsDir:        doc.IsDir,
			Sort:         doc.Sort,
			LikeCount:    doc.LikeCount,
			DislikeCount: doc.DislikeCount,
			CommentCount: doc.CommentCount,
		}
	}
	return results, nil
}

func (r *DocumentContentRepository) UpdateDocumentContentById(ctx context.Context, id bson.ObjectID, doc domain.DocumentContent) error {
	return r.dao.UpdateDocumentContentById(ctx, id, dao.DocumentContent{
		DocumentId:   doc.DocumentId,
		Title:        doc.Title,
		Content:      doc.Content,
		Description:  doc.Description,
		Version:      doc.Version,
		Alias:        doc.Alias,
		ParentId:     doc.ParentId,
		IsDir:        doc.IsDir,
		Sort:         doc.Sort,
	})
}
