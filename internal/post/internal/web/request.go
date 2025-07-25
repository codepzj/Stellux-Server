package web

import (
	"time"

	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
)

type BsonId = apiwrap.BsonId

type PostDto struct {
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title" binding:"required"`
	Content     string    `json:"content" binding:"required"`
	Description string    `json:"description"`
	Author      string    `json:"author" binding:"required"`
	Alias       string    `json:"alias" binding:"required"`
	CategoryID  string    `json:"category_id" binding:"required"`
	TagsID      []string  `json:"tags_id" binding:"required"`
	IsPublish   bool      `json:"is_publish"`
	IsTop       bool      `json:"is_top"`
	Thumbnail   string    `json:"thumbnail"`
}

type PostUpdateDto struct {
	Id          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Alias       string    `json:"alias" binding:"required"`
	CategoryID  string    `json:"category_id" binding:"required"`
	TagsID      []string  `json:"tags_id" binding:"required"`
	IsPublish   bool      `json:"is_publish"`
	IsTop       bool      `json:"is_top"`
	Thumbnail   string    `json:"thumbnail"`
}

type PostIdRequest struct {
	Id string `uri:"id" binding:"required"`
}

type PostPublishStatusRequest struct {
	ID        string `json:"id" binding:"required"`
	IsPublish *bool  `json:"is_publish" binding:"required"`
}

type PostIDListRequest struct {
	IDList []string `json:"id_list" binding:"required"`
}
