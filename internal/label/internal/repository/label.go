package repository

import (
	"context"

	"github.com/codepzj/stellux/server/internal/label/internal/domain"
	"github.com/codepzj/stellux/server/internal/label/internal/repository/dao"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ILabelRepository interface {
	Create(ctx context.Context, label *domain.Label) error
	Update(ctx context.Context, id string, label *domain.Label) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*domain.Label, error)
	GetList(ctx context.Context, labelType string, pageNo int64, pageSize int64) ([]*domain.Label, int64, error)
	GetAllByType(ctx context.Context, labelType string) ([]*domain.Label, error)
}

var _ ILabelRepository = (*LabelRepository)(nil)

func NewLabelRepository(dao dao.ILabelDao) *LabelRepository {
	return &LabelRepository{dao: dao}
}

type LabelRepository struct {
	dao dao.ILabelDao
}

func (r *LabelRepository) Create(ctx context.Context, label *domain.Label) error {
	return r.dao.Create(ctx, r.LabelDomainToLabelDO(label))
}

func (r *LabelRepository) Update(ctx context.Context, id string, label *domain.Label) error {
	bid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.Update(ctx, bid, r.LabelDomainToLabelDO(label))
}

func (r *LabelRepository) Delete(ctx context.Context, id string) error {
	bid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.Delete(ctx, bid)
}

func (r *LabelRepository) Get(ctx context.Context, id string) (*domain.Label, error) {
	bid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	label, err := r.dao.Get(ctx, bid)
	if err != nil {
		return nil, err
	}
	return r.LabelDoToDomain(label), nil
}

func (r *LabelRepository) GetList(ctx context.Context, labelType string, pageNo int64, pageSize int64) ([]*domain.Label, int64, error) {
	labels, count, err := r.dao.GetList(ctx, labelType, pageSize, (pageNo-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	return r.LabelDoToDomainList(labels), count, nil
}

func (r *LabelRepository) GetAllByType(ctx context.Context, labelType string) ([]*domain.Label, error) {
	labels, err := r.dao.GetAllByType(ctx, labelType)
	if err != nil {
		return nil, err
	}
	return r.LabelDoToDomainList(labels), nil
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
