package web

import "github.com/codepzj/stellux/server/internal/pkg/apiwrap"

// 创建文档内容请求
type CreateDocumentContentDto struct {
	DocumentId   apiwrap.BsonId `json:"documentId" binding:"required"`
	Title        string         `json:"title" binding:"required"`
	Content      string         `json:"content" binding:"required"`
	Version      int            `json:"version" binding:"required"`
	Alias        string         `json:"alias" binding:"required"`
	ParentID     apiwrap.BsonId `json:"parentId" binding:"required"`
	IsDir        bool           `json:"isDir"`
	LikeCount    int            `json:"likeCount"`
	DislikeCount int            `json:"dislikeCount"`
	CommentCount int            `json:"commentCount"`
}
