package repository

import (
	"context"

	"github.com/codepzj/Stellux-Server/internal/friend/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/friend/internal/repository/dao"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IFriendRepository interface {
	Create(ctx context.Context, friend *domain.Friend) error
	ExistsBySiteUrl(ctx context.Context, siteUrl string) (bool, error)
	ExistsBySiteUrlExceptID(ctx context.Context, siteUrl string, excludeID bson.ObjectID) (bool, error)
	FindAllActive(ctx context.Context) ([]*domain.Friend, error)
	FindAll(ctx context.Context) ([]*domain.Friend, error)
	Update(ctx context.Context, id bson.ObjectID, friend *domain.Friend) error
	Delete(ctx context.Context, id bson.ObjectID) error
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

func (r *FriendRepository) ExistsBySiteUrl(ctx context.Context, siteUrl string) (bool, error) {
	return r.dao.ExistsBySiteUrl(ctx, siteUrl)
}

func (r *FriendRepository) ExistsBySiteUrlExceptID(ctx context.Context, siteUrl string, excludeID bson.ObjectID) (bool, error) {
	return r.dao.ExistsBySiteUrlExceptID(ctx, siteUrl, excludeID)
}

// FindAllActive 查询所有活跃的友链
func (r *FriendRepository) FindAllActive(ctx context.Context) ([]*domain.Friend, error) {
	friends, err := r.dao.FindAllActive(ctx)
	if err != nil {
		return nil, err
	}
	return r.FriendDaoToDomainList(friends), nil
}

// FindAll 查询所有友链
func (r *FriendRepository) FindAll(ctx context.Context) ([]*domain.Friend, error) {
	friends, err := r.dao.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return r.FriendDaoToDomainList(friends), nil
}

func (r *FriendRepository) Update(ctx context.Context, id bson.ObjectID, friend *domain.Friend) error {
	return r.dao.Update(ctx, id, r.FriendDomainToDao(friend))
}

func (r *FriendRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	return r.dao.Delete(ctx, id)
}

func (r *FriendRepository) FriendDomainToDao(friend *domain.Friend) *dao.Friend {
	return &dao.Friend{
		Name:        friend.Name,
		Description: friend.Description,
		SiteUrl:     friend.SiteUrl,
		AvatarUrl:   friend.AvatarUrl,
		WebsiteType: friend.WebsiteType,
		IsActive:    friend.IsActive,
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
