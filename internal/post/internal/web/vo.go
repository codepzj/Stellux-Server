package web

import (
	"time"

	"github.com/codepzj/stellux/server/internal/label"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/codepzj/stellux/server/internal/post/internal/domain"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PostVO struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
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
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	IsPublish   bool      `json:"is_publish"`
	IsTop       bool      `json:"is_top"`
	Thumbnail   string    `json:"thumbnail"`
}

type SiteMapVO struct {
	Loc        string  `json:"loc"`
	Lastmod    string  `json:"lastmod"`
	Changefreq string  `json:"changefreq"`
	Priority   float64 `json:"priority"`
}

type SeoSettingVO struct {
	SiteAuthor      string `json:"site_author"`
	SiteUrl         string `json:"site_url"`
	SiteDescription string `json:"site_description"`
	SiteKeywords    string `json:"site_keywords"`
	Robots          string `json:"robots"`
	OgImage         string `json:"og_image"`
	OgType          string `json:"og_type"`
	TwitterCard     string `json:"twitter_card"`
	TwitterSite     string `json:"twitter_site"`
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
		CategoryID:  apiwrap.ConvertBsonID(postReq.CategoryID).ToObjectID(),
		TagsID:      apiwrap.ToObjectIDList(apiwrap.ConvertBsonIDList(postReq.TagsID)),
		IsPublish:   postReq.IsPublish,
		IsTop:       postReq.IsTop,
		Thumbnail:   postReq.Thumbnail,
	}
}

func (h *PostHandler) PostUpdateDTOToDomain(postUpdateReq PostUpdateDto) *domain.Post {
	return &domain.Post{
		ID:          apiwrap.ConvertBsonID(postUpdateReq.Id).ToObjectID(),
		CreatedAt:   postUpdateReq.CreatedAt,
		Title:       postUpdateReq.Title,
		Content:     postUpdateReq.Content,
		Description: postUpdateReq.Description,
		Author:      postUpdateReq.Author,
		CategoryID:  apiwrap.ConvertBsonID(postUpdateReq.CategoryID).ToObjectID(),
		TagsID:      apiwrap.ToObjectIDList(apiwrap.ConvertBsonIDList(postUpdateReq.TagsID)),
		IsPublish:   postUpdateReq.IsPublish,
		IsTop:       postUpdateReq.IsTop,
		Thumbnail:   postUpdateReq.Thumbnail,
	}
}

func (h *PostHandler) PostDetailToVO(post *domain.PostDetail) *PostDetailVO {
	return &PostDetailVO{
		ID:          post.ID.Hex(),
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
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
	return &PostVO{
		ID:          post.ID.Hex(),
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		CategoryID:  post.CategoryID.Hex(),
		TagsID: lo.Map(post.TagsID, func(id bson.ObjectID, _ int) string {
			return id.Hex()
		}),
		IsPublish: post.IsPublish,
		IsTop:     post.IsTop,
		Thumbnail: post.Thumbnail,
	}
}

func (h *PostHandler) PostListToVOList(posts []*domain.Post) []*PostVO {
	return lo.Map(posts, func(post *domain.Post, _ int) *PostVO {
		return h.PostToVO(post)
	})
}

func (h *PostHandler) PostListToSiteMapVOList(posts []*domain.Post, siteUrl string) []*SiteMapVO {
	return lo.Map(posts, func(post *domain.Post, _ int) *SiteMapVO {
		return &SiteMapVO{
			Loc:        siteUrl + "/post/" + post.ID.Hex(),
			Lastmod:    post.UpdatedAt.String(),
			Changefreq: "weekly",
			Priority:   0.8,
		}
	})
}
