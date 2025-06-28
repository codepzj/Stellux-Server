package web

var WebsiteType = map[int]string{
	0: "大佬",
	1: "技术型",
	2: "设计型",
	3: "生活型",
}

type FriendRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	SiteUrl     string `json:"site_url" binding:"required"`
	AvatarUrl   string `json:"avatar_url" binding:"required"`
	WebsiteType int    `json:"website_type" binding:"min=0,max=3"`
}

type FriendUpdateRequest struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	SiteUrl     string `json:"site_url" binding:"required"`
	AvatarUrl   string `json:"avatar_url" binding:"required"`
	WebsiteType int    `json:"website_type" binding:"min=0,max=3"`
	IsActive    bool   `json:"is_active"`
}
