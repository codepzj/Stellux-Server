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
	SoftDeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error
	RestoreDocumentContentById(ctx context.Context, id bson.ObjectID) error
	FindDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]domain.DocumentContent, error)
	FindDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error)
	UpdateDocumentContentById(ctx context.Context, id bson.ObjectID, doc domain.DocumentContent) error
	GetDocumentContentList(ctx context.Context, page *domain.Page) ([]domain.DocumentContent, int64, error)
	GetPublicDocumentContentList(ctx context.Context, page *domain.Page) ([]domain.DocumentContent, int64, error)
	SearchDocumentContent(ctx context.Context, keyword string) ([]domain.DocumentContent, error)
	SearchPublicDocumentContent(ctx context.Context, keyword string) ([]domain.DocumentContent, error)
	FindPublicDocumentContentById(ctx context.Context, id bson.ObjectID) (domain.DocumentContent, error)
	FindPublicDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]domain.DocumentContent, error)
	FindPublicDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error)
	UpdateLikeCount(ctx context.Context, id bson.ObjectID) error
	UpdateDislikeCount(ctx context.Context, id bson.ObjectID) error
	UpdateCommentCount(ctx context.Context, id bson.ObjectID) error
	DeleteDocumentContentList(ctx context.Context, ids []string) error
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
		Alias:        doc.Alias,
		ParentId:     doc.ParentId,
		IsDir:        doc.IsDir,
		Sort:         doc.Sort,
		LikeCount:    doc.LikeCount,
		DislikeCount: doc.DislikeCount,
		CommentCount: doc.CommentCount,
		IsDeleted:    doc.IsDeleted,
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
		DeletedAt:    doc.DeletedAt,
		DocumentId:   doc.DocumentId,
		Title:        doc.Title,
		Content:      doc.Content,
		Description:  doc.Description,
		Alias:        doc.Alias,
		ParentId:     doc.ParentId,
		IsDir:        doc.IsDir,
		Sort:         doc.Sort,
		LikeCount:    doc.LikeCount,
		DislikeCount: doc.DislikeCount,
		CommentCount: doc.CommentCount,
		IsDeleted:    doc.IsDeleted,
	}, nil
}

// DeleteDocumentContentById 根据id删除文档内容
func (r *DocumentContentRepository) DeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	return r.dao.DeleteDocumentContentById(ctx, id)
}

// SoftDeleteDocumentContentById 根据id软删除文档内容
func (r *DocumentContentRepository) SoftDeleteDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	return r.dao.SoftDeleteDocumentContentById(ctx, id)
}

// RestoreDocumentContentById 根据id恢复文档内容
func (r *DocumentContentRepository) RestoreDocumentContentById(ctx context.Context, id bson.ObjectID) error {
	return r.dao.RestoreDocumentContentById(ctx, id)
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
			DeletedAt:    doc.DeletedAt,
			DocumentId:   doc.DocumentId,
			Title:        doc.Title,
			Content:      doc.Content,
			Description:  doc.Description,
			Alias:        doc.Alias,
			ParentId:     doc.ParentId,
			IsDir:        doc.IsDir,
			Sort:         doc.Sort,
			LikeCount:    doc.LikeCount,
			DislikeCount: doc.DislikeCount,
			CommentCount: doc.CommentCount,
			IsDeleted:    doc.IsDeleted,
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
			DeletedAt:    doc.DeletedAt,
			DocumentId:   doc.DocumentId,
			Title:        doc.Title,
			Content:      doc.Content,
			Description:  doc.Description,
			Alias:        doc.Alias,
			ParentId:     doc.ParentId,
			IsDir:        doc.IsDir,
			Sort:         doc.Sort,
			LikeCount:    doc.LikeCount,
			DislikeCount: doc.DislikeCount,
			CommentCount: doc.CommentCount,
			IsDeleted:    doc.IsDeleted,
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
		Alias:        doc.Alias,
		ParentId:     doc.ParentId,
		IsDir:        doc.IsDir,
		Sort:         doc.Sort,
		LikeCount:    doc.LikeCount,
		DislikeCount: doc.DislikeCount,
		CommentCount: doc.CommentCount,
		IsDeleted:    doc.IsDeleted,
	})
}

