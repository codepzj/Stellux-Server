package service

import (
	"context"
	"errors"
	"regexp"

	"github.com/codepzj/Stellux-Server/internal/friend/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/friend/internal/repository"
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
		return err
	}
	// 唯一性校验：相同站点URL不允许重复
	exists, err := s.repo.ExistsBySiteUrl(ctx, friend.SiteUrl)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("站点URL已存在")
	}
	return s.repo.Create(ctx, friend)
}

// FindFriendList 查询友链列表
func (s *FriendService) FindFriendList(ctx context.Context) ([]*domain.Friend, error) {
	return s.repo.FindAllActive(ctx)
}

// FindAllFriends 查询所有好友
func (s *FriendService) FindAllFriends(ctx context.Context) ([]*domain.Friend, error) {
	return s.repo.FindAll(ctx)
}

// UpdateFriend 更新好友信息
func (s *FriendService) UpdateFriend(ctx context.Context, id bson.ObjectID, friend *domain.Friend) error {
	if err := s.validateFriend(friend); err != nil {
		return err
	}
	// 唯一性校验（排除自身ID）
	exists, err := s.repo.ExistsBySiteUrlExceptID(ctx, friend.SiteUrl, id)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("站点URL已存在")
	}
	return s.repo.Update(ctx, id, friend)
}

// DeleteFriend 删除好友
func (s *FriendService) DeleteFriend(ctx context.Context, id bson.ObjectID) error {
	return s.repo.Delete(ctx, id)
}
