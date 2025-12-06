package service

import (
	"context"
	"errors"
	"time"

	"github.com/codepzj/Stellux-Server/internal/config/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/config/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IConfigService interface {
	CreateConfig(ctx context.Context, config *domain.Config) error
	UpdateConfig(ctx context.Context, config *domain.Config) error
	GetConfigByID(ctx context.Context, id bson.ObjectID) (*domain.Config, error)
	GetConfigByType(ctx context.Context, configType string) (*domain.Config, error)
	ListConfigs(ctx context.Context) ([]*domain.Config, error)
	DeleteConfig(ctx context.Context, id bson.ObjectID) error
}

var _ IConfigService = (*ConfigService)(nil)

func NewConfigService(repo repository.IConfigRepository) *ConfigService {
	return &ConfigService{
		repo: repo,
	}
}

type ConfigService struct {
	repo repository.IConfigRepository
}

// CreateConfig 创建网站配置
func (s *ConfigService) CreateConfig(ctx context.Context, config *domain.Config) error {
	logger.Info("创建网站配置", logger.WithString("type", config.Type))

	// 检查是否已存在相同类型的配置
	if existing, err := s.repo.GetByType(ctx, config.Type); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		logger.Error("检查配置类型是否存在失败", logger.WithError(err), logger.WithString("type", config.Type))
		return err
	} else if existing != nil {
		logger.Warn("配置类型已存在", logger.WithString("type", config.Type))
		return errors.New("config type already exists")
	}

	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, config); err != nil {
		logger.Error("创建网站配置失败", logger.WithError(err), logger.WithString("type", config.Type))
		return err
	}

	logger.Info("网站配置创建成功", logger.WithString("type", config.Type))
	return nil
}

// GetConfigByID 根据ID获取网站配置
func (s *ConfigService) GetConfigByID(ctx context.Context, id bson.ObjectID) (*domain.Config, error) {
	logger.Info("根据ID获取网站配置", logger.WithString("id", id.Hex()))

	config, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("根据ID获取网站配置失败", logger.WithError(err), logger.WithString("id", id.Hex()))
		return nil, err
	}

	logger.Info("根据ID获取网站配置成功", logger.WithString("id", id.Hex()))
	return config, nil
}

// UpdateConfig 更新网站配置
func (s *ConfigService) UpdateConfig(ctx context.Context, config *domain.Config) error {
	logger.Info("更新页面配置", logger.WithString("id", config.Id.Hex()), logger.WithString("type", config.Type))

	// 检查配置是否存在
	if existing, err := s.repo.GetByType(ctx, config.Type); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		logger.Error("检查配置是否存在失败", logger.WithError(err), logger.WithString("type", config.Type))
		return err
	} else if existing != nil && existing.Id != config.Id {
		logger.Warn("配置类型冲突", logger.WithString("type", config.Type))
		return errors.New("config type already exists")
	}

	config.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, config); err != nil {
		logger.Error("更新页面配置失败", logger.WithError(err), logger.WithString("id", config.Id.Hex()))
		return err
	}

	logger.Info("页面配置更新成功", logger.WithString("id", config.Id.Hex()))
	return nil
}

// GetConfigByType 根据类型获取网站配置
func (s *ConfigService) GetConfigByType(ctx context.Context, configType string) (*domain.Config, error) {
	logger.Info("获取页面配置", logger.WithString("type", configType))

	config, err := s.repo.GetByType(ctx, configType)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Warn("页面配置不存在", logger.WithString("type", configType))
			return nil, errors.New("config not found")
		}
		logger.Error("获取页面配置失败", logger.WithError(err), logger.WithString("type", configType))
		return nil, err
	}

	logger.Info("页面配置获取成功", logger.WithString("type", configType))
	return config, nil
}

// ListConfigs 获取所有网站配置
func (s *ConfigService) ListConfigs(ctx context.Context) ([]*domain.Config, error) {
	logger.Info("获取所有网站配置")

	configs, err := s.repo.List(ctx)
	if err != nil {
		logger.Error("获取所有网站配置失败", logger.WithError(err))
		return nil, err
	}

	logger.Info("获取所有网站配置成功", logger.WithInt("count", len(configs)))
	return configs, nil
}

// DeleteConfig 删除网站配置
func (s *ConfigService) DeleteConfig(ctx context.Context, id bson.ObjectID) error {
	logger.Info("删除网站配置", logger.WithString("id", id.Hex()))

	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Error("删除网站配置失败", logger.WithError(err), logger.WithString("id", id.Hex()))
		return err
	}

	logger.Info("网站配置删除成功", logger.WithString("id", id.Hex()))
	return nil
}
