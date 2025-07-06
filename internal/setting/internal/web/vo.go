package web

type SiteConfigSettingVO struct {
	SiteTitle string `json:"siteTitle"`
	SiteSubtitle  string `json:"siteSubtitle"`
	SiteUrl string `json:"siteUrl"`
	SiteFavicon string `json:"siteFavicon"`
	SiteAvatar string `json:"siteAvatar"`
	SiteDescription string `json:"siteDescription"`
	SiteKeywords string `json:"siteKeywords"`
	SiteCopyright string `json:"siteCopyright"`
	SiteIcp string `json:"siteIcp"`
	SiteIcpLink string `json:"siteIcpLink"`
	GithubUsername string `json:"githubUsername"`
}
