package service

import (
	"context"

	"github.com/codepzj/stellux/server/internal/label/internal/domain"
	"github.com/codepzj/stellux/server/internal/label/internal/repository"
)

type ILabelService interface {
	Create(ctx context.Context, label *domain.Label) error
	Update(ctx context.Context, id string, label *domain.Label) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*domain.Label, error)
	QueryLabelList(ctx context.Context, labelType string, pageNo int64, pageSize int64) ([]*domain.Label, int64, error)
	QueryAllByType(ctx context.Context, labelType string) ([]*domain.Label, error)
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

func (s *LabelService) Create(ctx context.Context, label *domain.Label) error {
	return s.repo.Create(ctx, label)
}

func (s *LabelService) Update(ctx context.Context, id string, label *domain.Label) error {
	return s.repo.Update(ctx, id, label)
}

func (s *LabelService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *LabelService) Get(ctx context.Context, id string) (*domain.Label, error) {
	return s.repo.Get(ctx, id)
}

func (s *LabelService) QueryLabelList(ctx context.Context, labelType string, pageNo int64, pageSize int64) ([]*domain.Label, int64, error) {
	return s.repo.GetList(ctx, labelType, pageNo, pageSize)
}

func (s *LabelService) QueryAllByType(ctx context.Context, labelType string) ([]*domain.Label, error) {
	return s.repo.GetAllByType(ctx, labelType)
}
