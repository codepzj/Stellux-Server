package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/codepzj/Stellux-Server/internal/pkg/utils"
	"github.com/codepzj/Stellux-Server/internal/user/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/user/internal/repository"
	"gorm.io/gorm"
)

type IUserService interface {
	VerifyUser(ctx context.Context, user *domain.User) (*domain.User, error)
	AdminCreate(ctx context.Context, user *domain.User) error
	AdminUpdatePassword(ctx context.Context, id uint, oldPassword string, newPassword string) error
	AdminUpdate(ctx context.Context, user *domain.User) error
	AdminDelete(ctx context.Context, id uint) error
	GetUserList(ctx context.Context, page *domain.Page) ([]*domain.User, int64, error)
	GetUserInfo(ctx context.Context, id uint) (*domain.User, error)
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

func (s *UserService) VerifyUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	u, err := s.repo.GetByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("用户不存在")
	}
	has := utils.CompareHashAndPassword(u.Password, user.Password)
	if !has {
		return nil, errors.New("用户名或密码错误")
	}
	return u, nil
}

// 管理员创建用户
func (s *UserService) AdminCreate(ctx context.Context, user *domain.User) error {
	_, err := s.repo.GetByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	// 使用bcrypt生成hash密码
	user.Password, err = utils.GenerateHashPassword(user.Password)
	if err != nil {
		return err
	}
	return s.repo.Create(ctx, user)
}

// 管理员更新用户密码
func (s *UserService) AdminUpdatePassword(ctx context.Context, id uint, oldPassword string, newPassword string) error {
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
func (s *UserService) AdminDelete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// 获取用户列表
func (s *UserService) GetUserList(ctx context.Context, page *domain.Page) ([]*domain.User, int64, error) {
	return s.repo.FindByPage(ctx, page)
}

// 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, id uint) (*domain.User, error) {
	fmt.Println("id", id)
	return s.repo.GetByID(ctx, id)
}
