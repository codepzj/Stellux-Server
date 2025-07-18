package service

import (
	"context"

	"github.com/codepzj/stellux/server/internal/document_content/internal/domain"
	"github.com/codepzj/stellux/server/internal/document_content/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)


type IDocumentContentService interface {
	CreateDocumentContent(ctx context.Context, doc domain.DocumentContent) (bson.ObjectID, error)
	FindDocumentContentById(ctx context.Context, id bson.ObjectID) (domain.DocumentContent, error)
	DeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error
	FindDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]domain.DocumentContent, error)
	FindDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error)
	UpdateDocumentContentById(ctx context.Context, id bson.ObjectID, doc domain.DocumentContent) error
}

var _ IDocumentContentService = (*DocumentContentService)(nil)

func NewDocumentContentService(repo repository.IDocumentContentRepository) *DocumentContentService {
	return &DocumentContentService{
		repo: repo,
	}
}

type DocumentContentService struct {
	repo repository.IDocumentContentRepository
}

func (s *DocumentContentService) CreateDocumentContent(ctx context.Context, doc domain.DocumentContent) (bson.ObjectID, error) {
	return s.repo.CreateDocumentContent(ctx, doc)
}

func (s *DocumentContentService) FindDocumentContentById(ctx context.Context, id bson.ObjectID) (domain.DocumentContent, error) {
	return s.repo.FindDocumentContentById(ctx, id)
}

func (s *DocumentContentService) DeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	return s.repo.DeleteDocumentContentById(ctx, id)
}

func (s *DocumentContentService) FindDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]domain.DocumentContent, error) {
	return s.repo.FindDocumentContentByParentId(ctx, parentId)
}

func (s *DocumentContentService) FindDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error) {
	return s.repo.FindDocumentContentByDocumentId(ctx, documentId)
}

func (s *DocumentContentService) UpdateDocumentContentById(ctx context.Context, id bson.ObjectID, doc domain.DocumentContent) error {
	return s.repo.UpdateDocumentContentById(ctx, id, doc)
}