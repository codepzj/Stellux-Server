package repository

import (
	"context"

	"github.com/codepzj/stellux/server/internal/friend/internal/domain"
	"github.com/codepzj/stellux/server/internal/friend/internal/repository/dao"
	"github.com/samber/lo"
)

type IFriendRepository interface {
	Create(ctx context.Context, friend *domain.Friend) error
	FindAll(ctx context.Context) ([]*domain.Friend, error)
}

var _ IFriendRepository = (*FriendRepository)(nil)

func NewFriendRepository(dao dao.IFriendDao) *FriendRepository {
	return &FriendRepository{dao: dao}
}

type FriendRepository struct {
	dao dao.IFriendDao
}

func (r *FriendRepository) Create(ctx context.Context, friend *domain.Friend) error {
	return r.dao.Create(ctx, r.FriendDomainToDao(friend))
}

func (r *FriendRepository) FindAll(ctx context.Context) ([]*domain.Friend, error) {
	friends, err := r.dao.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return r.FriendDaoToDomainList(friends), nil
}

func (r *FriendRepository) FriendDomainToDao(friend *domain.Friend) *dao.Friend {
	return &dao.Friend{
		Name:        friend.Name,
		Description: friend.Description,
		SiteUrl:     friend.SiteUrl,
		AvatarUrl:   friend.AvatarUrl,
		WebsiteType: friend.WebsiteType,
	}
}

func (r *FriendRepository) FriendDaoToDomain(friend *dao.Friend) *domain.Friend {
	return &domain.Friend{
		ID:          friend.ID,
		Name:        friend.Name,
		Description: friend.Description,
		SiteUrl:     friend.SiteUrl,
		AvatarUrl:   friend.AvatarUrl,
		WebsiteType: friend.WebsiteType,
		IsActive:    friend.IsActive,
	}
}

func (r *FriendRepository) FriendDaoToDomainList(friends []*dao.Friend) []*domain.Friend {
	return lo.Map(friends, func(friend *dao.Friend, _ int) *domain.Friend {
		return r.FriendDaoToDomain(friend)
	})
}