// GetDocumentContentList 获取文档内容列表
func (r *DocumentContentRepository) GetDocumentContentList(ctx context.Context, page *domain.Page) ([]domain.DocumentContent, int64, error) {
	docs, count, err := r.dao.GetDocumentContentList(ctx, &dao.Page{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	})
	if err != nil {
		return nil, 0, err
	}

	results := make([]domain.DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = domain.DocumentContent{
			Id:           doc.ID,
			CreatedAt:    doc.CreatedAt,
			UpdatedAt:    doc.UpdatedAt,
			DeletedAt:    doc.DeletedAt,
			DocumentId:   doc.DocumentId,
			Title:        doc.Title,
			Content:      doc.Content,
			Description:  doc.Description,
			Alias:        doc.Alias,
			ParentId:     doc.ParentId,
			IsDir:        doc.IsDir,
			Sort:         doc.Sort,
			LikeCount:    doc.LikeCount,
			DislikeCount: doc.DislikeCount,
			CommentCount: doc.CommentCount,
			IsDeleted:    doc.IsDeleted,
		}
	}
	return results, count, nil
}

// GetPublicDocumentContentList 获取公开文档内容列表
func (r *DocumentContentRepository) GetPublicDocumentContentList(ctx context.Context, page *domain.Page) ([]domain.DocumentContent, int64, error) {
	docs, count, err := r.dao.GetPublicDocumentContentList(ctx, &dao.Page{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	})
	if err != nil {
		return nil, 0, err
	}

	results := make([]domain.DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = domain.DocumentContent{
			Id:           doc.ID,
			CreatedAt:    doc.CreatedAt,
			UpdatedAt:    doc.UpdatedAt,
			DeletedAt:    doc.DeletedAt,
			DocumentId:   doc.DocumentId,
			Title:        doc.Title,
			Content:      doc.Content,
			Description:  doc.Description,
			Alias:        doc.Alias,
			ParentId:     doc.ParentId,
			IsDir:        doc.IsDir,
			Sort:         doc.Sort,
			LikeCount:    doc.LikeCount,
			DislikeCount: doc.DislikeCount,
			CommentCount: doc.CommentCount,
			IsDeleted:    doc.IsDeleted,
		}
	}
	return results, count, nil
}

// SearchDocumentContent 搜索文档内容
func (r *DocumentContentRepository) SearchDocumentContent(ctx context.Context, keyword string) ([]domain.DocumentContent, error) {
	docs, err := r.dao.SearchDocumentContent(ctx, keyword)
	if err != nil {
		return nil, err
	}

	results := make([]domain.DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = domain.DocumentContent{
			Id:           doc.ID,
			CreatedAt:    doc.CreatedAt,
			UpdatedAt:    doc.UpdatedAt,
			DeletedAt:    doc.DeletedAt,
			DocumentId:   doc.DocumentId,
			Title:        doc.Title,
			Content:      doc.Content,
			Description:  doc.Description,
			Alias:        doc.Alias,
			ParentId:     doc.ParentId,
			IsDir:        doc.IsDir,
			Sort:         doc.Sort,
			LikeCount:    doc.LikeCount,
			DislikeCount: doc.DislikeCount,
			CommentCount: doc.CommentCount,
			IsDeleted:    doc.IsDeleted,
		}
	}
	return results, nil
}

// SearchPublicDocumentContent 搜索公开文档内容
func (r *DocumentContentRepository) SearchPublicDocumentContent(ctx context.Context, keyword string) ([]domain.DocumentContent, error) {
	docs, err := r.dao.SearchPublicDocumentContent(ctx, keyword)
	if err != nil {
		return nil, err
	}

	results := make([]domain.DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = domain.DocumentContent{
			Id:           doc.ID,
			CreatedAt:    doc.CreatedAt,
			UpdatedAt:    doc.UpdatedAt,
			DeletedAt:    doc.DeletedAt,
			DocumentId:   doc.DocumentId,
			Title:        doc.Title,
			Content:      doc.Content,
			Description:  doc.Description,
			Alias:        doc.Alias,
			ParentId:     doc.ParentId,
			IsDir:        doc.IsDir,
			Sort:         doc.Sort,
			LikeCount:    doc.LikeCount,
			DislikeCount: doc.DislikeCount,
			CommentCount: doc.CommentCount,
			IsDeleted:    doc.IsDeleted,
		}
	}
	return results, nil
}

