package repository

import (
	"context"

	"github.com/codepzj/stellux/server/internal/setting/internal/domain"
	"github.com/codepzj/stellux/server/internal/setting/internal/repository/dao"
)

type ISettingRepository interface {
	Upsert(ctx context.Context, setting domain.SiteSetting) error
	GetSetting(ctx context.Context, key string) (*domain.SiteSetting, error)
}

var _ ISettingRepository = (*SettingRepository)(nil)

func NewSettingRepository(dao dao.ISettingDao) *SettingRepository {
	return &SettingRepository{dao: dao}
}

type SettingRepository struct {
	dao dao.ISettingDao
}

func (r *SettingRepository) Upsert(ctx context.Context, setting domain.SiteSetting) error {
	return r.dao.Upsert(ctx, r.DomainToDao(&setting))
}

func (r *SettingRepository) GetSetting(ctx context.Context, key string) (*domain.SiteSetting, error) {
	setting, err := r.dao.GetSetting(ctx, key)
	if err != nil {
		return nil, err
	}
	return r.DaoToDomain(setting), nil
}

func (r *SettingRepository) DomainToDao(domain *domain.SiteSetting) *dao.Setting {
	return &dao.Setting{
		Key:   domain.Key,
		Value: dao.SiteConfig{
			SiteTitle:   domain.Value.SiteTitle,
			SiteSubTitle: domain.Value.SiteSubTitle,
			SiteFavicon:  domain.Value.SiteFavicon,
			SiteAvatar:   domain.Value.SiteAvatar,
			SiteKeywords: domain.Value.SiteKeywords,
			SiteDescription: domain.Value.SiteDescription,
			SiteCopyright:   domain.Value.SiteCopyright,
			SiteICP:         domain.Value.SiteICP,
			SiteICPLink:     domain.Value.SiteICPLink,
			GithubUsername:  domain.Value.GithubUsername,
		},
	}
}

func (r *SettingRepository) DaoToDomain(dao *dao.Setting) *domain.SiteSetting {
	return &domain.SiteSetting{
		Key:   dao.Key,
		Value: domain.SiteConfig{
			SiteTitle:   dao.Value.SiteTitle,
			SiteSubTitle: dao.Value.SiteSubTitle,
			SiteFavicon:  dao.Value.SiteFavicon,
			SiteAvatar:   dao.Value.SiteAvatar,
			SiteKeywords: dao.Value.SiteKeywords,
			SiteDescription: dao.Value.SiteDescription,
			SiteCopyright:   dao.Value.SiteCopyright,
			SiteICP:         dao.Value.SiteICP,
			SiteICPLink:     dao.Value.SiteICPLink,
			GithubUsername:  dao.Value.GithubUsername,
		},
	}
}