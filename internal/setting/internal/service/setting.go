package service

import (
	"context"

	"github.com/codepzj/stellux/server/internal/setting/internal/domain"
	"github.com/codepzj/stellux/server/internal/setting/internal/repository"
)

type ISettingService interface {
	AdminUpsertSetting(ctx context.Context, setting domain.Setting) error
	GetSetting(ctx context.Context, key string) (*domain.Setting, error)
}

var _ ISettingService = (*SettingService)(nil)

func NewSettingService(repo repository.ISettingRepository) *SettingService {
	return &SettingService{
		repo: repo,
	}
}

type SettingService struct {
	repo repository.ISettingRepository
}

func (s *SettingService) AdminUpsertSetting(ctx context.Context, setting domain.Setting) error {
	return s.repo.Upsert(ctx, setting)
}

func (s *SettingService) GetSetting(ctx context.Context, key string) (*domain.Setting, error) {
	return s.repo.GetSetting(ctx, key)
}