package web

import "github.com/codepzj/stellux/server/internal/pkg/apiwrap"

type DocumentVO struct {
	ID           string         `json:"id"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	DocumentType string         `json:"document_type"`
	IsPublic     bool           `json:"is_public"`
	ParentID     apiwrap.BsonID `json:"parent_id,omitzero"`
	DocumentID   apiwrap.BsonID `json:"document_id,omitzero"`
}

type DocumentRootVO struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Title       string `json:"title"`
	Alias       string `json:"alias"`
	DocumentType string `json:"document_type"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	IsPublic    bool   `json:"is_public"`
}

type DocumentTreeVO struct {
	ID           string         `json:"id"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
	Title        string         `json:"title"`
	DocumentType string         `json:"document_type"`
	IsPublic     bool           `json:"is_public"`
	ParentID     apiwrap.BsonID `json:"parent_id,omitzero"`
	DocumentID   apiwrap.BsonID `json:"document_id,omitzero"`
}

type DocumentSitemapVO struct {
	Loc        string `json:"loc"`
	Lastmod    string `json:"lastmod"`
	Changefreq string `json:"changefreq"`
	Priority   float64 `json:"priority"`
}

type SeoSettingVO struct {
	SiteAuthor    string `json:"site_author"`
	SiteUrl       string `json:"site_url"`
	SiteDescription string `json:"site_description"`
	SiteKeywords    string `json:"site_keywords"`
	Robots          string `json:"robots"`
	OgImage         string `json:"og_image"`
	OgType          string `json:"og_type"`
	TwitterCard     string `json:"twitter_card"`
	TwitterSite     string `json:"twitter_site"`
}