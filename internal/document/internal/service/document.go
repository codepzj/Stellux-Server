package service

import (
	"context"
	"fmt"

	"github.com/codepzj/stellux/server/global"
	"github.com/codepzj/stellux/server/internal/document/internal/domain"
	"github.com/codepzj/stellux/server/internal/document/internal/repository"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IDocumentService interface {
	FindAllPublic(ctx context.Context) ([]*domain.DocumentRoot, error)
	FindAllPublicByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error)
	FindByKeyword(ctx context.Context, keyword string, documentID bson.ObjectID) ([]*domain.Document, error)
	GetDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.Document, error)
	GetRootDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.DocumentRoot, error)
	GenerateSitemap(ctx context.Context) ([]*domain.DocumentSitemap, error)
	AdminCreate(ctx context.Context, doc *domain.Document) error
	AdminCreateRoot(ctx context.Context, doc *domain.DocumentRoot) error
	AdminFindAllRoot(ctx context.Context) ([]*domain.DocumentRoot, error)
	AdminFindAllParent(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error)
	AdminFindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error)
	AdminGetDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.Document, error)
	AdminGetRootDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.DocumentRoot, error)
	AdminUpdateDocumentByID(ctx context.Context, id bson.ObjectID, title string, content string) error
	AdminDeleteByID(ctx context.Context, id bson.ObjectID) error
	AdminDeleteByIDList(ctx context.Context, idList []bson.ObjectID) error
	AdminRenameDocumentByID(ctx context.Context, id bson.ObjectID, title string) error
	AdminEditRootDocumentByID(ctx context.Context, id bson.ObjectID, doc *domain.DocumentRoot) error
	AdminDeleteRootDocumentByID(ctx context.Context, id bson.ObjectID) error
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

// 获取所有公共根文档(用于文档概览)
func (s *DocumentService) FindAllPublic(ctx context.Context) ([]*domain.DocumentRoot, error) {
	return s.repo.FindAllPublicRootDocument(ctx)
}

// 根据文档id获取所有公共子文档
func (s *DocumentService) FindAllPublicByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error) {
	documentIsPublic, err := s.repo.FindDocumentIsPublic(ctx, documentID)
	if err != nil {
		return nil, err
	}
	if !documentIsPublic {
		return nil, global.ErrDocumentNotPublic
	}
	return s.repo.FindAllByDocumentID(ctx, documentID)
}

// 根据关键词查询文档
func (s *DocumentService) FindByKeyword(ctx context.Context, keyword string, documentID bson.ObjectID) ([]*domain.Document, error) {
	// TODO: 获取根文档的权限
	return s.repo.FindByKeyword(ctx, keyword, documentID)
}

// 根据id获取根文档
func (s *DocumentService) GetRootDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.DocumentRoot, error) {
	//TODO: 获取根文档的权限
	return s.repo.GetRootDocumentByID(ctx, id)
}

// 根据id获取文档
func (s *DocumentService) GetDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.Document, error) {
	// TODO: 获取文档的权限
	return s.repo.GetDocumentByID(ctx, id)
}

// 生成站点地图
func (s *DocumentService) GenerateSitemap(ctx context.Context) ([]*domain.DocumentSitemap, error) {
	rootDocumentList, err := s.repo.FindAllPublicRootDocument(ctx)
	if err != nil {
		return nil, err
	}
	rootIdList := lo.Map(rootDocumentList, func(item *domain.DocumentRoot, _ int) bson.ObjectID {
		return item.ID
	})
	childDocumentList, err := s.repo.FindAllByDocumentIDList(ctx, rootIdList)
	if err != nil {
		return nil, err
	}
	sitemapList := DocumentListToSitemap(rootDocumentList, childDocumentList)
	return sitemapList, nil
}

// 新增文档
func (s *DocumentService) AdminCreate(ctx context.Context, doc *domain.Document) error {
	return s.repo.Create(ctx, doc)
}

