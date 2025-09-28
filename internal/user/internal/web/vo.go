package web

import (
	"github.com/codepzj/Stellux-Server/internal/user/internal/domain"
	"github.com/samber/lo"
)

type UserVO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	RoleId   int    `json:"role_id"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

type LoginVO struct {
	AccessToken string `json:"access_token"`
}

func (h *UserHandler) UserDomainToVO(user *domain.User) *UserVO {
	return &UserVO{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		RoleId:   user.RoleId,
		Avatar:   user.Avatar,
		Email:    user.Email,
	}
}

func (h *UserHandler) UserDomainToVOList(users []*domain.User) []*UserVO {
	return lo.Map(users, func(user *domain.User, _ int) *UserVO {
		return h.UserDomainToVO(user)
	})
}
