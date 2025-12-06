package web

// ConfigDto 网站配置DTO
type ConfigDto struct {
	Type    string  `json:"type" binding:"required,oneof=home about seo"`
	Content Content `json:"content" binding:"required"`
}

// Content 网站内容
type Content struct {
	// 通用配置
	Title       string `json:"title"`
	Description string `json:"description"`

	// 主页配置
	Avatar           string   `json:"avatar,omitempty"`
	Name             string   `json:"name,omitempty"`
	Bio              string   `json:"bio,omitempty"`
	Github           string   `json:"github,omitempty"`
	Blog             string   `json:"blog,omitempty"`
	Location         string   `json:"location,omitempty"`
	TechStacks       []string `json:"tech_stacks,omitempty"`
	Repositories     []Repo   `json:"repositories,omitempty"`
	Quote            string   `json:"quote,omitempty"`
	ShowRecentPosts  bool     `json:"show_recent_posts,omitempty"`
	RecentPostsCount int      `json:"recent_posts_count,omitempty" binding:"omitempty,gte=1,lte=50"`

	// About页面配置
	Skills     []Skill    `json:"skills,omitempty"`
	Timeline   []Timeline `json:"timeline,omitempty"`
	Interests  []string   `json:"interests,omitempty"`
	FocusItems []string   `json:"focus_items,omitempty"`

	// SEO配置
	SEOTitle        string   `json:"seo_title,omitempty"`
	SEOKeywords     []string `json:"seo_keywords,omitempty"`
	SEODescription  string   `json:"seo_description,omitempty"`
	RobotsMeta      string   `json:"robots_meta,omitempty"`
	CanonicalURL    string   `json:"canonical_url,omitempty"`
	OGTitle         string   `json:"og_title,omitempty"`
	OGDescription   string   `json:"og_description,omitempty"`
	OGImage         string   `json:"og_image,omitempty"`
	TwitterCard     string   `json:"twitter_card,omitempty"`
}

// Repo 开源项目
type Repo struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
	Desc string `json:"desc"`
}

// Skill 技能
type Skill struct {
	Category string   `json:"category" binding:"required"`
	Items    []string `json:"items" binding:"required"`
}

// Timeline 时间线项目
type Timeline struct {
	Year  string `json:"year" binding:"required"`
	Title string `json:"title" binding:"required"`
	Desc  string `json:"desc"`
}

// ConfigUpdateDto 更新网站配置DTO
type ConfigUpdateDto struct {
	ID      string  `json:"id" binding:"required"`
	Type    string  `json:"type" binding:"required,oneof=home about seo"`
	Content Content `json:"content" binding:"required"`
}

// ConfigIdRequest ID请求
type ConfigIdRequest struct {
	ID string `uri:"id" binding:"required"`
}

// ConfigTypeRequest 类型请求
type ConfigTypeRequest struct {
	Type string `uri:"type" binding:"required,oneof=home about seo"`
}
