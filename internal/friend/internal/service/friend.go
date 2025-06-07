package service

import (
	"context"

	"github.com/codepzj/stellux/server/internal/friend/internal/domain"
	"github.com/codepzj/stellux/server/internal/friend/internal/repository"
)


type IFriendService interface {
	Create(ctx context.Context, friend *domain.Friend) error
	FindAll(ctx context.Context) ([]*domain.Friend, error)
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

func (s *FriendService) Create(ctx context.Context, friend *domain.Friend) error {
	return s.repo.Create(ctx, friend)
}

func (s *FriendService) FindAll(ctx context.Context) ([]*domain.Friend, error) {
	return s.repo.FindAll(ctx)
}