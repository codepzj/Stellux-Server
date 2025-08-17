package service

import (
	"context"
	"errors"

	"github.com/codepzj/stellux/server/internal/document/internal/domain"
	"github.com/codepzj/stellux/server/internal/document/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IDocumentService interface {
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

var _ IDocumentService = (*DocumentService)(nil)

func NewDocumentService(repo repository.IDocumentRepository) *DocumentService {
	return &DocumentService{
		repo: repo,
	}
}

type DocumentService struct {
	repo repository.IDocumentRepository
}

func (s *DocumentService) CreateDocument(ctx context.Context, doc *domain.Document) (bson.ObjectID, error) {
	existDoc, err := s.repo.FindDocumentByAlias(ctx, doc.Alias)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return bson.ObjectID{}, err
	}
	if existDoc != nil {
		return bson.ObjectID{}, errors.New("别名已存在")
	}
	return s.repo.CreateDocument(ctx, doc)
}

func (s *DocumentService) FindDocumentById(ctx context.Context, id bson.ObjectID) (*domain.Document, error) {
	return s.repo.FindDocumentById(ctx, id)
}

func (s *DocumentService) UpdateDocumentById(ctx context.Context, id bson.ObjectID, doc *domain.Document) error {
	existDoc, err := s.repo.FindDocumentByAlias(ctx, doc.Alias)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if existDoc != nil && existDoc.Id != id {
		return errors.New("别名已存在")
	}
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

func (s *DocumentService) FindDocumentByAlias(ctx context.Context, alias string) (*domain.Document, error) {
	return s.repo.FindDocumentByAlias(ctx, alias)
}

func (s *DocumentService) GetDocumentList(ctx context.Context, page *domain.Page) ([]*domain.Document, int64, error) {
	return s.repo.GetDocumentList(ctx, page)
}

func (s *DocumentService) GetDocumentBinList(ctx context.Context, page *domain.Page) ([]*domain.Document, int64, error) {
	return s.repo.GetDocumentBinList(ctx, page)
}

func (s *DocumentService) GetPublicDocumentList(ctx context.Context, page *domain.Page) ([]*domain.Document, int64, error) {
	return s.repo.GetPublicDocumentList(ctx, page)
}

func (s *DocumentService) GetAllPublicDocuments(ctx context.Context) ([]*domain.Document, error) {
	return s.repo.GetAllPublicDocuments(ctx)
}
