package web

import "time"

type DocumentContentVO struct {
	Id           string    `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	DocumentId   string    `json:"documentId"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Description  string    `json:"description"`
	Version      string    `json:"version"`
	Alias        string    `json:"alias"`
	ParentId     string    `json:"parentId"`
	IsDir        bool      `json:"isDir"`
	Sort         int       `json:"sort"`
	LikeCount    int       `json:"likeCount"`
	DislikeCount int       `json:"dislikeCount"`
	CommentCount int       `json:"commentCount"`
}
