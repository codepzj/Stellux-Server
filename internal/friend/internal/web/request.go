package web

type FriendRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	SiteUrl     string `json:"site_url" binding:"required"`
	AvatarUrl   string `json:"avatar_url" binding:"required"`
	WebsiteType string `json:"website_type" binding:"required"`
}

type FriendUpdateRequest struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	SiteUrl     string `json:"site_url" binding:"required"`
	AvatarUrl   string `json:"avatar_url" binding:"required"`
	WebsiteType string `json:"website_type" binding:"required"`
	IsActive    bool   `json:"is_active"`
}
