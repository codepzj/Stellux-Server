package dao

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name string `gorm:"name"`
	Color string `gorm:"color"`
}

func (Tag) TableName() string {
	return "tag"
}

type ITagDao interface {
	CreateTag(ctx context.Context, tag *Tag) error
	UpdateTag(ctx context.Context, id uint, tag *Tag) error
	DeleteTag(ctx context.Context, id uint) error
	GetTagById(ctx context.Context, id uint) (*Tag, error)
	GetAllCount(ctx context.Context) (int64, error)
	QueryTagList(ctx context.Context, limit int64, skip int64) ([]*Tag, error)
}

var _ ITagDao = (*TagDao)(nil)

func NewTagDao(db *gorm.DB) *TagDao {
	db.AutoMigrate(&Tag{})
	return &TagDao{db: db}
}

type TagDao struct {
	db *gorm.DB
}

// CreateTag 创建Tag
func (d *TagDao) CreateTag(ctx context.Context, tag *Tag) error {
	return d.db.Model(&Tag{}).Create(tag).Error
}

// UpdateTag 更新Tag
func (d *TagDao) UpdateTag(ctx context.Context, id uint, tag *Tag) error {
	res := d.db.Model(&Tag{}).Where("id = ?", id).Updates(tag)
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("更新Tag失败")
	}
	return nil
}

// DeleteTag 删除Tag
func (d *TagDao) DeleteTag(ctx context.Context, id uint) error {
	res := d.db.Model(&Tag{}).Where("id = ?", id).Delete(&Tag{})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("删除Tag失败")
	}
	return nil
}

// GetTagById 根据id获取Tag
func (d *TagDao) GetTagById(ctx context.Context, id uint) (*Tag, error) {
	var tag Tag
	err := d.db.Model(&Tag{}).Where("id = ?", id).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetAllCount 获取所有Tag数量
func (d *TagDao) GetAllCount(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.Model(&Tag{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// QueryTagList 分页查询Tag
func (d *TagDao) QueryTagList(ctx context.Context, limit int64, skip int64) ([]*Tag, error) {
	tagList := make([]*Tag, 0)
	err := d.db.Model(&Tag{}).Limit(int(limit)).Offset(int(skip)).Find(&tagList).Error
	if err != nil {
		return nil, err
	}
	return tagList, nil
}
