package service

import (
	"context"
	"errors"

	"github.com/codepzj/stellux/server/internal/pkg/utils"
	"github.com/codepzj/stellux/server/internal/user/internal/domain"
	"github.com/codepzj/stellux/server/internal/user/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IUserService interface {
	CheckUserExist(ctx context.Context, user *domain.User) (bool, string)
	AdminCreate(ctx context.Context, user *domain.User) error
	AdminUpdatePassword(ctx context.Context, id string, oldPassword string, newPassword string) error
	AdminUpdate(ctx context.Context, user *domain.User) error
	AdminDelete(ctx context.Context, id string) error
	GetUserList(ctx context.Context, page *domain.Page) ([]*domain.User, int64, error)
	GetUserInfo(ctx context.Context, id string) (*domain.User, error)
}

var _ IUserService = (*UserService)(nil)

func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

type UserService struct {
	repo repository.IUserRepository
}

func (s *UserService) CheckUserExist(ctx context.Context, user *domain.User) (bool, string) {
	u, err := s.repo.GetByUsername(ctx, user.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		return false, ""
	}
	if u == nil {
		return false, ""
	}
	return utils.CompareHashAndPassword(u.Password, user.Password), u.ID.Hex()
}

// 管理员创建用户
func (s *UserService) AdminCreate(ctx context.Context, user *domain.User) error {
	u, err := s.repo.GetByUsername(ctx, user.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if u != nil {
		return errors.New("用户已存在")
	}
	// 使用bcrypt生成hash密码
	user.Password, err = utils.GenerateHashPassword(user.Password)
	if err != nil {
		return err
	}
	return s.repo.Create(ctx, user)
}

// 管理员更新用户密码
func (s *UserService) AdminUpdatePassword(ctx context.Context, id string, oldPassword string, newPassword string) error {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !utils.CompareHashAndPassword(u.Password, oldPassword) {
		return errors.New("旧密码错误")
	}
	newPassword, err = utils.GenerateHashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, id, newPassword)
}

// 管理员更新用户
func (s *UserService) AdminUpdate(ctx context.Context, user *domain.User) error {
	return s.repo.Update(ctx, user)
}

// 管理员删除用户
func (s *UserService) AdminDelete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// 获取用户列表
func (s *UserService) GetUserList(ctx context.Context, page *domain.Page) ([]*domain.User, int64, error) {
	return s.repo.FindByPage(ctx, page)
}

// 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}
