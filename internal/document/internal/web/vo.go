package web

import "time"

type DocumentVO struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Alias       string    `json:"alias"`
	Sort        int       `json:"sort"`
	IsPublic    bool      `json:"is_public"`
	IsDeleted   bool      `json:"is_deleted"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DocumentRootVO struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Alias       string    `json:"alias"`
	Thumbnail   string    `json:"thumbnail"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DocumentTreeVO struct {
	Id           string    `json:"id"`
	Title        string    `json:"title"`
	DocumentType string    `json:"document_type"`
	IsDir        bool      `json:"is_dir"`
	IsPublic     bool      `json:"is_public"`
	ParentId     string    `json:"parent_id"`
	DocumentId   string    `json:"document_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DocumentSitemapVO struct {
	Loc        string  `json:"loc"`
	Lastmod    string  `json:"lastmod"`
	Changefreq string  `json:"changefreq"`
	Priority   float64 `json:"priority"`
}
