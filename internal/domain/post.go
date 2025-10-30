package domain

import (
	"context"
	"errors"
	"time"

	"github.com/codepzj/Stellux-Server/internal/domain/model"
	"go.uber.org/zap"
)

type Post struct {
	ID          uint       // 文章ID
	CreatedAt   time.Time  // 创建时间
	UpdatedAt   time.Time  // 更新时间
	DeletedAt   *time.Time // 删除时间
	Title       string     // 标题
	Content     string     // 内容
	Description string     // 描述
	Author      string     // 作者
	Alias       string     // 别名
	CategoryID  uint       // 分类ID
	IsPublish   bool       // 是否发布
	IsTop       bool       // 是否置顶
	Thumbnail   string     // 缩略图
}

type PostRepo struct {
	data *Data
	log  *zap.SugaredLogger
}

func NewPostRepo(data *Data, log *zap.SugaredLogger) *PostRepo {
	return &PostRepo{data: data, log: log}
}

func (r *PostRepo) CreatePost(ctx context.Context, post *model.Post) error {
	return r.data.query.Post.WithContext(ctx).Create(post)
}

func (r *PostRepo) UpdatePost(ctx context.Context, post *model.Post) error {
	updateRes, err := r.data.query.Post.WithContext(ctx).Where(r.data.query.Post.ID.Eq(post.ID)).Updates(post)
	if err != nil {
		return err
	}
	if updateRes.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	return nil
}

func (r *PostRepo) DeletePost(ctx context.Context, id int64) error {
	deleteRes, err := r.data.query.Post.WithContext(ctx).Where(r.data.query.Post.ID.Eq(id)).Delete()
	if err != nil {
		return err
	}
	if deleteRes.RowsAffected == 0 {
		return errors.New("删除失败")
	}
	return nil
}

func (r *PostRepo) GetPost(ctx context.Context, id int64) (*model.Post, error) {
	post, err := r.data.query.Post.WithContext(ctx).Where(r.data.query.Post.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return post, nil
}
