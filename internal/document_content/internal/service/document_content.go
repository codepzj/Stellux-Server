package service

import (
	"context"

	"github.com/codepzj/stellux/server/internal/document_content/internal/domain"
	"github.com/codepzj/stellux/server/internal/document_content/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)


type IDocumentContentService interface {
	CreateDocumentContent(ctx context.Context, doc *domain.DocumentContent) error
	FindDocumentContentByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.DocumentContent, error)
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

func (s *DocumentContentService) CreateDocumentContent(ctx context.Context, doc *domain.DocumentContent) error {
	return s.repo.CreateDocumentContent(ctx, doc)
}

func (s *DocumentContentService) FindDocumentContentByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.DocumentContent, error) {
	return s.repo.FindDocumentContentByDocumentID(ctx, documentID)
}