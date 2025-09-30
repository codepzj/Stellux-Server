package service

import (
	"context"
	"errors"

	"github.com/codepzj/Stellux-Server/internal/post/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/post/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IPostService interface {
	AdminCreatePost(ctx context.Context, post *domain.Post) error
	AdminUpdatePost(ctx context.Context, post *domain.Post) error
	AdminUpdatePostPublishStatus(ctx context.Context, id uint, isPublish bool) error
	AdminSoftDeletePost(ctx context.Context, id uint) error
	AdminSoftDeletePostBatch(ctx context.Context, ids []uint) error
	AdminDeletePost(ctx context.Context, id uint) error
	AdminDeletePostBatch(ctx context.Context, ids []uint) error
	AdminRestorePost(ctx context.Context, id uint) error
	AdminRestorePostBatch(ctx context.Context, ids []uint) error
	GetPostById(ctx context.Context, id uint) (*domain.Post, error)
	GetPostByKeyWord(ctx context.Context, keyWord string) ([]*domain.Post, error)
	GetPostDetailById(ctx context.Context, id uint) (*domain.PostDetail, error)
	GetPostList(ctx context.Context, page *domain.Page, postType string) ([]*domain.PostDetail, int64, error)
	GetAllPublishPost(ctx context.Context) ([]*domain.PostDetail, error)
	FindByAlias(ctx context.Context, alias string) (*domain.Post, error)
}

var _ IPostService = (*PostService)(nil)

func NewPostService(repo repository.IPostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

type PostService struct {
	repo repository.IPostRepository
}

func (s *PostService) AdminCreatePost(ctx context.Context, post *domain.Post) error {
	if existPost, err := s.repo.FindByAlias(ctx, post.Alias); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	} else if existPost != nil {
		return errors.New("别名已存在")
	}
	return s.repo.Create(ctx, post)
}

func (s *PostService) AdminUpdatePost(ctx context.Context, post *domain.Post) error {
	if existPost, err := s.repo.FindByAlias(ctx, post.Alias); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	} else if existPost != nil && existPost.ID != post.ID {
		return errors.New("别名已存在")
	}
	return s.repo.Update(ctx, post)
}

func (s *PostService) AdminUpdatePostPublishStatus(ctx context.Context, id uint, isPublish bool) error {
	return s.repo.UpdatePublishStatus(ctx, id, isPublish)
}

func (s *PostService) AdminSoftDeletePost(ctx context.Context, id uint) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *PostService) AdminSoftDeletePostBatch(ctx context.Context, ids []uint) error {
	return s.repo.SoftDeleteBatch(ctx, ids)
}

func (s *PostService) AdminDeletePost(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *PostService) AdminDeletePostBatch(ctx context.Context, ids []uint) error {
	return s.repo.DeleteBatch(ctx, ids)
}

func (s *PostService) AdminRestorePost(ctx context.Context, id uint) error {
	return s.repo.Restore(ctx, id)
}

func (s *PostService) AdminRestorePostBatch(ctx context.Context, ids []uint) error {
	return s.repo.RestoreBatch(ctx, ids)
}

func (s *PostService) GetPostById(ctx context.Context, id uint) (*domain.Post, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PostService) GetPostByKeyWord(ctx context.Context, keyWord string) ([]*domain.Post, error) {
	return s.repo.GetByKeyWord(ctx, keyWord)
}

func (s *PostService) GetPostDetailById(ctx context.Context, id uint) (*domain.PostDetail, error) {
	return s.repo.GetDetailByID(ctx, id)
}

func (s *PostService) GetPostList(ctx context.Context, page *domain.Page, postType string) ([]*domain.PostDetail, int64, error) {
	switch postType {
	case "publish":
		return s.repo.GetListByFilter(ctx, page, "publish")
	case "draft":
		return s.repo.GetListByFilter(ctx, page, "draft")
	case "bin":
		return s.repo.GetListByFilter(ctx, page, "bin")
	default:
		return nil, 0, errors.New("invalid post type")
	}
}

// GetAllPublishPost 获取所有发布文章
func (s *PostService) GetAllPublishPost(ctx context.Context) ([]*domain.PostDetail, error) {
	return s.repo.GetAllPublishPost(ctx)
}

// FindByAlias 根据别名获取文章
func (s *PostService) FindByAlias(ctx context.Context, alias string) (*domain.Post, error) {
	return s.repo.FindByAlias(ctx, alias)
}
