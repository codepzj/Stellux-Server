package dao

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Label struct {
	gorm.Model
	LabelType string `gorm:"label_type"`
	Name      string `gorm:"name"`
}

func (Label) TableName() string {
	return "label"
}

// LabelPostCount 每种标签文章的数量
type LabelPostCount struct {
	gorm.Model
	LabelType string `gorm:"label_type"`
	Name      string `gorm:"name"`
	Count     int    `gorm:"post_count"`
}

type ILabelDao interface {
	CreateLabel(ctx context.Context, label *Label) error
	UpdateLabel(ctx context.Context, id uint, label *Label) error
	DeleteLabel(ctx context.Context, id uint) error
	GetLabelById(ctx context.Context, id uint) (*Label, error)
	GetAllCount(ctx context.Context) (int64, error)
	QueryLabelList(ctx context.Context, labelType string, limit int64, skip int64) ([]*Label, error)
	GetAllLabelsByType(ctx context.Context, labelType string) ([]*Label, error)
	GetLabelByName(ctx context.Context, name string) (*Label, error)
}

var _ ILabelDao = (*LabelDao)(nil)

func NewLabelDao(db *gorm.DB) *LabelDao {
	db.AutoMigrate(&Label{})
	return &LabelDao{db: db}
}

type LabelDao struct {
	db *gorm.DB
}

// CreateLabel 创建标签
func (d *LabelDao) CreateLabel(ctx context.Context, label *Label) error {
	return d.db.Model(&Label{}).Create(label).Error
}

// UpdateLabel 更新标签
func (d *LabelDao) UpdateLabel(ctx context.Context, id uint, label *Label) error {
	res := d.db.Model(&Label{}).Where("id = ?", id).Updates(label)
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("更新标签失败")
	}
	return nil
}

// DeleteLabel 删除标签
func (d *LabelDao) DeleteLabel(ctx context.Context, id uint) error {
	res := d.db.Model(&Label{}).Where("id = ?", id).Delete(&Label{})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("删除标签失败")
	}
	return nil
}

// GetLabelById 根据id获取标签
func (d *LabelDao) GetLabelById(ctx context.Context, id uint) (*Label, error) {
	var label Label
	err := d.db.Model(&Label{}).Where("id = ?", id).First(&label).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

// GetAllCount 获取所有标签数量
func (d *LabelDao) GetAllCount(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.Model(&Label{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// QueryLabelList 分页查询标签
func (d *LabelDao) QueryLabelList(ctx context.Context, labelType string, limit int64, skip int64) ([]*Label, error) {
	labelList := make([]*Label, 0)
	err := d.db.Model(&Label{}).Where("label_type = ?", labelType).Limit(int(limit)).Offset(int(skip)).Find(&labelList).Error
	if err != nil {
		return nil, err
	}
	return labelList, nil
}

// GetAllLabelsByType 通过类型获取所有标签
func (d *LabelDao) GetAllLabelsByType(ctx context.Context, labelType string) ([]*Label, error) {
	labelList := make([]*Label, 0)
	err := d.db.Model(&Label{}).Where("label_type = ?", labelType).Find(&labelList).Error
	if err != nil {
		return nil, err
	}
	return labelList, nil
}

// GetLabelByName 根据名称获取标签
func (d *LabelDao) GetLabelByName(ctx context.Context, name string) (*Label, error) {
	var label Label
	err := d.db.Model(&Label{}).Where("name = ?", name).First(&label).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}
