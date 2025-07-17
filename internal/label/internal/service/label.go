package service

import (
	"context"

	"github.com/codepzj/stellux/server/internal/label/internal/domain"
	"github.com/codepzj/stellux/server/internal/label/internal/repository"
)

type ILabelService interface {
	CreateLabel(ctx context.Context, label *domain.Label) error
	UpdateLabel(ctx context.Context, id string, label *domain.Label) error
	DeleteLabel(ctx context.Context, id string) error
	GetLabelById(ctx context.Context, id string) (*domain.Label, error)
	QueryLabelList(ctx context.Context, labelType string, pageNo int64, pageSize int64) ([]*domain.Label, int64, error)
	GetAllLabelsByType(ctx context.Context, labelType string) ([]*domain.Label, error)
	GetAllLabelsWithCount(ctx context.Context) ([]*domain.LabelPostCount, error)
	GetAllTagsLabelWithCount(ctx context.Context) ([]*domain.LabelPostCount, error)
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
	return s.repo.CreateLabel(ctx, label)
}

// UpdateLabel 更新标签
func (s *LabelService) UpdateLabel(ctx context.Context, id string, label *domain.Label) error {
	return s.repo.UpdateLabel(ctx, id, label)
}

// DeleteLabel 删除标签
func (s *LabelService) DeleteLabel(ctx context.Context, id string) error {
	return s.repo.DeleteLabel(ctx, id)
}

// GetLabelById 根据id获取标签
func (s *LabelService) GetLabelById(ctx context.Context, id string) (*domain.Label, error) {
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

// GetAllLabelsWithCount 获取所有分类标签及其文章数量
func (s *LabelService) GetAllLabelsWithCount(ctx context.Context) ([]*domain.LabelPostCount, error) {
	return s.repo.GetCategoryLabelWithCount(ctx)
}

// GetAllTagsLabelWithCount 获取所有标签及其文章数量
func (s *LabelService) GetAllTagsLabelWithCount(ctx context.Context) ([]*domain.LabelPostCount, error) {
	return s.repo.GetTagsLabelWithCount(ctx)
}
