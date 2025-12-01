package service

import (
	"context"
	"errors"
	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"

	"github.com/codepzj/Stellux-Server/internal/document_content/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/document_content/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/pkg/logger"
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
	GetDocumentContentList(ctx context.Context, page *apiwrap.Page, documentId bson.ObjectID) ([]domain.DocumentContent, int64, error)
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
			logger.Error("查询别名失败",
				logger.WithError(err),
				logger.WithString("alias", doc.Alias),
			)
			return bson.ObjectID{}, err
		}
		if len(docContentList) > 0 {
			logger.Warn("别名已存在",
				logger.WithString("alias", doc.Alias),
			)
			return bson.ObjectID{}, errors.New("别名已存在")
		}
	}

	id, err := s.repo.CreateDocumentContent(ctx, doc)
	if err != nil {
		logger.Error("创建文档内容失败",
			logger.WithError(err),
			logger.WithString("title", doc.Title),
		)
		return bson.ObjectID{}, err
	}

	logger.Info("创建文档内容成功",
		logger.WithString("contentId", id.Hex()),
		logger.WithString("title", doc.Title),
	)

	return id, nil
}

func (s *DocumentContentService) FindDocumentContentById(ctx context.Context, id bson.ObjectID) (domain.DocumentContent, error) {
	logger.Info("查询文档内容",
		logger.WithString("method", "FindDocumentContentById"),
		logger.WithString("contentId", id.Hex()),
	)

	content, err := s.repo.FindDocumentContentById(ctx, id)
	if err != nil {
		logger.Error("查询文档内容失败",
			logger.WithError(err),
			logger.WithString("contentId", id.Hex()),
		)
		return domain.DocumentContent{}, err
	}

	return content, nil
}

func (s *DocumentContentService) DeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	err := s.repo.DeleteDocumentContentById(ctx, id)
	if err != nil {
		logger.Error("删除文档内容失败",
			logger.WithError(err),
			logger.WithString("contentId", id.Hex()),
		)
		return err
	}

	logger.Info("删除文档内容成功",
		logger.WithString("contentId", id.Hex()),
	)

	return nil
}

func (s *DocumentContentService) SoftDeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	err := s.repo.SoftDeleteDocumentContentById(ctx, id)
	if err != nil {
		logger.Error("软删除文档内容失败",
			logger.WithError(err),
			logger.WithString("contentId", id.Hex()),
		)
		return err
	}

	logger.Info("软删除文档内容成功",
		logger.WithString("contentId", id.Hex()),
	)

	return nil
}

func (s *DocumentContentService) RestoreDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	err := s.repo.RestoreDocumentContentById(ctx, id)
	if err != nil {
		logger.Error("恢复文档内容失败",
			logger.WithError(err),
			logger.WithString("contentId", id.Hex()),
		)
		return err
	}

	logger.Info("恢复文档内容成功",
		logger.WithString("contentId", id.Hex()),
	)

	return nil
}

func (s *DocumentContentService) FindDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]domain.DocumentContent, error) {
	logger.Info("根据父节点查询文档内容",
		logger.WithString("method", "FindDocumentContentByParentId"),
		logger.WithString("parentId", parentId.Hex()),
	)

	contents, err := s.repo.FindDocumentContentByParentId(ctx, parentId)
	if err != nil {
		logger.Error("查询文档内容失败",
			logger.WithError(err),
			logger.WithString("parentId", parentId.Hex()),
		)
		return nil, err
	}

	return contents, nil
}

func (s *DocumentContentService) FindDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error) {
	logger.Info("根据文档ID查询所有内容",
		logger.WithString("method", "FindDocumentContentByDocumentId"),
		logger.WithString("documentId", documentId.Hex()),
	)

	contents, err := s.repo.FindDocumentContentByDocumentId(ctx, documentId)
	if err != nil {
		logger.Error("查询文档内容失败",
			logger.WithError(err),
			logger.WithString("documentId", documentId.Hex()),
		)
		return nil, err
	}

	return contents, nil
}

func (s *DocumentContentService) UpdateDocumentContentById(ctx context.Context, id bson.ObjectID, doc domain.DocumentContent) error {
	docContentList, err := s.repo.GetDocumentContentListByAlias(ctx, doc.Alias, doc.DocumentId)
	if err != nil {
		logger.Error("查询别名失败",
			logger.WithError(err),
			logger.WithString("alias", doc.Alias),
		)
		return err
	}

	// 如果别名不存在,或者别名存在且是当前文档的别名,则更新或者当前文档是目录
	if len(docContentList) == 0 || (len(docContentList) == 1 && docContentList[0].Id.Hex() == id.Hex()) || doc.IsDir {
		err = s.repo.UpdateDocumentContentById(ctx, id, doc)
		if err != nil {
			logger.Error("更新文档内容失败",
				logger.WithError(err),
				logger.WithString("contentId", id.Hex()),
			)
			return err
		}
		logger.Info("更新文档内容成功",
			logger.WithString("contentId", id.Hex()),
			logger.WithString("title", doc.Title),
		)
		return nil
	}

	logger.Warn("别名已存在",
		logger.WithString("alias", doc.Alias),
	)
	return errors.New("别名已存在")
}

