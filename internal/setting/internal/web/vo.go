package web

type SiteConfigSettingVO struct {
	SiteTitle string `json:"site_title"`
	SiteSubTitle  string `json:"site_sub_title"`
	SiteLogo string `json:"site_logo"`
	SiteFavicon string `json:"site_favicon"`
	SiteDescription string `json:"site_description"`
	SiteKeywords string `json:"site_keywords"`
	SiteCopyright string `json:"site_copyright"`
	SiteIcp string `json:"site_icp"`
	SiteIcpUrl string `json:"site_icp_url"`
	GithubUsername string `json:"github_username"`
}
