package repository

import (
	"context"

	"github.com/codepzj/stellux/server/internal/label/internal/domain"
	"github.com/codepzj/stellux/server/internal/label/internal/repository/dao"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ILabelRepository interface {
	CreateLabel(ctx context.Context, label *domain.Label) error
	UpdateLabel(ctx context.Context, id string, label *domain.Label) error
	DeleteLabel(ctx context.Context, id string) error
	GetLabelById(ctx context.Context, id string) (*domain.Label, error)
	QueryLabelList(ctx context.Context, labelType string, pageNo int64, pageSize int64) ([]*domain.Label, int64, error)
	GetAllLabelsByType(ctx context.Context, labelType string) ([]*domain.Label, error)
	GetCategoryLabelWithCount(ctx context.Context) ([]*domain.LabelPostCount, error)
	GetTagsLabelWithCount(ctx context.Context) ([]*domain.LabelPostCount, error)
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

func (r *LabelRepository) UpdateLabel(ctx context.Context, id string, label *domain.Label) error {
	bid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.UpdateLabel(ctx, bid, r.LabelDomainToLabelDO(label))
}

func (r *LabelRepository) DeleteLabel(ctx context.Context, id string) error {
	bid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.DeleteLabel(ctx, bid)
}

func (r *LabelRepository) GetLabelById(ctx context.Context, id string) (*domain.Label, error) {
	bid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	label, err := r.dao.GetLabelById(ctx, bid)
	if err != nil {
		return nil, err
	}
	return r.LabelDoToDomain(label), nil
}

func (r *LabelRepository) QueryLabelList(ctx context.Context, labelType string, pageNo int64, pageSize int64) ([]*domain.Label, int64, error) {
	labels, count, err := r.dao.QueryLabelList(ctx, labelType, pageSize, (pageNo-1)*pageSize)
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

func (r *LabelRepository) GetCategoryLabelWithCount(ctx context.Context) ([]*domain.LabelPostCount, error) {
	labelWithCount, err := r.dao.GetCategoryLabelWithCount(ctx)
	if err != nil {
		return nil, err
	}

	return r.LabelPostCountDoToDomainList(labelWithCount), nil
}

func (r *LabelRepository) GetTagsLabelWithCount(ctx context.Context) ([]*domain.LabelPostCount, error) {
	labelWithCount, err := r.dao.GetTagsLabelWithCount(ctx)
	if err != nil {
		return nil, err
	}
	return r.LabelPostCountDoToDomainList(labelWithCount), nil
}

func (r *LabelRepository) LabelDomainToLabelDO(label *domain.Label) *dao.Label {
	return &dao.Label{
		LabelType: label.LabelType,
		Name:      label.Name,
	}
}

func (r *LabelRepository) LabelDoToDomain(label *dao.Label) *domain.Label {
	return &domain.Label{
		Id:        label.ID,
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
			Id:        labelPostCount.ID,
			LabelType: labelPostCount.LabelType,
			Name:      labelPostCount.Name,
			Count:     labelPostCount.Count,
		}
	})
}
