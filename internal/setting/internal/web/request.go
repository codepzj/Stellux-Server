package web

type SiteConfigRequest struct {
	SiteTitle       string `json:"siteTitle" binding:"required"`
	SiteSubtitle    string `json:"siteSubtitle" binding:"required"`
	SiteUrl         string `json:"siteUrl" binding:"required"`
	SiteFavicon     string `json:"siteFavicon"`
	SiteAuthor      string `json:"siteAuthor"`
	SiteAnimateText string `json:"siteAnimateText"`
	SiteAvatar      string `json:"siteAvatar"`
	SiteKeywords    string `json:"siteKeywords"`
	SiteDescription string `json:"siteDescription"`
	SiteCopyright   string `json:"siteCopyright"`
	SiteIcp         string `json:"siteIcp"`
	SiteIcpLink     string `json:"siteIcpLink"`
	GithubUsername  string `json:"githubUsername" binding:"required"`
}
