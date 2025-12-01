package service

import (
	"context"
	"errors"

	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"

	"github.com/codepzj/Stellux-Server/internal/pkg/logger"
	"github.com/codepzj/Stellux-Server/internal/post/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/post/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IPostService interface {
	AdminCreatePost(ctx context.Context, post *domain.Post) error
	AdminUpdatePost(ctx context.Context, post *domain.Post) error
	AdminUpdatePostPublishStatus(ctx context.Context, id bson.ObjectID, isPublish bool) error
	AdminSoftDeletePost(ctx context.Context, id bson.ObjectID) error
	AdminSoftDeletePostBatch(ctx context.Context, ids []bson.ObjectID) error
	AdminDeletePost(ctx context.Context, id bson.ObjectID) error
	AdminDeletePostBatch(ctx context.Context, ids []bson.ObjectID) error
	AdminRestorePost(ctx context.Context, id bson.ObjectID) error
	AdminRestorePostBatch(ctx context.Context, ids []bson.ObjectID) error
	GetPostById(ctx context.Context, id bson.ObjectID) (*domain.Post, error)
	GetPostByKeyWord(ctx context.Context, keyWord string) ([]*domain.Post, error)
	GetPostDetailById(ctx context.Context, id bson.ObjectID) (*domain.PostDetail, error)
	GetPostList(ctx context.Context, page *apiwrap.Page, labelName, categoryName, postType string) ([]*domain.PostDetail, int64, error)
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
		logger.Error("查询别名失败",
			logger.WithError(err),
			logger.WithString("alias", post.Alias),
		)
		return err
	} else if existPost != nil {
		logger.Warn("文章别名已存在",
			logger.WithString("alias", post.Alias),
			logger.WithString("existPostId", existPost.Id.Hex()),
		)
		return errors.New("别名已存在")
	}

	err := s.repo.Create(ctx, post)
	if err != nil {
		logger.Error("创建文章失败",
			logger.WithError(err),
			logger.WithString("title", post.Title),
		)
		return err
	}

	logger.Info("创建文章成功",
		logger.WithString("postId", post.Id.Hex()),
		logger.WithString("title", post.Title),
	)

	return nil
}

func (s *PostService) AdminUpdatePost(ctx context.Context, post *domain.Post) error {
	if existPost, err := s.repo.FindByAlias(ctx, post.Alias); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		logger.Error("查询别名失败",
			logger.WithError(err),
			logger.WithString("alias", post.Alias),
		)
		return err
	} else if existPost != nil && existPost.Id != post.Id {
		logger.Warn("别名已被其他文章使用",
			logger.WithString("alias", post.Alias),
			logger.WithString("existPostId", existPost.Id.Hex()),
		)
		return errors.New("别名已存在")
	}

	err := s.repo.Update(ctx, post)
	if err != nil {
		logger.Error("更新文章失败",
			logger.WithError(err),
			logger.WithString("postId", post.Id.Hex()),
		)
		return err
	}

	logger.Info("更新文章成功",
		logger.WithString("postId", post.Id.Hex()),
		logger.WithString("title", post.Title),
	)

	return nil
}

func (s *PostService) AdminUpdatePostPublishStatus(ctx context.Context, id bson.ObjectID, isPublish bool) error {
	err := s.repo.UpdatePublishStatus(ctx, id, isPublish)
	if err != nil {
		logger.Error("更新发布状态失败",
			logger.WithError(err),
			logger.WithString("postId", id.Hex()),
		)
		return err
	}

	logger.Info("更新发布状态成功",
		logger.WithString("postId", id.Hex()),
		logger.WithAny("isPublish", isPublish),
	)

	return nil
}

func (s *PostService) AdminSoftDeletePost(ctx context.Context, id bson.ObjectID) error {
	err := s.repo.SoftDelete(ctx, id)
	if err != nil {
		logger.Error("软删除文章失败",
			logger.WithError(err),
			logger.WithString("postId", id.Hex()),
		)
		return err
	}

	logger.Info("软删除文章成功",
		logger.WithString("postId", id.Hex()),
	)

	return nil
}

func (s *PostService) AdminSoftDeletePostBatch(ctx context.Context, ids []bson.ObjectID) error {
	err := s.repo.SoftDeleteBatch(ctx, ids)
	if err != nil {
		logger.Error("批量软删除失败",
			logger.WithError(err),
			logger.WithInt("count", len(ids)),
		)
		return err
	}

	logger.Info("批量软删除成功",
		logger.WithInt("count", len(ids)),
	)

	return nil
}