func (s *DocumentContentService) GetDocumentContentList(ctx context.Context, page *apiwrap.Page, documentId bson.ObjectID) ([]domain.DocumentContent, int64, error) {
	logger.Info("查询文档内容列表",
		logger.WithString("method", "GetDocumentContentList"),
		logger.WithInt("pageNo", int(page.PageNo)),
		logger.WithInt("pageSize", int(page.PageSize)),
	)

	contents, total, err := s.repo.GetDocumentContentList(ctx, page, documentId)
	if err != nil {
		logger.Error("查询文档内容列表失败",
			logger.WithError(err),
		)
		return nil, 0, err
	}

	return contents, total, nil
}

func (s *DocumentContentService) GetPublicDocumentContentListByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error) {
	logger.Info("查询公开文档内容列表",
		logger.WithString("method", "GetPublicDocumentContentListByDocumentId"),
		logger.WithString("documentId", documentId.Hex()),
	)

	contents, err := s.repo.GetPublicDocumentContentListByDocumentId(ctx, documentId)
	if err != nil {
		logger.Error("查询公开文档内容列表失败",
			logger.WithError(err),
			logger.WithString("documentId", documentId.Hex()),
		)
		return nil, err
	}

	return contents, nil
}

func (s *DocumentContentService) SearchDocumentContent(ctx context.Context, keyword string) ([]domain.DocumentContent, error) {
	logger.Info("搜索文档内容",
		logger.WithString("method", "SearchDocumentContent"),
		logger.WithString("keyword", keyword),
	)

	contents, err := s.repo.SearchDocumentContent(ctx, keyword)
	if err != nil {
		logger.Error("搜索文档内容失败",
			logger.WithError(err),
			logger.WithString("keyword", keyword),
		)
		return nil, err
	}

	logger.Info("搜索文档内容成功",
		logger.WithInt("count", len(contents)),
		logger.WithString("keyword", keyword),
	)

	return contents, nil
}

func (s *DocumentContentService) SearchPublicDocumentContent(ctx context.Context, keyword string) ([]domain.DocumentContent, error) {
	logger.Info("搜索公开文档内容",
		logger.WithString("method", "SearchPublicDocumentContent"),
		logger.WithString("keyword", keyword),
	)

	contents, err := s.repo.SearchPublicDocumentContent(ctx, keyword)
	if err != nil {
		logger.Error("搜索公开文档内容失败",
			logger.WithError(err),
			logger.WithString("keyword", keyword),
		)
		return nil, err
	}

	logger.Info("搜索公开文档内容成功",
		logger.WithInt("count", len(contents)),
		logger.WithString("keyword", keyword),
	)

	return contents, nil
}

func (s *DocumentContentService) FindPublicDocumentContentById(ctx context.Context, id bson.ObjectID) (domain.DocumentContent, error) {
	logger.Info("查询公开文档内容",
		logger.WithString("method", "FindPublicDocumentContentById"),
		logger.WithString("contentId", id.Hex()),
	)

	content, err := s.repo.FindPublicDocumentContentById(ctx, id)
	if err != nil {
		logger.Error("查询公开文档内容失败",
			logger.WithError(err),
			logger.WithString("contentId", id.Hex()),
		)
		return domain.DocumentContent{}, err
	}

	return content, nil
}

func (s *DocumentContentService) FindPublicDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]domain.DocumentContent, error) {
	logger.Info("根据父节点查询公开文档内容",
		logger.WithString("method", "FindPublicDocumentContentByParentId"),
		logger.WithString("parentId", parentId.Hex()),
	)

	contents, err := s.repo.FindPublicDocumentContentByParentId(ctx, parentId)
	if err != nil {
		logger.Error("查询公开文档内容失败",
			logger.WithError(err),
			logger.WithString("parentId", parentId.Hex()),
		)
		return nil, err
	}

	return contents, nil
}

func (s *DocumentContentService) FindPublicDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error) {
	logger.Info("根据文档ID查询公开内容",
		logger.WithString("method", "FindPublicDocumentContentByDocumentId"),
		logger.WithString("documentId", documentId.Hex()),
	)

	contents, err := s.repo.FindPublicDocumentContentByDocumentId(ctx, documentId)
	if err != nil {
		logger.Error("查询公开文档内容失败",
			logger.WithError(err),
			logger.WithString("documentId", documentId.Hex()),
		)
		return nil, err
	}

	return contents, nil
}

func (s *DocumentContentService) DeleteDocumentContentList(ctx context.Context, ids []string) error {
	err := s.repo.DeleteDocumentContentList(ctx, ids)
	if err != nil {
		logger.Error("批量删除文档内容失败",
			logger.WithError(err),
			logger.WithInt("count", len(ids)),
		)
		return err
	}

	logger.Info("批量删除文档内容成功",
		logger.WithInt("count", len(ids)),
	)

	return nil
}

func (s *DocumentContentService) FindPublicDocumentContentByRootIdAndAlias(ctx context.Context, documentId bson.ObjectID, alias string) (domain.DocumentContent, error) {
	logger.Info("根据别名查询公开文档内容",
		logger.WithString("method", "FindPublicDocumentContentByRootIdAndAlias"),
		logger.WithString("documentId", documentId.Hex()),
		logger.WithString("alias", alias),
	)

	content, err := s.repo.FindPublicDocumentContentByRootIdAndAlias(ctx, documentId, alias)
	if err != nil {
		logger.Error("查询公开文档内容失败",
			logger.WithError(err),
			logger.WithString("alias", alias),
		)
		return domain.DocumentContent{}, err
	}

	return content, nil
}
