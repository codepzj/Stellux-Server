package web

import (
	"time"
)

type Page struct {
	// 当前页
	PageNo int64 `form:"page_no" json:"page_no" binding:"required,gte=1"`
	// 每页条数
	PageSize int64 `form:"page_size" json:"page_size" binding:"required,gte=1"`
	// 排序字段
	Field string `form:"field" json:"field,omitempty"`
	// 排序方式
	Order string `form:"order" json:"order,omitempty" binding:"omitempty,oneof=ASC DESC"`
	// 搜索内容
	Keyword string `form:"keyword" json:"keyword,omitempty"`
	// 标签名称
	LabelName string `form:"label_name" json:"label_name,omitempty"`
	// 分类名称
	CategoryName string `form:"category_name" json:"category_name,omitempty"`
}

type PostDto struct {
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title" binding:"required"`
	Content     string    `json:"content" binding:"required"`
	Description string    `json:"description"`
	Author      string    `json:"author" binding:"required"`
	Alias       string    `json:"alias" binding:"required"`
	CategoryID  string    `json:"category_id"`
	TagsID      []string  `json:"tags_id"`
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
	CategoryID  string    `json:"category_id"`
	TagsID      []string  `json:"tags_id"`
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