func (s *PostService) AdminDeletePost(ctx context.Context, id bson.ObjectID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		logger.Error("删除文章失败",
			logger.WithError(err),
			logger.WithString("postId", id.Hex()),
		)
		return err
	}

	logger.Info("删除文章成功",
		logger.WithString("postId", id.Hex()),
	)

	return nil
}

func (s *PostService) AdminDeletePostBatch(ctx context.Context, ids []bson.ObjectID) error {
	err := s.repo.DeleteBatch(ctx, ids)
	if err != nil {
		logger.Error("批量删除失败",
			logger.WithError(err),
			logger.WithInt("count", len(ids)),
		)
		return err
	}

	logger.Info("批量删除成功",
		logger.WithInt("count", len(ids)),
	)

	return nil
}

func (s *PostService) AdminRestorePost(ctx context.Context, id bson.ObjectID) error {
	err := s.repo.Restore(ctx, id)
	if err != nil {
		logger.Error("恢复文章失败",
			logger.WithError(err),
			logger.WithString("postId", id.Hex()),
		)
		return err
	}

	logger.Info("恢复文章成功",
		logger.WithString("postId", id.Hex()),
	)

	return nil
}

func (s *PostService) AdminRestorePostBatch(ctx context.Context, ids []bson.ObjectID) error {
	err := s.repo.RestoreBatch(ctx, ids)
	if err != nil {
		logger.Error("批量恢复失败",
			logger.WithError(err),
			logger.WithInt("count", len(ids)),
		)
		return err
	}

	logger.Info("批量恢复成功",
		logger.WithInt("count", len(ids)),
	)
	return nil
}

func (s *PostService) GetPostById(ctx context.Context, id bson.ObjectID) (*domain.Post, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("查询文章失败",
			logger.WithError(err),
			logger.WithString("postId", id.Hex()),
		)
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetPostByKeyWord(ctx context.Context, keyWord string) ([]*domain.Post, error) {
	posts, err := s.repo.GetByKeyWord(ctx, keyWord)
	if err != nil {
		logger.Error("搜索文章失败",
			logger.WithError(err),
			logger.WithString("keyword", keyWord),
		)
		return nil, err
	}
	return posts, nil
}

func (s *PostService) GetPostDetailById(ctx context.Context, id bson.ObjectID) (*domain.PostDetail, error) {
	detail, err := s.repo.GetDetailByID(ctx, id)
	if err != nil {
		logger.Error("查询文章详情失败",
			logger.WithError(err),
			logger.WithString("postId", id.Hex()),
		)
		return nil, err
	}
	return detail, nil
}

func (s *PostService) GetPostList(ctx context.Context, page *apiwrap.Page, labelName, categoryName, postType string) ([]*domain.PostDetail, int64, error) {
	var posts []*domain.PostDetail
	var total int64
	var err error

	// 将apiwrap.Page转换为domain.PostQueryPage
	queryPage := &domain.PostQueryPage{
		Page:         *page,
		LabelName:    labelName,
		CategoryName: categoryName,
	}

	switch postType {
	case "publish":
		posts, total, err = s.repo.GetList(ctx, queryPage, "publish")
	case "draft":
		posts, total, err = s.repo.GetList(ctx, queryPage, "draft")
	case "bin":
		posts, total, err = s.repo.GetList(ctx, queryPage, "bin")
	default:
		logger.Warn("无效的文章类型",
			logger.WithString("postType", postType),
		)
		return nil, 0, errors.New("invalid post type")
	}

	if err != nil {
		logger.Error("查询文章列表失败",
			logger.WithError(err),
			logger.WithString("postType", postType),
		)
		return nil, 0, err
	}
	return posts, total, nil
}

// GetAllPublishPost 获取所有发布文章
func (s *PostService) GetAllPublishPost(ctx context.Context) ([]*domain.PostDetail, error) {
	posts, err := s.repo.GetAllPublishPost(ctx)
	if err != nil {
		logger.Error("查询所有发布文章失败",
			logger.WithError(err),
		)
		return nil, err
	}
	return posts, nil
}

// FindByAlias 根据别名获取文章
func (s *PostService) FindByAlias(ctx context.Context, alias string) (*domain.Post, error) {
	post, err := s.repo.FindByAlias(ctx, alias)
	if err != nil {
		logger.Error("查询文章失败",
			logger.WithError(err),
			logger.WithString("alias", alias),
		)
		return nil, err
	}
	return post, nil
}