// 新增根文档
func (s *DocumentService) AdminCreateRoot(ctx context.Context, doc *domain.DocumentRoot) error {
	return s.repo.CreateRoot(ctx, doc)
}

// 管理员编辑根文档
func (s *DocumentService) AdminEditRootDocumentByID(ctx context.Context, id bson.ObjectID, doc *domain.DocumentRoot) error {
	return s.repo.UpdateRootDocumentByID(ctx, id, doc)
}

// 管理员删除根文档
func (s *DocumentService) AdminDeleteRootDocumentByID(ctx context.Context, id bson.ObjectID) error {
	return s.repo.DeleteRootDocumentByID(ctx, id)
}

// 管理员获取所有根文档
func (s *DocumentService) AdminFindAllRoot(ctx context.Context) ([]*domain.DocumentRoot, error) {
	return s.repo.FindAllRootDocument(ctx)
}

// 管理员获取所有父文档
func (s *DocumentService) AdminFindAllParent(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error) {
	return s.repo.FindAllByTypeAndDocumentID(ctx, "parent", documentID)
}

// 管理员获取一个文档的所有子文档(包含非直接子文档)
func (s *DocumentService) AdminFindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error) {
	return s.repo.FindAllByDocumentID(ctx, documentID)
}

// 管理员根据id获取文档
func (s *DocumentService) AdminGetDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.Document, error) {
	return s.repo.GetDocumentByID(ctx, id)
}

// 管理员根据id获取根文档
func (s *DocumentService) AdminGetRootDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.DocumentRoot, error) {
	return s.repo.GetRootDocumentByID(ctx, id)
}

// 管理员更新文档
func (s *DocumentService) AdminUpdateDocumentByID(ctx context.Context, id bson.ObjectID, title string, content string) error {
	return s.repo.UpdateDocumentByID(ctx, id, title, content)
}

// 管理员重命名文档
func (s *DocumentService) AdminRenameDocumentByID(ctx context.Context, id bson.ObjectID, title string) error {
	return s.repo.RenameDocumentByID(ctx, id, title)
}

// 管理员删除文档
func (s *DocumentService) AdminDeleteByID(ctx context.Context, id bson.ObjectID) error {
	return s.repo.DeleteByID(ctx, id)
}

// 管理员删除多个文档
func (s *DocumentService) AdminDeleteByIDList(ctx context.Context, idList []bson.ObjectID) error {
	return s.repo.DeleteByIDList(ctx, idList)
}

func DocumentListToSitemap(
	rootDocumentList []*domain.DocumentRoot,
	childDocumentList []*domain.Document,
) []*domain.DocumentSitemap {
	// 构建 root ID -> Alias 映射，避免在循环中重复查找
	rootAliasMap := make(map[bson.ObjectID]string, len(rootDocumentList))
	rootSitemapList := make([]*domain.DocumentSitemap, 0, len(rootDocumentList))

	for _, root := range rootDocumentList {
		rootAliasMap[root.ID] = root.Alias
		rootSitemapList = append(rootSitemapList, &domain.DocumentSitemap{
			ID:           root.ID,
			UpdatedAt:    root.UpdatedAt,
			Alias:        root.Alias,
			DocumentType: root.DocumentType,
			ParentID:     bson.ObjectID{},
			DocumentID:   bson.ObjectID{},
		})
	}

	childSitemapList := make([]*domain.DocumentSitemap, 0, len(childDocumentList))
	for _, doc := range childDocumentList {
		alias := rootAliasMap[doc.DocumentID] // 如果不存在，alias 会为空字符串
		fmt.Println(alias, doc.DocumentID)
		childSitemapList = append(childSitemapList, &domain.DocumentSitemap{
			ID:           doc.ID,
			UpdatedAt:    doc.UpdatedAt,
			Alias:        alias,
			DocumentType: doc.DocumentType,
			ParentID:     doc.ParentID,
			DocumentID:   doc.DocumentID,
		})
	}

	return append(rootSitemapList, childSitemapList...)
}
