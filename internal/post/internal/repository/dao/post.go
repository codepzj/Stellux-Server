package dao

import (
	"context"
	"fmt"

	"github.com/codepzj/Stellux-Server/internal/label"
	"github.com/codepzj/Stellux-Server/internal/pkg/utils"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title       string         `gorm:"column:title"`
	Content     string         `gorm:"column:content"`
	Description string         `gorm:"column:description"`
	Author      string         `gorm:"column:author"`
	Alias       string         `gorm:"column:alias"`
	CategoryID  uint           `gorm:"column:category_id"`
	Category    label.Domain   `gorm:"-"`
	TagsID      pq.Int64Array  `gorm:"column:tags_id;type:integer[]"`
	Tags        []label.Domain `gorm:"-"`
	IsPublish   bool           `gorm:"column:is_publish"`
	IsTop       bool           `gorm:"column:is_top"`
	Thumbnail   string         `gorm:"column:thumbnail"`
}

func (Post) TableName() string {
	return "post"
}

type PostUpdate struct {
	Title       string
	Content     string
	Description string
	Author      string
	Alias       string
	CategoryID  uint
	TagsID      pq.Int64Array `gorm:"type:integer[]"`
	IsPublish   bool
	IsTop       bool
	Thumbnail   string
}

func (PostUpdate) TableName() string {
	return "post"
}

type IPostDao interface {
	Create(ctx context.Context, post *Post) error
	Update(ctx context.Context, id uint, post *PostUpdate) error
	UpdatePostPublishStatus(ctx context.Context, id uint, isPublish bool) error
	SoftDelete(ctx context.Context, id uint) error
	SoftDeleteBatch(ctx context.Context, ids []uint) error
	Delete(ctx context.Context, id uint) error
	DeleteBatch(ctx context.Context, ids []uint) error
	Restore(ctx context.Context, id uint) error
	RestoreBatch(ctx context.Context, ids []uint) error
	GetByID(ctx context.Context, id uint) (*Post, error)
	GetListByKeyWord(ctx context.Context, keyWord string) ([]*Post, error)
	GetWithCategoryAndTags(ctx context.Context, id uint) (*Post, error)
	GetListByFilter(ctx context.Context, filter func(tx *gorm.DB) *gorm.DB) ([]*Post, error)
	GetAllCount(ctx context.Context) (int64, error)
	GetAllPublishPost(ctx context.Context) ([]*Post, error)
	FindByAlias(ctx context.Context, alias string) (*Post, error)
}

var _ IPostDao = (*PostDao)(nil)

func NewPostDao(db *gorm.DB) *PostDao {
	fmt.Println("NewPostDao")
	err := db.AutoMigrate(&Post{})
	if err != nil {
		fmt.Println("AutoMigrate Post error", err)
	}
	return &PostDao{db: db}
}

type PostDao struct {
	db *gorm.DB
}

// Create 创建文章
func (d *PostDao) Create(ctx context.Context, post *Post) error {
	return d.db.Model(&Post{}).Create(post).Error
}

// Update 更新文章
func (d *PostDao) Update(ctx context.Context, id uint, post *PostUpdate) error {
	res := d.db.Model(&Post{}).Where("id = ?", id).Updates(post)
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("文章更新失败")
	}
	return nil
}

// UpdatePostPublishStatus 更新文章发布状态
func (d *PostDao) UpdatePostPublishStatus(ctx context.Context, id uint, isPublish bool) error {
	res := d.db.Model(&Post{}).Where("id = ?", id).Update("is_publish", isPublish)
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("文章发布状态更新失败")
	}
	return nil
}

// SoftDelete 软删除文章
func (d *PostDao) SoftDelete(ctx context.Context, id uint) error {
	res := d.db.Delete(&Post{}, id)
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("文章软删除失败")
	}
	return nil
}

// Delete 删除文章
func (d *PostDao) Delete(ctx context.Context, id uint) error {
	res := d.db.Unscoped().Delete(&Post{}, id)
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("文章删除失败")
	}
	return nil
}

// DeleteBatch 批量删除文章
func (d *PostDao) DeleteBatch(ctx context.Context, ids []uint) error {
	res := d.db.Unscoped().Delete(&Post{}, ids)
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected != int64(len(ids)) {
		return errors.New("文章删除失败")
	}
	return nil
}

// Restore 恢复文章
func (d *PostDao) Restore(ctx context.Context, id uint) error {
	res := d.db.Model(&Post{}).Where("id = ?", id).Update("deleted_at", nil)
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("文章恢复失败")
	}
	return nil
}

// GetWithCategoryAndTags 获取有分类和标签的文章
func (d *PostDao) GetWithCategoryAndTags(ctx context.Context, id uint) (*Post, error) {
	var post Post
	err := d.db.Model(&Post{}).Where("id = ?", id).Preload("Category").Preload("Tags").First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetByID 获取文章
func (d *PostDao) GetByID(ctx context.Context, id uint) (*Post, error) {
	var post Post
	err := d.db.Model(&Post{}).Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetListByKeyWord 根据关键词获取文章列表
func (d *PostDao) GetListByKeyWord(ctx context.Context, keyWord string) ([]*Post, error) {
	var posts []Post
	err := d.db.Model(&Post{}).Where("title LIKE ? OR description LIKE ?", "%"+keyWord+"%", "%"+keyWord+"%").Where("is_publish = ?", true).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return utils.ValToPtrList(posts), nil
}

// GetListByFilter 根据条件获取文章列表
func (d *PostDao) GetListByFilter(ctx context.Context, filter func(tx *gorm.DB) *gorm.DB) ([]*Post, error) {
	var posts []Post
	err := d.db.Model(&Post{}).Scopes(filter).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return utils.ValToPtrList(posts), nil
}

func (d *PostDao) GetAllCount(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.Model(&Post{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SoftDeleteBatch 批量软删除文章
func (d *PostDao) SoftDeleteBatch(ctx context.Context, ids []uint) error {
	res := d.db.Delete(&Post{}, ids)
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected != int64(len(ids)) {
		return errors.New("文章软删除失败")
	}
	return nil
}

// RestoreBatch 批量恢复文章
func (d *PostDao) RestoreBatch(ctx context.Context, ids []uint) error {
	res := d.db.Model(&Post{}).Where("id IN (?)", ids).Update("deleted_at", nil)
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected != int64(len(ids)) {
		return errors.New("文章恢复失败")
	}
	return nil
}

// GetAllPublishPost 获取所有发布文章
func (d *PostDao) GetAllPublishPost(ctx context.Context) ([]*Post, error) {
	var posts []Post
	err := d.db.Model(&Post{}).Where("is_publish = ?", true).Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return utils.ValToPtrList(posts), nil
}

// FindByAlias 根据别名获取文章
func (d *PostDao) FindByAlias(ctx context.Context, alias string) (*Post, error) {
	var post Post
	err := d.db.Model(&Post{}).Where("alias = ?", alias).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}
