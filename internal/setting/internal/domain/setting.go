package domain

type Setting struct {
	Key       string
	Value     any
}

// 站点配置
type SiteSetting struct {
	Key       string
	Value     SiteConfig
}

type SiteConfig struct {
	SiteTitle   string // 网站标题
	SiteSubTitle string // 网站副标题
	SiteFavicon  string // 网站图标
	SiteAvatar   string // 网站头像
	SiteKeywords string // 网站关键词
	SiteDescription string // 网站描述
	SiteCopyright   string // 网站版权
	SiteICP         string // 网站备案号
	SiteICPLink     string // 网站备案号链接
	GithubUsername  string // Github用户名
}