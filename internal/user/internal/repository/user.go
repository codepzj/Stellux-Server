package repository

import (
	"context"

	"github.com/codepzj/stellux/server/internal/user/internal/domain"
	"github.com/codepzj/stellux/server/internal/user/internal/repository/dao"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, id string, password string) error
	Delete(ctx context.Context, id string) error
	FindByPage(ctx context.Context, page *domain.Page) ([]*domain.User, int64, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

var _ IUserRepository = (*UserRepository)(nil)

func NewUserRepository(dao dao.IUserDao) *UserRepository {
	return &UserRepository{dao: dao}
}

type UserRepository struct {
	dao dao.IUserDao
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.dao.Create(ctx, r.UserDomainToUserDO(user))
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := r.dao.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return r.UserDoToUserDomain(user), nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	return r.dao.Update(ctx, user.ID, &dao.User{
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Email:    user.Email,
	})
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id string, password string) error {
	bid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.UpdatePassword(ctx, bid, password)
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	bid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.Delete(ctx, bid)
}

func (r *UserRepository) FindByPage(ctx context.Context, page *domain.Page) ([]*domain.User, int64, error) {
	findOptions := options.Find().SetSkip((page.PageNo - 1) * page.PageSize).SetLimit(page.PageSize).SetSort(bson.M{"role_id": 1})
	users, count, err := r.dao.FindByCondition(ctx, findOptions)
	if err != nil {
		return nil, 0, err
	}
	return r.UserDoToUserDomainList(users), count, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	bid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	user, err := r.dao.GetByID(ctx, bid)
	if err != nil {
		return nil, err
	}
	return r.UserDoToUserDomain(user), nil
}

func (r *UserRepository) UserDomainToUserDO(user *domain.User) *dao.User {
	return &dao.User{
		Username: user.Username,
		Password: user.Password,
		Nickname: user.Nickname,
		RoleId:   user.RoleId,
		Avatar:   user.Avatar,
		Email:    user.Email,
	}
}

func (r *UserRepository) UserDoToUserDomain(user *dao.User) *domain.User {
	return &domain.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Nickname: user.Nickname,
		RoleId:   user.RoleId,
		Avatar:   user.Avatar,
		Email:    user.Email,
	}
}

func (r *UserRepository) UserDoToUserDomainList(users []*dao.User) []*domain.User {
	return lo.Map(users, func(user *dao.User, _ int) *domain.User {
		return r.UserDoToUserDomain(user)
	})
}
