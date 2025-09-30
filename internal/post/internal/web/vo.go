package web

import (
	"time"

	"github.com/codepzj/Stellux-Server/internal/label"
	"github.com/codepzj/Stellux-Server/internal/post/internal/domain"
	"github.com/samber/lo"
)

type PostVO struct {
	ID          uint    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Alias       string    `json:"alias"`
	CategoryID  uint    `json:"category_id"`
	TagsID      []uint  `json:"tags_id"`
	IsPublish   bool      `json:"is_publish"`
	IsTop       bool      `json:"is_top"`
	Thumbnail   string    `json:"thumbnail"`
}

type PostDetailVO struct {
	ID          uint    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Alias       string    `json:"alias"`
	CategoryID  uint    `json:"category_id"`
	TagsID      []uint  `json:"tags_id"`
	IsPublish   bool      `json:"is_publish"`
	IsTop       bool      `json:"is_top"`
	Thumbnail   string    `json:"thumbnail"`
}

func GetCategoryNameFromLabel(label label.Domain) uint {
	return label.ID
}

func GetTagNamesFromLabels(labels []label.Domain) []uint {
	return lo.Map(labels, func(label label.Domain, _ int) uint {
		return label.ID
	})
}

func (h *PostHandler) PostDTOToDomain(postReq PostDto) *domain.Post {
	return &domain.Post{
		CreatedAt:   postReq.CreatedAt,
		Title:       postReq.Title,
		Content:     postReq.Content,
		Description: postReq.Description,
		Author:      postReq.Author,
		Alias:       postReq.Alias,
		CategoryID:  postReq.CategoryID,
		TagsID:      postReq.TagsID,
		IsPublish:   postReq.IsPublish,
		IsTop:       postReq.IsTop,
		Thumbnail:   postReq.Thumbnail,
	}
}

func (h *PostHandler) PostUpdateDTOToDomain(postUpdateReq PostUpdateDto) *domain.Post {
	return &domain.Post{
		ID:          postUpdateReq.ID,
		CreatedAt:   postUpdateReq.CreatedAt,
		Title:       postUpdateReq.Title,
		Content:     postUpdateReq.Content,
		Description: postUpdateReq.Description,
		Author:      postUpdateReq.Author,
		Alias:       postUpdateReq.Alias,
		CategoryID:  postUpdateReq.CategoryID,
		TagsID:      postUpdateReq.TagsID,
		IsPublish:   postUpdateReq.IsPublish,
		IsTop:       postUpdateReq.IsTop,
		Thumbnail:   postUpdateReq.Thumbnail,
	}
}

func (h *PostHandler) PostDetailToVO(post *domain.PostDetail) *PostDetailVO {
	return &PostDetailVO{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		Alias:       post.Alias,
		CategoryID:  post.CategoryID,
		TagsID:      GetTagNamesFromLabels(post.Tags),
		IsPublish:   post.IsPublish,
		IsTop:       post.IsTop,
		Thumbnail:   post.Thumbnail,
	}
}

func (h *PostHandler) PostDetailListToVOList(posts []*domain.PostDetail) []*PostDetailVO {
	return lo.Map(posts, func(post *domain.PostDetail, _ int) *PostDetailVO {
		return h.PostDetailToVO(post)
	})
}

func (h *PostHandler) PostToVO(post *domain.Post) *PostVO {
	return &PostVO{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		Alias:       post.Alias,
		CategoryID:  post.CategoryID,
		TagsID:      post.TagsID,
		IsPublish:   post.IsPublish,
		IsTop:       post.IsTop,
		Thumbnail:   post.Thumbnail,
	}
}

func (h *PostHandler) PostListToVOList(posts []*domain.Post) []*PostVO {
	return lo.Map(posts, func(post *domain.Post, _ int) *PostVO {
		return h.PostToVO(post)
	})
}
