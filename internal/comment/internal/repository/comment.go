package repository

import (
	"context"

	"github.com/codepzj/Stellux-Server/internal/comment/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/comment/internal/repository/dao"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ICommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	GetListByPath(ctx context.Context, path string) ([]*domain.CommentShow, error)
	Update(ctx context.Context, comment *domain.Comment) error
	Delete(ctx context.Context, id bson.ObjectID) error
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
		Path:     comment.Path,
		Content:  comment.Content,
		RootId:   comment.RootId,
		ParentId: comment.ParentId,
		Nickname: comment.Nickname,
		Avatar:   comment.Avatar,
		Email:    comment.Email,
		SiteUrl:  comment.SiteUrl,
		IsAdmin:  comment.IsAdmin,
	})
}

// GetListByPath 根据路径获取评论列表
func (r *CommentRepository) GetListByPath(ctx context.Context, path string) ([]*domain.CommentShow, error) {
	comments, err := r.dao.GetList(ctx, bson.D{{Key: "path", Value: path}})
	if err != nil {
		return nil, err
	}
	return r.CommentDaoToShowDomainList(comments), nil
}

// Update 更新评论
func (r *CommentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	return r.dao.Update(ctx, comment.Id, &dao.Comment{
		Content:  comment.Content,
		Nickname: comment.Nickname,
		Avatar:   comment.Avatar,
		SiteUrl:  comment.SiteUrl,
		IsAdmin:  comment.IsAdmin,
	})
}

// Delete 删除评论
func (r *CommentRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	return r.dao.Delete(ctx, id)
}

func (r *CommentRepository) CommentDaoToShowDomainList(comments []*dao.Comment) []*domain.CommentShow {
	return lo.Map(comments, func(comment *dao.Comment, _ int) *domain.CommentShow {
		return r.CommentDaoToShowDomain(comment)
	})
}

// 转成前端展示的评论
func (r *CommentRepository) CommentDaoToShowDomain(comment *dao.Comment) *domain.CommentShow {
	return &domain.CommentShow{
		Id:        comment.ID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		DeletedAt: comment.DeletedAt,
		Path:      comment.Path,
		Content:   comment.Content,
		Nickname:  comment.Nickname,
		Avatar:    comment.Avatar,
		SiteUrl:   comment.SiteUrl,
		IsAdmin:   comment.IsAdmin,
	}
}
