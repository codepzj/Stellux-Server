package web

type SiteConfigSettingVO struct {
	SiteTitle string `json:"site_title"`
	SiteSubTitle  string `json:"site_subtitle"`
	SiteFavicon string `json:"site_favicon"`
	SiteAvatar string `json:"site_avatar"`
	SiteDescription string `json:"site_description"`
	SiteKeywords string `json:"site_keywords"`
	SiteCopyright string `json:"site_copyright"`
	SiteIcp string `json:"site_icp"`
	SiteIcpLink string `json:"site_icplink"`
	GithubUsername string `json:"github_username"`
}
