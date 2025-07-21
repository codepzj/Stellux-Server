package web

import "time"

type DocumentContentVO struct {
	Id           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DocumentId   string    `json:"document_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Description  string    `json:"description"`
	Alias        string    `json:"alias"`
	ParentId     string    `json:"parent_id"`
	IsDir        bool      `json:"is_dir"`
	Sort         int       `json:"sort"`
	LikeCount    int       `json:"like_count"`
	DislikeCount int       `json:"dislike_count"`
	CommentCount int       `json:"comment_count"`
}
