package web

type SiteConfigSettingVO struct {
	SiteTitle string `json:"siteTitle"`
	SiteSubtitle  string `json:"siteSubtitle"`
	SiteUrl string `json:"siteUrl"`
	SiteFavicon string `json:"siteFavicon"`
	SiteAnimateText string `json:"siteAnimateText"`
	SiteAvatar string `json:"siteAvatar"`
	SiteAuthor string `json:"siteAuthor"`
	SiteDescription string `json:"siteDescription"`
	SiteKeywords string `json:"siteKeywords"`
	SiteCopyright string `json:"siteCopyright"`
	SiteIcp string `json:"siteIcp"`
	SiteIcpLink string `json:"siteIcpLink"`
	GithubUsername string `json:"githubUsername"`
}
