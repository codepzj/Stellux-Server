package service

import (
	"context"

	"github.com/codepzj/stellux/server/internal/comment/internal/domain"
	"github.com/codepzj/stellux/server/internal/comment/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ICommentService interface {
	Create(ctx context.Context, comment *domain.Comment) error
	AdminEdit(ctx context.Context, comment *domain.Comment) error
	AdminDelete(ctx context.Context, id bson.ObjectID) error
	GetListByPostId(ctx context.Context, postId bson.ObjectID) ([]*domain.Comment, error)
}

var _ ICommentService = (*CommentService)(nil)

func NewCommentService(repo repository.ICommentRepository) *CommentService {
	return &CommentService{
		repo: repo,
	}
}

type CommentService struct {
	repo repository.ICommentRepository
}

// Create 创建评论
func (s *CommentService) Create(ctx context.Context, comment *domain.Comment) error {
	return s.repo.Create(ctx, comment)
}

// AdminEdit 管理员编辑评论
func (s *CommentService) AdminEdit(ctx context.Context, comment *domain.Comment) error {
	return s.repo.Edit(ctx, comment)
}

// AdminDelete 管理员删除评论
func (s *CommentService) AdminDelete(ctx context.Context, id bson.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

// GetListByPostId 根据帖子id获取评论列表
func (s *CommentService) GetListByPostId(ctx context.Context, postId bson.ObjectID) ([]*domain.Comment, error) {
	return s.repo.GetListByPostId(ctx, postId)
}