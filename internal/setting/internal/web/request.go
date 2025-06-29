package web

type SiteConfigRequest struct {
	SiteTitle    string `json:"site_title" binding:"required"`
	SiteSubtitle string `json:"site_sub_title" binding:"required"`
	SiteFavicon  string `json:"site_favicon"`
	SiteAvatar   string `json:"site_avatar"`
	SiteKeywords string `json:"site_keywords"`
	SiteDescription string `json:"site_description"`
	SiteCopyright string `json:"site_copyright"`
	SiteICP string `json:"site_icp"`
	SiteICPLink string `json:"site_icp_link"`
	GithubUsername string `json:"github_username"`
}