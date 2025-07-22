package web

import "time"

type DocumentVO struct {
	Id          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	Alias       string    `json:"alias"`
	Sort        int       `json:"sort"`
	IsPublic    bool      `json:"is_public"`
	IsDeleted   bool      `json:"is_deleted"`
}

type DocumentSitemapVO struct {
	Loc        string  `json:"loc"`
	Lastmod    string  `json:"lastmod"`
	Changefreq string  `json:"changefreq"`
	Priority   float64 `json:"priority"`
}
