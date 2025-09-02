package service

import (
	"context"

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

// CreateFriend 创建好友
func (s *FriendService) CreateFriend(ctx context.Context, friend *domain.Friend) error {
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
	return s.repo.Update(ctx, id, friend)
}

// DeleteFriend 删除好友
func (s *FriendService) DeleteFriend(ctx context.Context, id bson.ObjectID) error {
	return s.repo.Delete(ctx, id)
}
