package repository

import (
	"context"

	"github.com/codepzj/stellux/server/internal/document_content/internal/domain"
	"github.com/codepzj/stellux/server/internal/document_content/internal/repository/dao"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IDocumentContentRepository interface {
	CreateDocumentContent(ctx context.Context, doc *domain.DocumentContent) error
	FindDocumentContentByDocumentID(ctx context.Context, documentId bson.ObjectID) ([]*domain.DocumentContent, error)
}

var _ IDocumentContentRepository = (*DocumentContentRepository)(nil)

func NewDocumentContentRepository(dao dao.IDocumentContentDao) *DocumentContentRepository {
	return &DocumentContentRepository{dao: dao}
}

type DocumentContentRepository struct {
	dao dao.IDocumentContentDao
}

// CreateDocumentContent 创建文档内容
func (r *DocumentContentRepository) CreateDocumentContent(ctx context.Context, doc *domain.DocumentContent) error {
	return r.dao.CreateDocumentContent(ctx, &dao.DocumentContent{
		DocumentId: doc.DocumentId,
		Title: doc.Title,
		Content: doc.Content,
		Version: doc.Version,
		Alias: doc.Alias,
		ParentID: doc.ParentID,
		IsDir: doc.IsDir,
		LikeCount: doc.LikeCount,
		DislikeCount: doc.DislikeCount,
		CommentCount: doc.CommentCount,
	})
}

// FindDocumentContentByDocumentID 根据文档ID查询文档内容
func (r *DocumentContentRepository) FindDocumentContentByDocumentID(ctx context.Context, documentId bson.ObjectID) ([]*domain.DocumentContent, error) {
	docs, err := r.dao.FindDocumentContentByDocumentID(ctx, documentId)
	if err != nil {
		return nil, err
	}

	results := make([]*domain.DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = &domain.DocumentContent{
			Id: doc.DocumentId,
			Title: doc.Title,
			Content: doc.Content,
			Version: doc.Version,
			Alias: doc.Alias,
			ParentID: doc.ParentID,
			IsDir: doc.IsDir,
			LikeCount: doc.LikeCount,
			DislikeCount: doc.DislikeCount,
			CommentCount: doc.CommentCount,
		}
	}
	return results, nil
}