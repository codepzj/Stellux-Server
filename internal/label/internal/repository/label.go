package repository

import (
	"context"

	"github.com/codepzj/Stellux-Server/internal/label/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/label/internal/repository/dao"
	"github.com/samber/lo"
)

type ILabelRepository interface {
	CreateLabel(ctx context.Context, label *domain.Label) error
	UpdateLabel(ctx context.Context, id uint, label *domain.Label) error
	DeleteLabel(ctx context.Context, id uint) error
	GetLabelById(ctx context.Context, id uint) (*domain.Label, error)
	QueryLabelList(ctx context.Context, labelType string, pageNo int64, pageSize int64) ([]*domain.Label, int64, error)
	GetAllLabelsByType(ctx context.Context, labelType string) ([]*domain.Label, error)
	GetLabelByName(ctx context.Context, name string) (*domain.Label, error)
}

var _ ILabelRepository = (*LabelRepository)(nil)

func NewLabelRepository(dao dao.ILabelDao) *LabelRepository {
	return &LabelRepository{dao: dao}
}

type LabelRepository struct {
	dao dao.ILabelDao
}

func (r *LabelRepository) CreateLabel(ctx context.Context, label *domain.Label) error {
	return r.dao.CreateLabel(ctx, r.LabelDomainToLabelDO(label))
}

func (r *LabelRepository) UpdateLabel(ctx context.Context, id uint, label *domain.Label) error {
	return r.dao.UpdateLabel(ctx, id, r.LabelDomainToLabelDO(label))
}

func (r *LabelRepository) DeleteLabel(ctx context.Context, id uint) error {
	return r.dao.DeleteLabel(ctx, id)
}

func (r *LabelRepository) GetLabelById(ctx context.Context, id uint) (*domain.Label, error) {
	label, err := r.dao.GetLabelById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.LabelDoToDomain(label), nil
}

func (r *LabelRepository) QueryLabelList(ctx context.Context, labelType string, pageNo int64, pageSize int64) ([]*domain.Label, int64, error) {
	count, err := r.dao.GetAllCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	labels, err := r.dao.QueryLabelList(ctx, labelType, pageSize, (pageNo-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	return r.LabelDoToDomainList(labels), count, nil
}

func (r *LabelRepository) GetAllLabelsByType(ctx context.Context, labelType string) ([]*domain.Label, error) {
	labels, err := r.dao.GetAllLabelsByType(ctx, labelType)
	if err != nil {
		return nil, err
	}
	return r.LabelDoToDomainList(labels), nil
}

func (r *LabelRepository) GetLabelByName(ctx context.Context, name string) (*domain.Label, error) {
	label, err := r.dao.GetLabelByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return r.LabelDoToDomain(label), nil
}

func (r *LabelRepository) LabelDomainToLabelDO(label *domain.Label) *dao.Label {
	return &dao.Label{
		LabelType: label.LabelType,
		Name:      label.Name,
	}
}

func (r *LabelRepository) LabelDoToDomain(label *dao.Label) *domain.Label {
	return &domain.Label{
		ID:        label.ID,
		LabelType: label.LabelType,
		Name:      label.Name,
	}
}

func (r *LabelRepository) LabelDoToDomainList(labels []*dao.Label) []*domain.Label {
	return lo.Map(labels, func(label *dao.Label, _ int) *domain.Label {
		return r.LabelDoToDomain(label)
	})
}

func (r *LabelRepository) LabelPostCountDoToDomainList(labelPostCounts []*dao.LabelPostCount) []*domain.LabelPostCount {
	return lo.Map(labelPostCounts, func(labelPostCount *dao.LabelPostCount, _ int) *domain.LabelPostCount {
		return &domain.LabelPostCount{
			ID:        labelPostCount.ID,
			LabelType: labelPostCount.LabelType,
			Name:      labelPostCount.Name,
			Count:     labelPostCount.Count,
		}
	})
}
