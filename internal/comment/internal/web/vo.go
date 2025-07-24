package web

import "time"

type CommentVO struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Url       string    `json:"url"`
	Avatar    string    `json:"avatar"`
	Content   string    `json:"content"`
	PostId    string    `json:"post_id"`
	CommentId string    `json:"comment_id"`
}
