package web

import "time"

type DocumentVO struct {
	Id          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	Alias       string    `json:"alias"`
	Sort        int       `json:"sort"`
	IsPublic    bool      `json:"is_public"`
	IsDeleted   bool      `json:"is_deleted"`
}
