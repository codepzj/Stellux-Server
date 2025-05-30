package web

type BasicSettingRequest struct {
	SiteTitle    string `json:"site_title" binding:"required"`
	SiteSubtitle string `json:"site_subtitle" binding:"required"`
	SiteFavicon  string `json:"site_favicon"`
}

type SeoSettingRequest struct {
	SiteAuthor    string `json:"site_author" binding:"required"`
	SiteUrl       string `json:"site_url" binding:"required"`
	SiteDescription string `json:"site_description" binding:"required"`
	SiteKeywords    string `json:"site_keywords" binding:"required"`
	Robots          string `json:"robots" binding:"required"`
	OgImage         string `json:"og_image" binding:"required"`
	OgType          string `json:"og_type" binding:"required"`
	TwitterCard     string `json:"twitter_card" binding:"required"`
	TwitterSite     string `json:"twitter_site" binding:"required"`
}

type BlogSettingRequest struct {
	BlogAvatar    string `json:"blog_avatar" binding:"required"`
	BlogTitle     string `json:"blog_title" binding:"required"`
	BlogSubtitle  string `json:"blog_subtitle" binding:"required"`
	BlogWelcome   string `json:"blog_welcome" binding:"required"`
	BlogMotto     string `json:"blog_motto" binding:"required"`
}