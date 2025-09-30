package repository

import (
	"context"

	"github.com/codepzj/Stellux-Server/internal/post/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/post/internal/repository/dao"
	"github.com/codepzj/gokit/slice"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type IPostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	Update(ctx context.Context, post *domain.Post) error
	UpdatePublishStatus(ctx context.Context, id uint, isPublish bool) error
	SoftDelete(ctx context.Context, id uint) error
	SoftDeleteBatch(ctx context.Context, ids []uint) error
	Delete(ctx context.Context, id uint) error
	DeleteBatch(ctx context.Context, ids []uint) error
	Restore(ctx context.Context, id uint) error
	RestoreBatch(ctx context.Context, ids []uint) error
	GetByID(ctx context.Context, id uint) (*domain.Post, error)
	GetByKeyWord(ctx context.Context, keyWord string) ([]*domain.Post, error)
	GetDetailByID(ctx context.Context, id uint) (*domain.PostDetail, error)
	GetListByFilter(ctx context.Context, page *domain.Page, postType string) ([]*domain.PostDetail, int64, error)
	GetAllPublishPost(ctx context.Context) ([]*domain.PostDetail, error)
	FindByAlias(ctx context.Context, alias string) (*domain.Post, error)
}

var _ IPostRepository = (*PostRepository)(nil)

func NewPostRepository(dao dao.IPostDao) *PostRepository {
	return &PostRepository{dao: dao}
}

type PostRepository struct {
	dao dao.IPostDao
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	return r.dao.Create(ctx, r.PostDomainToPostDO(post))
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	return r.dao.Update(ctx, post.ID, r.PostDomainToUpdatePostDO(post))
}

// UpdatePostPublishStatus 更新文章发布状态
func (r *PostRepository) UpdatePublishStatus(ctx context.Context, id uint, isPublish bool) error {
	return r.dao.UpdatePostPublishStatus(ctx, id, isPublish)
}

func (r *PostRepository) SoftDelete(ctx context.Context, id uint) error {
	return r.dao.SoftDelete(ctx, id)
}

// SoftDeleteBatch 批量软删除文章
func (r *PostRepository) SoftDeleteBatch(ctx context.Context, ids []uint) error {
	return r.dao.SoftDeleteBatch(ctx, ids)
}

// Delete 删除文章
func (r *PostRepository) Delete(ctx context.Context, id uint) error {
	return r.dao.Delete(ctx, id)
}

// DeleteBatch 批量删除文章
func (r *PostRepository) DeleteBatch(ctx context.Context, ids []uint) error {
	return r.dao.DeleteBatch(ctx, ids)
}

// Restore 恢复文章
func (r *PostRepository) Restore(ctx context.Context, id uint) error {
	return r.dao.Restore(ctx, id)
}

// RestoreBatch 批量恢复文章
func (r *PostRepository) RestoreBatch(ctx context.Context, ids []uint) error {
	return r.dao.RestoreBatch(ctx, ids)
}

// GetByID 获取文章
func (r *PostRepository) GetByID(ctx context.Context, id uint) (*domain.Post, error) {
	post, err := r.dao.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.PostDOToPostDomain(post), nil
}

// GetDetailByID 获取文章
func (r *PostRepository) GetDetailByID(ctx context.Context, id uint) (*domain.PostDetail, error) {
	postCategoryTags, err := r.dao.GetWithCategoryAndTags(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.PostDOToPostDetail(postCategoryTags), nil
}

// GetByKeyWord 获取文章
func (r *PostRepository) GetByKeyWord(ctx context.Context, keyWord string) ([]*domain.Post, error) {
	posts, err := r.dao.GetListByKeyWord(ctx, keyWord)
	if err != nil {
		return nil, err
	}
	return lo.Map(posts, func(post *dao.Post, _ int) *domain.Post {
		return r.PostDOToPostDomain(post)
	}), nil
}

// buildPostQueryCondition 构建文章查询条件
func (r *PostRepository) buildPostQueryCondition(page *domain.Page, postType string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if page.Keyword != "" {
			tx = tx.Where("title LIKE ? OR description LIKE ?", "%"+page.Keyword+"%", "%"+page.Keyword+"%")
		}
		if page.Field != "" && page.Order != "" {
			tx = tx.Order(page.Field + " " + page.Order)
		}
		if page.PageNo != 0 && page.PageSize != 0 {
			tx = tx.Offset(int((page.PageNo - 1) * page.PageSize)).Limit(int(page.PageSize))
		}
		if postType == "publish" {
			tx = tx.Where("is_publish = ?", true)
		} else if postType == "draft" {
			tx = tx.Where("is_publish = ?", false)
		} else if postType == "bin" {
			tx = tx
		}
		return tx
	}
}

