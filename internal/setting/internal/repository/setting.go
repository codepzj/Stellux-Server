package repository

import (
	"context"

	"github.com/codepzj/stellux/server/internal/setting/internal/domain"
	"github.com/codepzj/stellux/server/internal/setting/internal/repository/dao"
)

type ISettingRepository interface {
	Upsert(ctx context.Context, setting domain.Setting) error
	GetSetting(ctx context.Context, key string) (*domain.Setting, error)
}

var _ ISettingRepository = (*SettingRepository)(nil)

func NewSettingRepository(dao dao.ISettingDao) *SettingRepository {
	return &SettingRepository{dao: dao}
}

type SettingRepository struct {
	dao dao.ISettingDao
}

func (r *SettingRepository) Upsert(ctx context.Context, setting domain.Setting) error {
	return r.dao.Upsert(ctx, r.DomainToDao(&setting))
}

func (r *SettingRepository) GetSetting(ctx context.Context, key string) (*domain.Setting, error) {
	setting, err := r.dao.GetSetting(ctx, key)
	if err != nil {
		return nil, err
	}
	return r.DaoToDomain(setting), nil
}

func (r *SettingRepository) DomainToDao(domain *domain.Setting) *dao.Setting {
	return &dao.Setting{
		Key:   domain.Key,
		Value: domain.Value,
	}
}

func (r *SettingRepository) DaoToDomain(dao *dao.Setting) *domain.Setting {
	return &domain.Setting{
		Key:   dao.Key,
		Value: dao.Value,
	}
}