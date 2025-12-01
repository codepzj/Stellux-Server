package service

import (
	"context"
	"errors"
	"regexp"

	"github.com/codepzj/Stellux-Server/internal/friend/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/friend/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IFriendService interface {
	CreateFriend(ctx context.Context, friend *domain.Friend) error
	FindFriendList(ctx context.Context) ([]*domain.Friend, error)
	FindAllFriends(ctx context.Context) ([]*domain.Friend, error)
	UpdateFriend(ctx context.Context, id bson.ObjectID, friend *domain.Friend) error
	DeleteFriend(ctx context.Context, id bson.ObjectID) error
}

var _ IFriendService = (*FriendService)(nil)

func NewFriendService(repo repository.IFriendRepository) *FriendService {
	return &FriendService{
		repo: repo,
	}
}

type FriendService struct {
	repo repository.IFriendRepository
}

// validateFriend 使用正则校验URL等字段
func (s *FriendService) validateFriend(friend *domain.Friend) error {
	// 放宽URL校验，只要以http或https开头即可
	urlRegex := regexp.MustCompile(`^https?://.+`)
	if friend.SiteUrl == "" || !urlRegex.MatchString(friend.SiteUrl) {
		return errors.New("站点URL格式不正确")
	}
	if friend.AvatarUrl != "" && !urlRegex.MatchString(friend.AvatarUrl) {
		return errors.New("头像URL格式不正确")
	}
	if friend.Name == "" {
		return errors.New("名称不能为空")
	}
	if friend.Description == "" {
		return errors.New("描述不能为空")
	}
	return nil
}

// CreateFriend 创建好友
func (s *FriendService) CreateFriend(ctx context.Context, friend *domain.Friend) error {
	if err := s.validateFriend(friend); err != nil {
		logger.Warn("友链校验失败",
			logger.WithError(err),
			logger.WithString("name", friend.Name),
		)
		return err
	}

	// 唯一性校验：相同站点URL不允许重复
	exists, err := s.repo.ExistsBySiteUrl(ctx, friend.SiteUrl)
	if err != nil {
		logger.Error("查询站点URL失败",
			logger.WithError(err),
			logger.WithString("siteUrl", friend.SiteUrl),
		)
		return err
	}

	if exists {
		logger.Warn("站点URL已存在",
			logger.WithString("siteUrl", friend.SiteUrl),
		)
		return errors.New("站点URL已存在")
	}

	err = s.repo.Create(ctx, friend)
	if err != nil {
		logger.Error("创建友链失败",
			logger.WithError(err),
			logger.WithString("name", friend.Name),
		)
		return err
	}

	logger.Info("创建友链成功",
		logger.WithString("name", friend.Name),
		logger.WithString("siteUrl", friend.SiteUrl),
	)

	return nil
}

// FindFriendList 查询友链列表
func (s *FriendService) FindFriendList(ctx context.Context) ([]*domain.Friend, error) {
	logger.Info("查询活跃友链列表",
		logger.WithString("method", "FindFriendList"),
	)

	friends, err := s.repo.FindAllActive(ctx)
	if err != nil {
		logger.Error("查询友链列表失败",
			logger.WithError(err),
		)
		return nil, err
	}

	return friends, nil
}

// FindAllFriends 查询所有好友
func (s *FriendService) FindAllFriends(ctx context.Context) ([]*domain.Friend, error) {
	logger.Info("查询所有友链",
		logger.WithString("method", "FindAllFriends"),
	)

	friends, err := s.repo.FindAll(ctx)
	if err != nil {
		logger.Error("查询友链失败",
			logger.WithError(err),
		)
		return nil, err
	}

	return friends, nil
}

// UpdateFriend 更新好友信息
func (s *FriendService) UpdateFriend(ctx context.Context, id bson.ObjectID, friend *domain.Friend) error {
	if err := s.validateFriend(friend); err != nil {
		logger.Warn("友链校验失败",
			logger.WithError(err),
			logger.WithString("friendId", id.Hex()),
		)
		return err
	}

	// 唯一性校验（排除自身ID）
	exists, err := s.repo.ExistsBySiteUrlExceptID(ctx, friend.SiteUrl, id)
	if err != nil {
		logger.Error("查询站点URL失败",
			logger.WithError(err),
			logger.WithString("siteUrl", friend.SiteUrl),
		)
		return err
	}

	if exists {
		logger.Warn("站点URL已被其他友链使用",
			logger.WithString("siteUrl", friend.SiteUrl),
		)
		return errors.New("站点URL已存在")
	}

	err = s.repo.Update(ctx, id, friend)
	if err != nil {
		logger.Error("更新友链失败",
			logger.WithError(err),
			logger.WithString("friendId", id.Hex()),
		)
		return err
	}

	logger.Info("更新友链成功",
		logger.WithString("friendId", id.Hex()),
		logger.WithString("name", friend.Name),
	)

	return nil
}

// DeleteFriend 删除好友
func (s *FriendService) DeleteFriend(ctx context.Context, id bson.ObjectID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		logger.Error("删除友链失败",
			logger.WithError(err),
			logger.WithString("friendId", id.Hex()),
		)
		return err
	}

	logger.Info("删除友链成功",
		logger.WithString("friendId", id.Hex()),
	)

	return nil
}