// GetList 获取文章列表
func (r *PostRepository) GetListByFilter(ctx context.Context, page *domain.Page, postType string) ([]*domain.PostDetail, int64, error) {
	count, err := r.dao.GetAllCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	filter := r.buildPostQueryCondition(page, postType)
	posts, err := r.dao.GetListByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return r.PostDOToPostDetailList(posts), count, nil
}

func (r *PostRepository) GetAllPublishPost(ctx context.Context) ([]*domain.PostDetail, error) {
	posts, err := r.dao.GetAllPublishPost(ctx)
	if err != nil {
		return nil, err
	}
	return r.PostDOToPostDetailList(posts), nil
}

func (r *PostRepository) FindByAlias(ctx context.Context, alias string) (*domain.Post, error) {
	post, err := r.dao.FindByAlias(ctx, alias)
	if err != nil {
		return nil, err
	}
	return r.PostDOToPostDomain(post), nil
}

// PostDomain2PostDO 将domain.Post转换为dao.Post
func (r *PostRepository) PostDomainToPostDO(post *domain.Post) *dao.Post {
	return &dao.Post{
		Model:       gorm.Model{CreatedAt: post.CreatedAt},
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		Alias:       post.Alias,
		CategoryID:  post.CategoryID,
		TagsID:      slice.UintToInt64(post.TagsID),
		IsPublish:   post.IsPublish,
		IsTop:       post.IsTop,
		Thumbnail:   post.Thumbnail,
	}
}

func (r *PostRepository) PostDomainToUpdatePostDO(post *domain.Post) *dao.PostUpdate {
	return &dao.PostUpdate{
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		Alias:       post.Alias,
		CategoryID:  post.CategoryID,
		TagsID:      slice.UintToInt64(post.TagsID),
		IsPublish:   post.IsPublish,
		IsTop:       post.IsTop,
		Thumbnail:   post.Thumbnail,
	}
}

func (r *PostRepository) PostDOToPostDomainList(posts []*dao.Post) []*domain.Post {
	return lo.Map(posts, func(post *dao.Post, _ int) *domain.Post {
		return r.PostDOToPostDomain(post)
	})
}

// PostDOToPostDomain 将dao.Post转换为domain.Post
func (r *PostRepository) PostDOToPostDomain(post *dao.Post) *domain.Post {
	return &domain.Post{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		Alias:       post.Alias,
		CategoryID:  post.CategoryID,
		TagsID:      slice.Int64ToUint(post.TagsID),
		IsPublish:   post.IsPublish,
		IsTop:       post.IsTop,
		Thumbnail:   post.Thumbnail,
	}
}

func (r *PostRepository) PostDOToPostDetailList(postCategoryTags []*dao.Post) []*domain.PostDetail {
	return lo.Map(postCategoryTags, func(postCategoryTag *dao.Post, _ int) *domain.PostDetail {
		return r.PostDOToPostDetail(postCategoryTag)
	})
}

func (r *PostRepository) PostDOToPostDetail(post *dao.Post) *domain.PostDetail {
	return &domain.PostDetail{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		Alias:       post.Alias,
		CategoryID:  post.CategoryID,
		TagsID:      slice.Int64ToUint(post.TagsID),
		IsPublish:   post.IsPublish,
		IsTop:       post.IsTop,
		Thumbnail:   post.Thumbnail,
	}
}

func (r *PostRepository) OrderConvertToInt(order string) int {
	switch order {
	case "DESC":
		return -1
	case "ASC":
		return 1
	default:
		return -1
	}
}
