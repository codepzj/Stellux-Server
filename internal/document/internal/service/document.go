package service

import (
	"context"

	"github.com/codepzj/stellux/server/internal/document/internal/domain"
	"github.com/codepzj/stellux/server/internal/document/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IDocumentService interface {
	CreateDocument(ctx context.Context, doc domain.Document) (bson.ObjectID, error)
	FindDocumentById(ctx context.Context, id bson.ObjectID) (domain.Document, error)
	UpdateDocumentById(ctx context.Context, id bson.ObjectID, doc domain.Document) error
	DeleteDocumentById(ctx context.Context, id bson.ObjectID) error
	SoftDeleteDocumentById(ctx context.Context, id bson.ObjectID) error
	RestoreDocumentById(ctx context.Context, id bson.ObjectID) error
	FindDocumentByAlias(ctx context.Context, alias string, filter bson.D) (domain.Document, error)
	GetDocumentList(ctx context.Context, page *domain.Page) ([]domain.Document, int64, error)
	GetPublicDocumentList(ctx context.Context, page *domain.Page) ([]domain.Document, int64, error)
	GetAllPublicDocuments(ctx context.Context) ([]domain.Document, error)
}

var _ IDocumentService = (*DocumentService)(nil)

func NewDocumentService(repo repository.IDocumentRepository) *DocumentService {
	return &DocumentService{
		repo: repo,
	}
}

type DocumentService struct {
	repo repository.IDocumentRepository
}

func (s *DocumentService) CreateDocument(ctx context.Context, doc domain.Document) (bson.ObjectID, error) {
	return s.repo.CreateDocument(ctx, doc)
}

func (s *DocumentService) FindDocumentById(ctx context.Context, id bson.ObjectID) (domain.Document, error) {
	return s.repo.FindDocumentById(ctx, id)
}

func (s *DocumentService) UpdateDocumentById(ctx context.Context, id bson.ObjectID, doc domain.Document) error {
	return s.repo.UpdateDocumentById(ctx, id, doc)
}

func (s *DocumentService) DeleteDocumentById(ctx context.Context, id bson.ObjectID) error {
	return s.repo.DeleteDocumentById(ctx, id)
}

func (s *DocumentService) SoftDeleteDocumentById(ctx context.Context, id bson.ObjectID) error {
	return s.repo.SoftDeleteDocumentById(ctx, id)
}

func (s *DocumentService) RestoreDocumentById(ctx context.Context, id bson.ObjectID) error {
	return s.repo.RestoreDocumentById(ctx, id)
}

func (s *DocumentService) FindDocumentByAlias(ctx context.Context, alias string, filter bson.D) (domain.Document, error) {
	return s.repo.FindDocumentByAlias(ctx, alias, filter)
}

func (s *DocumentService) GetDocumentList(ctx context.Context, page *domain.Page) ([]domain.Document, int64, error) {
	return s.repo.GetDocumentList(ctx, page)
}

func (s *DocumentService) GetPublicDocumentList(ctx context.Context, page *domain.Page) ([]domain.Document, int64, error) {
	return s.repo.GetPublicDocumentList(ctx, page)
}

func (s *DocumentService) GetAllPublicDocuments(ctx context.Context) ([]domain.Document, error) {
	return s.repo.GetAllPublicDocuments(ctx)
}