// FindPublicDocumentContentById 根据id查询公开文档内容
func (r *DocumentContentRepository) FindPublicDocumentContentById(ctx context.Context, id bson.ObjectID) (domain.DocumentContent, error) {
	doc, err := r.dao.FindPublicDocumentContentById(ctx, id)
	if err != nil {
		return domain.DocumentContent{}, err
	}

	return domain.DocumentContent{
		Id:           doc.ID,
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
		DeletedAt:    doc.DeletedAt,
		DocumentId:   doc.DocumentId,
		Title:        doc.Title,
		Content:      doc.Content,
		Description:  doc.Description,
		Alias:        doc.Alias,
		ParentId:     doc.ParentId,
		IsDir:        doc.IsDir,
		Sort:         doc.Sort,
		LikeCount:    doc.LikeCount,
		DislikeCount: doc.DislikeCount,
		CommentCount: doc.CommentCount,
		IsDeleted:    doc.IsDeleted,
	}, nil
}

// FindPublicDocumentContentByParentId 根据父级ID查询公开文档内容
func (r *DocumentContentRepository) FindPublicDocumentContentByParentId(ctx context.Context, parentId bson.ObjectID) ([]domain.DocumentContent, error) {
	docs, err := r.dao.FindPublicDocumentContentByParentId(ctx, parentId)
	if err != nil {
		return nil, err
	}
	results := make([]domain.DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = domain.DocumentContent{
			Id:           doc.ID,
			CreatedAt:    doc.CreatedAt,
			UpdatedAt:    doc.UpdatedAt,
			DeletedAt:    doc.DeletedAt,
			DocumentId:   doc.DocumentId,
			Title:        doc.Title,
			Content:      doc.Content,
			Description:  doc.Description,
			Alias:        doc.Alias,
			ParentId:     doc.ParentId,
			IsDir:        doc.IsDir,
			Sort:         doc.Sort,
			LikeCount:    doc.LikeCount,
			DislikeCount: doc.DislikeCount,
			CommentCount: doc.CommentCount,
			IsDeleted:    doc.IsDeleted,
		}
	}
	return results, nil
}

// FindPublicDocumentContentByDocumentId 根据文档ID查询公开文档内容
func (r *DocumentContentRepository) FindPublicDocumentContentByDocumentId(ctx context.Context, documentId bson.ObjectID) ([]domain.DocumentContent, error) {
	docs, err := r.dao.FindPublicDocumentContentByDocumentId(ctx, documentId)
	if err != nil {
		return nil, err
	}
	results := make([]domain.DocumentContent, len(docs))
	for i, doc := range docs {
		results[i] = domain.DocumentContent{
			Id:           doc.ID,
			CreatedAt:    doc.CreatedAt,
			UpdatedAt:    doc.UpdatedAt,
			DeletedAt:    doc.DeletedAt,
			DocumentId:   doc.DocumentId,
			Title:        doc.Title,
			Content:      doc.Content,
			Description:  doc.Description,
			Alias:        doc.Alias,
			ParentId:     doc.ParentId,
			IsDir:        doc.IsDir,
			Sort:         doc.Sort,
			LikeCount:    doc.LikeCount,
			DislikeCount: doc.DislikeCount,
			CommentCount: doc.CommentCount,
			IsDeleted:    doc.IsDeleted,
		}
	}
	return results, nil
}

// UpdateLikeCount 更新点赞数
func (r *DocumentContentRepository) UpdateLikeCount(ctx context.Context, id bson.ObjectID) error {
	return r.dao.UpdateLikeCount(ctx, id)
}

// UpdateDislikeCount 更新反对数
func (r *DocumentContentRepository) UpdateDislikeCount(ctx context.Context, id bson.ObjectID) error {
	return r.dao.UpdateDislikeCount(ctx, id)
}

// UpdateCommentCount 更新评论数
func (r *DocumentContentRepository) UpdateCommentCount(ctx context.Context, id bson.ObjectID) error {
	return r.dao.UpdateCommentCount(ctx, id)
}

func (r *DocumentContentRepository) DeleteDocumentContentList(ctx context.Context, ids []string) error {
	return r.dao.DeleteDocumentContentList(ctx, ids)
}
