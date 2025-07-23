package repository

import (
	"context"

	"github.com/codepzj/stellux/server/internal/comment/internal/domain"
	"github.com/codepzj/stellux/server/internal/comment/internal/repository/dao"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ICommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	Edit(ctx context.Context, comment *domain.Comment) error
	Delete(ctx context.Context, id bson.ObjectID) error
	GetListByPostId(ctx context.Context, postId bson.ObjectID) ([]*domain.Comment, error)
}

var _ ICommentRepository = (*CommentRepository)(nil)

func NewCommentRepository(dao dao.ICommentDao) *CommentRepository {
	return &CommentRepository{dao: dao}
}

type CommentRepository struct {
	dao dao.ICommentDao
}

// Create 创建评论
func (r *CommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	return r.dao.Create(ctx, &dao.Comment{
		Nickname:  comment.Nickname,
		Email:     comment.Email,
		Avatar:    comment.Avatar,
		Url:       comment.Url,
		Content:   comment.Content,
		PostId:    comment.PostId,
		CommentId: comment.CommentId,
	})
}

// Edit 编辑评论
func (r *CommentRepository) Edit(ctx context.Context, comment *domain.Comment) error {
	err := r.dao.Edit(ctx, bson.D{
		{Key: "id", Value: comment.Id},
		{Key: "nickname", Value: comment.Nickname},
		{Key: "email", Value: comment.Email},
		{Key: "avatar", Value: comment.Avatar},
		{Key: "url", Value: comment.Url},
		{Key: "content", Value: comment.Content},
	})
	if err != nil {
		return err
	}
	return nil
}

// Delete 删除评论
func (r *CommentRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	return r.dao.Delete(ctx, id)
}

// GetListByPostId 根据帖子id获取评论列表
func (r *CommentRepository) GetListByPostId(ctx context.Context, postId bson.ObjectID) ([]*domain.Comment, error) {
	comments, err := r.dao.GetListByPostId(ctx, postId)
	if err != nil {
		return nil, err
	}
	domainComments := make([]*domain.Comment, len(comments))
	for i, comment := range comments {
		domainComments[i] = &domain.Comment{
			Id:        comment.ID,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Nickname:  comment.Nickname,
			Email:     comment.Email,
			Avatar:    comment.Avatar,
			Url:       comment.Url,
			Content:   comment.Content,
			PostId:    comment.PostId,
			CommentId: comment.CommentId,
		}
	}
	return domainComments, nil
}
