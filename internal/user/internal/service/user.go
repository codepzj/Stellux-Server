package service

import (
	"context"
	"errors"
	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"

	"github.com/codepzj/Stellux-Server/internal/pkg/logger"
	"github.com/codepzj/Stellux-Server/internal/pkg/utils"
	"github.com/codepzj/Stellux-Server/internal/user/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/user/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IUserService interface {
	CheckUserExist(ctx context.Context, user *domain.User) (bool, string)
	AdminCreate(ctx context.Context, user *domain.User) error
	AdminUpdatePassword(ctx context.Context, id string, oldPassword string, newPassword string) error
	AdminUpdate(ctx context.Context, user *domain.User) error
	AdminDelete(ctx context.Context, id string) error
	GetUserList(ctx context.Context, page *apiwrap.Page) ([]*domain.User, int64, error)
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
		logger.Error("查询用户失败",
			logger.WithError(err),
			logger.WithString("username", user.Username),
		)
		return false, ""
	}

	if u == nil {
		logger.Warn("用户不存在",
			logger.WithString("username", user.Username),
		)
		return false, ""
	}

	passwordMatch := utils.CompareHashAndPassword(u.Password, user.Password)
	if passwordMatch {
		logger.Info("用户验证成功",
			logger.WithString("username", user.Username),
			logger.WithString("userId", u.ID.Hex()),
		)
	} else {
		logger.Warn("密码不匹配",
			logger.WithString("username", user.Username),
		)
	}

	return passwordMatch, u.ID.Hex()
}

// 管理员创建用户
func (s *UserService) AdminCreate(ctx context.Context, user *domain.User) error {
	u, err := s.repo.GetByUsername(ctx, user.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		logger.Error("查询用户失败",
			logger.WithError(err),
			logger.WithString("username", user.Username),
		)
		return err
	}

	if u != nil {
		logger.Warn("用户已存在",
			logger.WithString("username", user.Username),
		)
		return errors.New("用户已存在")
	}

	// 使用bcrypt生成hash密码
	user.Password, err = utils.GenerateHashPassword(user.Password)
	if err != nil {
		logger.Error("密码加密失败",
			logger.WithError(err),
			logger.WithString("username", user.Username),
		)
		return err
	}

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		logger.Error("创建用户失败",
			logger.WithError(err),
			logger.WithString("username", user.Username),
		)
		return err
	}

	logger.Info("创建用户成功",
		logger.WithString("username", user.Username),
		logger.WithString("userId", id.Hex()),
	)

	return nil
}

// 管理员更新用户密码
func (s *UserService) AdminUpdatePassword(ctx context.Context, id string, oldPassword string, newPassword string) error {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("查询用户失败",
			logger.WithError(err),
			logger.WithString("userId", id),
		)
		return err
	}

	if !utils.CompareHashAndPassword(u.Password, oldPassword) {
		logger.Warn("旧密码错误",
			logger.WithString("userId", id),
			logger.WithString("username", u.Username),
		)
		return errors.New("旧密码错误")
	}

	newPassword, err = utils.GenerateHashPassword(newPassword)
	if err != nil {
		logger.Error("密码加密失败",
			logger.WithError(err),
			logger.WithString("userId", id),
		)
		return err
	}

	err = s.repo.UpdatePassword(ctx, id, newPassword)
	if err != nil {
		logger.Error("更新密码失败",
			logger.WithError(err),
			logger.WithString("userId", id),
		)
		return err
	}

	logger.Info("更新用户密码成功",
		logger.WithString("userId", id),
		logger.WithString("username", u.Username),
	)

	return nil
}

// 管理员更新用户
func (s *UserService) AdminUpdate(ctx context.Context, user *domain.User) error {
	err := s.repo.Update(ctx, user)
	if err != nil {
		logger.Error("更新用户失败",
			logger.WithError(err),
			logger.WithString("userId", user.ID.Hex()),
		)
		return err
	}

	logger.Info("更新用户成功",
		logger.WithString("userId", user.ID.Hex()),
		logger.WithString("username", user.Username),
	)

	return nil
}

// 管理员删除用户
func (s *UserService) AdminDelete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		logger.Error("删除用户失败",
			logger.WithError(err),
			logger.WithString("userId", id),
		)
		return err
	}

	logger.Info("删除用户成功",
		logger.WithString("userId", id),
	)

	return nil
}

// 获取用户列表
func (s *UserService) GetUserList(ctx context.Context, page *apiwrap.Page) ([]*domain.User, int64, error) {
	users, total, err := s.repo.FindByPage(ctx, page)
	if err != nil {
		logger.Error("查询用户列表失败",
			logger.WithError(err),
			logger.WithInt("pageNo", int(page.PageNo)),
		)
		return nil, 0, err
	}
	return users, total, nil
}

// 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("查询用户信息失败",
			logger.WithError(err),
			logger.WithString("userId", id),
		)
		return nil, err
	}
	return user, nil
}
