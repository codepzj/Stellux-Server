package repository

import (
	"context"
	"encoding/json"

	"github.com/codepzj/Stellux-Server/internal/config/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/config/internal/repository/dao"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IConfigRepository interface {
	Create(ctx context.Context, config *domain.Config) error
	Update(ctx context.Context, config *domain.Config) error
	GetByID(ctx context.Context, id bson.ObjectID) (*domain.Config, error)
	GetByType(ctx context.Context, configType string) (*domain.Config, error)
	List(ctx context.Context) ([]*domain.Config, error)
	Delete(ctx context.Context, id bson.ObjectID) error
}

var _ IConfigRepository = (*ConfigRepository)(nil)

func NewConfigRepository(dao dao.IConfigDao) *ConfigRepository {
	return &ConfigRepository{dao: dao}
}

type ConfigRepository struct {
	dao dao.IConfigDao
}

// Create 创建网站配置
func (r *ConfigRepository) Create(ctx context.Context, config *domain.Config) error {
	daoConfig := r.ConfigDomainToDao(config)
	return r.dao.Create(ctx, daoConfig)
}

// Update 更新网站配置
func (r *ConfigRepository) Update(ctx context.Context, config *domain.Config) error {
	daoConfig := r.ConfigDomainToDaoUpdate(config)
	return r.dao.Update(ctx, config.Id, daoConfig)
}

// GetByID 根据ID获取网站配置
func (r *ConfigRepository) GetByID(ctx context.Context, id bson.ObjectID) (*domain.Config, error) {
	daoConfig, err := r.dao.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.ConfigDaoToDomain(daoConfig), nil
}

// GetByType 根据类型获取网站配置
func (r *ConfigRepository) GetByType(ctx context.Context, configType string) (*domain.Config, error) {
	daoConfig, err := r.dao.GetByType(ctx, configType)
	if err != nil {
		return nil, err
	}
	return r.ConfigDaoToDomain(daoConfig), nil
}

// GetActiveByType 获取激活状态的页面配置
// List 获取所有网站配置
func (r *ConfigRepository) List(ctx context.Context) ([]*domain.Config, error) {
	daoConfigs, err := r.dao.List(ctx)
	if err != nil {
		return nil, err
	}

	var configs []*domain.Config
	for _, daoConfig := range daoConfigs {
		configs = append(configs, r.ConfigDaoToDomain(daoConfig))
	}
	return configs, nil
}

// Delete 删除网站配置
func (r *ConfigRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	return r.dao.Delete(ctx, id)
}

// ConfigDomainToDao 域模型转DAO模型
func (r *ConfigRepository) ConfigDomainToDao(domainConfig *domain.Config) *dao.Config {
	contentBytes, _ := json.Marshal(domainConfig.Content)
	var contentMap map[string]interface{}
	json.Unmarshal(contentBytes, &contentMap)

	return &dao.Config{
		ID:        domainConfig.Id,
		CreatedAt: domainConfig.CreatedAt,
		UpdatedAt: domainConfig.UpdatedAt,
		Type:      domainConfig.Type,
		Content:   contentMap,
	}
}

// ConfigDomainToDaoUpdate 域模型转DAO更新模型
func (r *ConfigRepository) ConfigDomainToDaoUpdate(domainConfig *domain.Config) *dao.ConfigUpdate {
	contentBytes, _ := json.Marshal(domainConfig.Content)
	var contentMap map[string]interface{}
	json.Unmarshal(contentBytes, &contentMap)

	return &dao.ConfigUpdate{
		Type:    domainConfig.Type,
		Content: contentMap,
	}
}

// ConfigDaoToDomain DAO模型转域模型
func (r *ConfigRepository) ConfigDaoToDomain(daoConfig *dao.Config) *domain.Config {
	contentBytes, _ := json.Marshal(daoConfig.Content)
	var content domain.Content
	json.Unmarshal(contentBytes, &content)

	return &domain.Config{
		Id:        daoConfig.ID,
		CreatedAt: daoConfig.CreatedAt,
		UpdatedAt: daoConfig.UpdatedAt,
		Type:      daoConfig.Type,
		Content:   content,
	}
}
