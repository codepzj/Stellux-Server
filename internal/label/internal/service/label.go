package service

import (
	"context"
	"errors"

	"github.com/codepzj/Stellux-Server/internal/label/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/label/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ILabelService interface {
	CreateLabel(ctx context.Context, label *domain.Label) error
	UpdateLabel(ctx context.Context, id string, label *domain.Label) error
	DeleteLabel(ctx context.Context, id string) error
	GetLabelById(ctx context.Context, id string) (*domain.Label, error)
	QueryLabelList(ctx context.Context, labelType string, keyword string, pageNo int64, pageSize int64) ([]*domain.Label, int64, error)
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
	// 检查标签是否存在
	existLabel, err := s.repo.GetLabelByName(ctx, label.Name)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = s.repo.CreateLabel(ctx, label)
			if err != nil {
				logger.Error("创建标签失败",
					logger.WithError(err),
					logger.WithString("name", label.Name),
				)
				return err
			}
			logger.Info("创建标签成功",
				logger.WithString("name", label.Name),
				logger.WithString("type", label.LabelType),
			)
			return nil
		}
		logger.Error("查询标签失败",
			logger.WithError(err),
			logger.WithString("name", label.Name),
		)
		return err
	}

	if existLabel != nil {
		if existLabel.Id.Hex() == label.Id.Hex() {
			logger.Info("标签已存在，跳过创建",
				logger.WithString("name", label.Name),
			)
			return nil
		}
		logger.Warn("标签名称已存在",
			logger.WithString("name", label.Name),
			logger.WithString("existLabelId", existLabel.Id.Hex()),
		)
		return errors.New("标签已存在")
	}

	err = s.repo.CreateLabel(ctx, label)
	if err != nil {
		logger.Error("创建标签失败",
			logger.WithError(err),
			logger.WithString("name", label.Name),
		)
		return err
	}

	logger.Info("创建标签成功",
		logger.WithString("name", label.Name),
		logger.WithString("type", label.LabelType),
	)

	return nil
}

// UpdateLabel 更新标签
func (s *LabelService) UpdateLabel(ctx context.Context, id string, label *domain.Label) error {
	// 检查标签是否存在
	existLabel, err := s.repo.GetLabelByName(ctx, label.Name)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = s.repo.UpdateLabel(ctx, id, label)
			if err != nil {
				logger.Error("更新标签失败",
					logger.WithError(err),
					logger.WithString("labelId", id),
				)
				return err
			}
			logger.Info("更新标签成功",
				logger.WithString("labelId", id),
				logger.WithString("name", label.Name),
			)
			return nil
		}
		logger.Error("查询标签失败",
			logger.WithError(err),
			logger.WithString("name", label.Name),
		)
		return err
	}

	if existLabel != nil {
		if existLabel.Id.Hex() == label.Id.Hex() {
			logger.Info("标签未变更，跳过更新",
				logger.WithString("labelId", id),
			)
			return nil
		}
		logger.Warn("标签名称已被其他标签使用",
			logger.WithString("name", label.Name),
			logger.WithString("existLabelId", existLabel.Id.Hex()),
		)
		return errors.New("标签已存在")
	}

	err = s.repo.UpdateLabel(ctx, id, label)
	if err != nil {
		logger.Error("更新标签失败",
			logger.WithError(err),
			logger.WithString("labelId", id),
		)
		return err
	}

	logger.Info("更新标签成功",
		logger.WithString("labelId", id),
		logger.WithString("name", label.Name),
	)

	return nil
}

// DeleteLabel 删除标签
func (s *LabelService) DeleteLabel(ctx context.Context, id string) error {
	err := s.repo.DeleteLabel(ctx, id)
	if err != nil {
		logger.Error("删除标签失败",
			logger.WithError(err),
			logger.WithString("labelId", id),
		)
		return err
	}

	logger.Info("删除标签成功",
		logger.WithString("labelId", id),
	)

	return nil
}

// GetLabelById 根据id获取标签
func (s *LabelService) GetLabelById(ctx context.Context, id string) (*domain.Label, error) {
	logger.Info("查询标签",
		logger.WithString("method", "GetLabelById"),
		logger.WithString("labelId", id),
	)

	label, err := s.repo.GetLabelById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Warn("标签不存在",
				logger.WithString("labelId", id),
			)
		} else {
			logger.Error("查询标签失败",
				logger.WithError(err),
				logger.WithString("labelId", id),
			)
		}
		return nil, err
	}

	return label, nil
}

// QueryLabelList 分页查询标签
func (s *LabelService) QueryLabelList(ctx context.Context, labelType string, keyword string, pageNo int64, pageSize int64) ([]*domain.Label, int64, error) {
	logger.Info("分页查询标签",
		logger.WithString("method", "QueryLabelList"),
		logger.WithString("labelType", labelType),
		logger.WithString("keyword", keyword),
		logger.WithInt("pageNo", int(pageNo)),
		logger.WithInt("pageSize", int(pageSize)),
	)

	labels, total, err := s.repo.QueryLabelList(ctx, labelType, keyword, pageNo, pageSize)
	if err != nil {
		logger.Error("查询标签列表失败",
			logger.WithError(err),
			logger.WithString("labelType", labelType),
		)
		return nil, 0, err
	}

	return labels, total, nil
}

// GetAllLabelsByType 获取所有标签
func (s *LabelService) GetAllLabelsByType(ctx context.Context, labelType string) ([]*domain.Label, error) {
	logger.Info("查询所有指定类型标签",
		logger.WithString("method", "GetAllLabelsByType"),
		logger.WithString("labelType", labelType),
	)

	labels, err := s.repo.GetAllLabelsByType(ctx, labelType)
	if err != nil {
		logger.Error("查询标签失败",
			logger.WithError(err),
			logger.WithString("labelType", labelType),
		)
		return nil, err
	}

	return labels, nil
}

// GetAllLabelsWithCount 获取所有分类标签及其文章数量
func (s *LabelService) GetAllLabelsWithCount(ctx context.Context) ([]*domain.LabelPostCount, error) {
	logger.Info("查询分类标签及文章数",
		logger.WithString("method", "GetAllLabelsWithCount"),
	)

	labels, err := s.repo.GetCategoryLabelWithCount(ctx)
	if err != nil {
		logger.Error("查询分类标签失败",
			logger.WithError(err),
		)
		return nil, err
	}

	return labels, nil
}

// GetAllTagsLabelWithCount 获取所有标签及其文章数量
func (s *LabelService) GetAllTagsLabelWithCount(ctx context.Context) ([]*domain.LabelPostCount, error) {
	logger.Info("查询标签及文章数",
		logger.WithString("method", "GetAllTagsLabelWithCount"),
	)

	labels, err := s.repo.GetTagsLabelWithCount(ctx)
	if err != nil {
		logger.Error("查询标签失败",
			logger.WithError(err),
		)
		return nil, err
	}

	return labels, nil
}
