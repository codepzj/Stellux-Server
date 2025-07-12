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

func (s *LabelService) CreateLabel(ctx context.Context, label *domain.Label) error {
	return s.repo.CreateLabel(ctx, label)
}

func (s *LabelService) UpdateLabel(ctx context.Context, id string, label *domain.Label) error {
	return s.repo.UpdateLabel(ctx, id, label)
}

func (s *LabelService) DeleteLabel(ctx context.Context, id string) error {
	return s.repo.DeleteLabel(ctx, id)
}

func (s *LabelService) GetLabelById(ctx context.Context, id string) (*domain.Label, error) {
	return s.repo.GetLabelById(ctx, id)
}

func (s *LabelService) QueryLabelList(ctx context.Context, labelType string, pageNo int64, pageSize int64) ([]*domain.Label, int64, error) {
	return s.repo.QueryLabelList(ctx, labelType, pageNo, pageSize)
}

func (s *LabelService) GetAllLabelsByType(ctx context.Context, labelType string) ([]*domain.Label, error) {
	return s.repo.GetAllLabelsByType(ctx, labelType)
}

func (s *LabelService) GetAllLabelsWithCount(ctx context.Context) ([]*domain.LabelPostCount, error) {
	return s.repo.GetCategoryLabelWithCount(ctx)
}
