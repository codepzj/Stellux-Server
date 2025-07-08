package domain

type Setting struct {
	Key   string
	Value any
}

// 站点配置
type SiteSetting struct {
	Key   string
	Value SiteConfig
}

type SiteConfig struct {
	SiteTitle       string // 网站标题(seo)
	SiteSubtitle    string // 网站副标题(seo)
	SiteUrl         string // 网站地址(seo)
	SiteFavicon     string // 网站图标(seo)
	SiteAuthor      string // 网站作者(首屏打招呼)
	SiteAnimateText string // 网站打字机文本(首屏打招呼介绍)
	SiteAvatar      string // 网站头像(首屏头像)
	SiteKeywords    string // 网站关键词(seo)
	SiteDescription string // 网站描述(seo)
	SiteCopyright   string // 网站版权(seo)
	SiteIcp         string // 网站备案号(seo)
	SiteIcpLink     string // 网站备案号链接
	GithubUsername  string // Github用户名
}
