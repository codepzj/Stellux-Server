package service

import (
	"context"
	"errors"

	"github.com/codepzj/Stellux-Server/internal/label/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/label/internal/repository"
	"gorm.io/gorm"
)

type ILabelService interface {
	CreateLabel(ctx context.Context, label *domain.Label) error
	UpdateLabel(ctx context.Context, id uint, label *domain.Label) error
	DeleteLabel(ctx context.Context, id uint) error
	GetLabelById(ctx context.Context, id uint) (*domain.Label, error)
	QueryLabelList(ctx context.Context, labelType string, pageNo int64, pageSize int64) ([]*domain.Label, int64, error)
	GetAllLabelsByType(ctx context.Context, labelType string) ([]*domain.Label, error)
}

var _ ILabelService = (*LabelService)(nil)

func NewLabelService(repo repository.ILabelRepository) *LabelService {
	return &LabelService{
		repo: repo,
	}
}

type LabelService struct {
	repo repository.ILabelRepository
}

// CreateLabel 创建标签
func (s *LabelService) CreateLabel(ctx context.Context, label *domain.Label) error {
	// 检查标签是否存在
	existLabel, err := s.repo.GetLabelByName(ctx, label.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existLabel != nil {
		return errors.New("标签已存在")
	}
	return s.repo.CreateLabel(ctx, label)
}

// UpdateLabel 更新标签
func (s *LabelService) UpdateLabel(ctx context.Context, id uint, label *domain.Label) error {
	// 检查标签是否存在
	existLabel, err := s.repo.GetLabelByName(ctx, label.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existLabel != nil && existLabel.ID != id {
		return errors.New("标签已存在")
	}
	return s.repo.UpdateLabel(ctx, id, label)
}

// DeleteLabel 删除标签
func (s *LabelService) DeleteLabel(ctx context.Context, id uint) error {
	return s.repo.DeleteLabel(ctx, id)
}

// GetLabelById 根据id获取标签
func (s *LabelService) GetLabelById(ctx context.Context, id uint) (*domain.Label, error) {
	return s.repo.GetLabelById(ctx, id)
}

// QueryLabelList 分页查询标签
func (s *LabelService) QueryLabelList(ctx context.Context, labelType string, pageNo int64, pageSize int64) ([]*domain.Label, int64, error) {
	return s.repo.QueryLabelList(ctx, labelType, pageNo, pageSize)
}

// GetAllLabelsByType 获取所有标签
func (s *LabelService) GetAllLabelsByType(ctx context.Context, labelType string) ([]*domain.Label, error) {
	return s.repo.GetAllLabelsByType(ctx, labelType)
}

