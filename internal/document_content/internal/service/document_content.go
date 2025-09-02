package service

import (
	"context"
	"errors"

	"github.com/codepzj/Stellux-Server/internal/document_content/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/document_content/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IDocumentContentService interface {
	CreateDocumentContent(ctx context.Context, doc domain.DocumentContent) (bson.ObjectID, error)
	FindDocumentContentById(ctx context.Context, id bson.ObjectID) (domain.DocumentContent, error)
	DeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error
	SoftDeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error
	RestoreDocumentContentById(ctx context.Context, id bson.ObjectID) error
	FindDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]domain.DocumentContent, error)
	FindDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error)
	UpdateDocumentContentById(ctx context.Context, id bson.ObjectID, doc domain.DocumentContent) error
	GetDocumentContentList(ctx context.Context, page *domain.Page) ([]domain.DocumentContent, int64, error)
	GetPublicDocumentContentListByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error)
	SearchDocumentContent(ctx context.Context, keyword string) ([]domain.DocumentContent, error)
	SearchPublicDocumentContent(ctx context.Context, keyword string) ([]domain.DocumentContent, error)
	FindPublicDocumentContentById(ctx context.Context, id bson.ObjectID) (domain.DocumentContent, error)
	FindPublicDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]domain.DocumentContent, error)
	FindPublicDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error)
	FindPublicDocumentContentByRootIdAndAlias(ctx context.Context, documentId bson.ObjectID, alias string) (domain.DocumentContent, error)
	DeleteDocumentContentList(ctx context.Context, ids []string) error
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
	// 如果非目录，则检查别名是否存在
	if !doc.IsDir {
		docContentList, err := s.repo.GetDocumentContentListByAlias(ctx, doc.Alias, doc.DocumentId)
		if err != nil {
			return bson.ObjectID{}, err
		}
		if len(docContentList) > 0 {
			return bson.ObjectID{}, errors.New("别名已存在")
		}
	}
	return s.repo.CreateDocumentContent(ctx, doc)
}

func (s *DocumentContentService) FindDocumentContentById(ctx context.Context, id bson.ObjectID) (domain.DocumentContent, error) {
	return s.repo.FindDocumentContentById(ctx, id)
}

func (s *DocumentContentService) DeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	return s.repo.DeleteDocumentContentById(ctx, id)
}

func (s *DocumentContentService) SoftDeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	return s.repo.SoftDeleteDocumentContentById(ctx, id)
}

func (s *DocumentContentService) RestoreDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	return s.repo.RestoreDocumentContentById(ctx, id)
}

func (s *DocumentContentService) FindDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]domain.DocumentContent, error) {
	return s.repo.FindDocumentContentByParentId(ctx, parentId)
}

func (s *DocumentContentService) FindDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error) {
	return s.repo.FindDocumentContentByDocumentId(ctx, documentId)
}

func (s *DocumentContentService) UpdateDocumentContentById(ctx context.Context, id bson.ObjectID, doc domain.DocumentContent) error {
	docContentList, err := s.repo.GetDocumentContentListByAlias(ctx, doc.Alias, doc.DocumentId)
	if err != nil {
		return err
	}
	// 如果别名不存在,或者别名存在且是当前文档的别名,则更新或者当前文档是目录
	if len(docContentList) == 0 || (len(docContentList) == 1 && docContentList[0].Id.Hex() == id.Hex()) || doc.IsDir {
		return s.repo.UpdateDocumentContentById(ctx, id, doc)
	}
	return errors.New("别名已存在")
}

func (s *DocumentContentService) GetDocumentContentList(ctx context.Context, page *domain.Page) ([]domain.DocumentContent, int64, error) {
	return s.repo.GetDocumentContentList(ctx, page)
}

func (s *DocumentContentService) GetPublicDocumentContentListByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error) {
	return s.repo.GetPublicDocumentContentListByDocumentId(ctx, documentId)
}

func (s *DocumentContentService) SearchDocumentContent(ctx context.Context, keyword string) ([]domain.DocumentContent, error) {
	return s.repo.SearchDocumentContent(ctx, keyword)
}

func (s *DocumentContentService) SearchPublicDocumentContent(ctx context.Context, keyword string) ([]domain.DocumentContent, error) {
	return s.repo.SearchPublicDocumentContent(ctx, keyword)
}

func (s *DocumentContentService) FindPublicDocumentContentById(ctx context.Context, id bson.ObjectID) (domain.DocumentContent, error) {
	return s.repo.FindPublicDocumentContentById(ctx, id)
}

func (s *DocumentContentService) FindPublicDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]domain.DocumentContent, error) {
	return s.repo.FindPublicDocumentContentByParentId(ctx, parentId)
}

func (s *DocumentContentService) FindPublicDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error) {
	return s.repo.FindPublicDocumentContentByDocumentId(ctx, documentId)
}

func (s *DocumentContentService) DeleteDocumentContentList(ctx context.Context, ids []string) error {
	return s.repo.DeleteDocumentContentList(ctx, ids)
}

func (s *DocumentContentService) FindPublicDocumentContentByRootIdAndAlias(ctx context.Context, documentId bson.ObjectID, alias string) (domain.DocumentContent, error) {
	return s.repo.FindPublicDocumentContentByRootIdAndAlias(ctx, documentId, alias)
}
