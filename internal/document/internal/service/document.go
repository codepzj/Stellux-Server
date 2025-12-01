package service

import (
	"context"
	"errors"
	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"

	"github.com/codepzj/Stellux-Server/internal/document/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/document/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/pkg/logger"
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
	GetDocumentList(ctx context.Context, page *apiwrap.Page) ([]*domain.Document, int64, error)
	GetDocumentBinList(ctx context.Context, page *apiwrap.Page) ([]*domain.Document, int64, error)
	GetPublicDocumentList(ctx context.Context, page *apiwrap.Page) ([]*domain.Document, int64, error)
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
		logger.Error("查询别名失败",
			logger.WithError(err),
			logger.WithString("alias", doc.Alias),
		)
		return bson.ObjectID{}, err
	}

	if existDoc != nil {
		logger.Warn("文档别名已存在",
			logger.WithString("alias", doc.Alias),
			logger.WithString("existDocId", existDoc.Id.Hex()),
		)
		return bson.ObjectID{}, errors.New("别名已存在")
	}

	id, err := s.repo.CreateDocument(ctx, doc)
	if err != nil {
		logger.Error("创建文档失败",
			logger.WithError(err),
			logger.WithString("title", doc.Title),
		)
		return bson.ObjectID{}, err
	}

	logger.Info("创建文档成功",
		logger.WithString("documentId", id.Hex()),
		logger.WithString("title", doc.Title),
	)

	return id, nil
}

func (s *DocumentService) FindDocumentById(ctx context.Context, id bson.ObjectID) (*domain.Document, error) {
	logger.Info("查询文档",
		logger.WithString("method", "FindDocumentById"),
		logger.WithString("documentId", id.Hex()),
	)

	doc, err := s.repo.FindDocumentById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Warn("文档不存在",
				logger.WithString("documentId", id.Hex()),
			)
		} else {
			logger.Error("查询文档失败",
				logger.WithError(err),
				logger.WithString("documentId", id.Hex()),
			)
		}
		return nil, err
	}

	return doc, nil
}

func (s *DocumentService) UpdateDocumentById(ctx context.Context, id bson.ObjectID, doc *domain.Document) error {
	existDoc, err := s.repo.FindDocumentByAlias(ctx, doc.Alias)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		logger.Error("查询别名失败",
			logger.WithError(err),
			logger.WithString("alias", doc.Alias),
		)
		return err
	}

	if existDoc != nil && existDoc.Id != id {
		logger.Warn("别名已被其他文档使用",
			logger.WithString("alias", doc.Alias),
			logger.WithString("existDocId", existDoc.Id.Hex()),
		)
		return errors.New("别名已存在")
	}

	err = s.repo.UpdateDocumentById(ctx, id, doc)
	if err != nil {
		logger.Error("更新文档失败",
			logger.WithError(err),
			logger.WithString("documentId", id.Hex()),
		)
		return err
	}

	logger.Info("更新文档成功",
		logger.WithString("documentId", id.Hex()),
		logger.WithString("title", doc.Title),
	)

	return nil
}

func (s *DocumentService) DeleteDocumentById(ctx context.Context, id bson.ObjectID) error {
	err := s.repo.DeleteDocumentById(ctx, id)
	if err != nil {
		logger.Error("删除文档失败",
			logger.WithError(err),
			logger.WithString("documentId", id.Hex()),
		)
		return err
	}

	logger.Info("删除文档成功",
		logger.WithString("documentId", id.Hex()),
	)

	return nil
}

func (s *DocumentService) SoftDeleteDocumentById(ctx context.Context, id bson.ObjectID) error {
	err := s.repo.SoftDeleteDocumentById(ctx, id)
	if err != nil {
		logger.Error("软删除文档失败",
			logger.WithError(err),
			logger.WithString("documentId", id.Hex()),
		)
		return err
	}

	logger.Info("软删除文档成功",
		logger.WithString("documentId", id.Hex()),
	)

	return nil
}

func (s *DocumentService) RestoreDocumentById(ctx context.Context, id bson.ObjectID) error {
	err := s.repo.RestoreDocumentById(ctx, id)
	if err != nil {
		logger.Error("恢复文档失败",
			logger.WithError(err),
			logger.WithString("documentId", id.Hex()),
		)
		return err
	}

	logger.Info("恢复文档成功",
		logger.WithString("documentId", id.Hex()),
	)

	return nil
}

func (s *DocumentService) FindDocumentByAlias(ctx context.Context, alias string) (*domain.Document, error) {
	logger.Info("根据别名查询文档",
		logger.WithString("method", "FindDocumentByAlias"),
		logger.WithString("alias", alias),
	)

	doc, err := s.repo.FindDocumentByAlias(ctx, alias)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Warn("文档不存在",
				logger.WithString("alias", alias),
			)
		} else {
			logger.Error("查询文档失败",
				logger.WithError(err),
				logger.WithString("alias", alias),
			)
		}
		return nil, err
	}

	return doc, nil
}

func (s *DocumentService) GetDocumentList(ctx context.Context, page *apiwrap.Page) ([]*domain.Document, int64, error) {
	logger.Info("查询文档列表",
		logger.WithString("method", "GetDocumentList"),
		logger.WithInt("pageNo", int(page.PageNo)),
		logger.WithInt("pageSize", int(page.PageSize)),
	)

	docs, total, err := s.repo.GetDocumentList(ctx, page)
	if err != nil {
		logger.Error("查询文档列表失败",
			logger.WithError(err),
		)
		return nil, 0, err
	}

	return docs, total, nil
}

func (s *DocumentService) GetDocumentBinList(ctx context.Context, page *apiwrap.Page) ([]*domain.Document, int64, error) {
	logger.Info("查询回收站文档列表",
		logger.WithString("method", "GetDocumentBinList"),
		logger.WithInt("pageNo", int(page.PageNo)),
		logger.WithInt("pageSize", int(page.PageSize)),
	)

	docs, total, err := s.repo.GetDocumentBinList(ctx, page)
	if err != nil {
		logger.Error("查询回收站列表失败",
			logger.WithError(err),
		)
		return nil, 0, err
	}

	return docs, total, nil
}

func (s *DocumentService) GetPublicDocumentList(ctx context.Context, page *apiwrap.Page) ([]*domain.Document, int64, error) {
	logger.Info("查询公开文档列表",
		logger.WithString("method", "GetPublicDocumentList"),
		logger.WithInt("pageNo", int(page.PageNo)),
		logger.WithInt("pageSize", int(page.PageSize)),
	)

	docs, total, err := s.repo.GetPublicDocumentList(ctx, page)
	if err != nil {
		logger.Error("查询公开文档列表失败",
			logger.WithError(err),
		)
		return nil, 0, err
	}

	return docs, total, nil
}

func (s *DocumentService) GetAllPublicDocuments(ctx context.Context) ([]*domain.Document, error) {
	logger.Info("查询所有公开文档",
		logger.WithString("method", "GetAllPublicDocuments"),
	)

	docs, err := s.repo.GetAllPublicDocuments(ctx)
	if err != nil {
		logger.Error("查询所有公开文档失败",
			logger.WithError(err),
		)
		return nil, err
	}

	return docs, nil
}
