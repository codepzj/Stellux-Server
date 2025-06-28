package web

import (
	"github.com/codepzj/stellux/server/internal/friend/internal/domain"
	"github.com/samber/lo"
)

type FriendVO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SiteUrl     string `json:"site_url"`
	AvatarUrl   string `json:"avatar_url"`
	WebsiteType int    `json:"website_type"`
	IsActive    bool   `json:"is_active"`
}

// FriendShowVO 友链展示VO
type FriendShowVO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SiteUrl     string `json:"site_url"`
	AvatarUrl   string `json:"avatar_url"`
	WebsiteType int    `json:"website_type"`
}

func FriendDomainToVO(friend *domain.Friend) *FriendVO {
	return &FriendVO{
		ID:          friend.ID.Hex(),
		Name:        friend.Name,
		Description: friend.Description,
		SiteUrl:     friend.SiteUrl,
		AvatarUrl:   friend.AvatarUrl,
		WebsiteType: friend.WebsiteType,
		IsActive:    friend.IsActive,
	}
}

func FriendDomainToVOList(friends []*domain.Friend) []*FriendVO {
	return lo.Map(friends, func(friend *domain.Friend, _ int) *FriendVO {
		return FriendDomainToVO(friend)
	})
}

func FriendDomainToShowVOList(friends []*domain.Friend) []*FriendShowVO {
	return lo.Map(friends, func(friend *domain.Friend, _ int) *FriendShowVO {
		return &FriendShowVO{
			ID:          friend.ID.Hex(),
			Name:        friend.Name,
			Description: friend.Description,
			SiteUrl:     friend.SiteUrl,
			AvatarUrl:   friend.AvatarUrl,
			WebsiteType: friend.WebsiteType,
		}
	})
}
