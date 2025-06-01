package web

type SettingVO struct{}

type BasicSettingVO struct {
	SiteTitle    string `json:"site_title"`
	SiteSubtitle string `json:"site_subtitle"`
	SiteLogo     string `json:"site_logo"`
	SiteFavicon  string `json:"site_favicon"`
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

type BlogSettingVO struct {
	BlogDarkLogo string `json:"blog_dark_logo"`
	BlogLightLogo string `json:"blog_light_logo"`
}

type AboutSettingVO struct {
	Author        string   `json:"author"`
	AvatarUrl     string   `json:"avatar_url"`
	LeftTags      []string `json:"left_tags"`
	RightTags     []string `json:"right_tags"`
	KnowMe        string   `json:"know_me"`
	GithubUsername string   `json:"github_username"`
}