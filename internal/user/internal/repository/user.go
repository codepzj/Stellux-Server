package repository

import (
	"context"

	"github.com/codepzj/Stellux-Server/internal/user/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/user/internal/repository/dao"
	"github.com/codepzj/gokit/slice"

	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, id uint, password string) error
	Delete(ctx context.Context, id uint) error
	FindByPage(ctx context.Context, page *domain.Page) ([]*domain.User, int64, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
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

func (r *UserRepository) UpdatePassword(ctx context.Context, id uint, password string) error {
	return r.dao.UpdatePassword(ctx, id, password)
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.dao.Delete(ctx, id)
}

// 分页查询用户
func (r *UserRepository) FindByPage(ctx context.Context, page *domain.Page) ([]*domain.User, int64, error) {
	users, err := r.dao.Find(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Offset(int((page.PageNo - 1) * page.PageSize)).Limit(int(page.PageSize))
	})
	if err != nil {
		return nil, 0, err
	}
	count, err := r.dao.GetAllCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return r.UserDoToUserDomainList(users), count, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	user, err := r.dao.GetByID(ctx, id)
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
	return slice.Map(users, func(user *dao.User) *domain.User {
		return r.UserDoToUserDomain(user)
	})
}
