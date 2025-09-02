package web

import (
	"time"

	"github.com/codepzj/Stellux-Server/internal/label"
	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"
	"github.com/codepzj/Stellux-Server/internal/post/internal/domain"
	"github.com/samber/lo"
)

type PostVO struct {
	Id          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Alias       string    `json:"alias"`
	CategoryID  string    `json:"category_id"`
	TagsID      []string  `json:"tags_id"`
	IsPublish   bool      `json:"is_publish"`
	IsTop       bool      `json:"is_top"`
	Thumbnail   string    `json:"thumbnail"`
}

type PostDetailVO struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Alias       string    `json:"alias"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	IsPublish   bool      `json:"is_publish"`
	IsTop       bool      `json:"is_top"`
	Thumbnail   string    `json:"thumbnail"`
}

func GetCategoryNameFromLabel(label label.Domain) string {
	return label.Name
}

func GetTagNamesFromLabels(labels []label.Domain) []string {
	return lo.Map(labels, func(label label.Domain, _ int) string {
		return label.Name
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
		CategoryId:  apiwrap.ConvertBsonID(postReq.CategoryID).ToObjectID(),
		TagsId:      apiwrap.ToObjectIDList(apiwrap.ConvertBsonIDList(postReq.TagsID)),
		IsPublish:   postReq.IsPublish,
		IsTop:       postReq.IsTop,
		Thumbnail:   postReq.Thumbnail,
	}
}

func (h *PostHandler) PostUpdateDTOToDomain(postUpdateReq PostUpdateDto) *domain.Post {
	return &domain.Post{
		Id:          apiwrap.ConvertBsonID(postUpdateReq.Id).ToObjectID(),
		CreatedAt:   postUpdateReq.CreatedAt,
		Title:       postUpdateReq.Title,
		Content:     postUpdateReq.Content,
		Description: postUpdateReq.Description,
		Author:      postUpdateReq.Author,
		Alias:       postUpdateReq.Alias,
		CategoryId:  apiwrap.ConvertBsonID(postUpdateReq.CategoryID).ToObjectID(),
		TagsId:      apiwrap.ToObjectIDList(apiwrap.ConvertBsonIDList(postUpdateReq.TagsID)),
		IsPublish:   postUpdateReq.IsPublish,
		IsTop:       postUpdateReq.IsTop,
		Thumbnail:   postUpdateReq.Thumbnail,
	}
}

func (h *PostHandler) PostDetailToVO(post *domain.PostDetail) *PostDetailVO {
	return &PostDetailVO{
		ID:          post.Id.Hex(),
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		Alias:       post.Alias,
		Category:    GetCategoryNameFromLabel(post.Category),
		Tags:        GetTagNamesFromLabels(post.Tags),
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
	categoryId := ""
	tagsId := []string{}
	// mongodb查询得到的categoryId和tagsId可能出现空值，转成字符串需要处理
	if !post.CategoryId.IsZero() {
		categoryId = post.CategoryId.Hex()
	}
	for i := 0; i < len(post.TagsId); i++ {
		if !post.TagsId[i].IsZero() {
			tagsId = append(tagsId, post.TagsId[i].Hex())
		}
	}
	return &PostVO{
		Id:          post.Id.Hex(),
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		Alias:       post.Alias,
		CategoryID:  categoryId,
		TagsID:      tagsId,
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
