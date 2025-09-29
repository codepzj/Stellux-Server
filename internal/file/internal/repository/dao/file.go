package dao

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	FileName string
	Url      string
	Dst      string
}

type IFileDao interface {
	Create(ctx context.Context, file *File) error
	Get(ctx context.Context, id uint) (*File, error)
	GetList(ctx context.Context, skip int64, limit int64) ([]*File, error)
	GetAllCount(ctx context.Context) (int64, error)
	GetListByIDList(ctx context.Context, idList []uint) ([]*File, error)
	Delete(ctx context.Context, id uint) error
	DeleteMany(ctx context.Context, idList []uint) error
}

var _ IFileDao = (*FileDao)(nil)

func NewFileDao(db *gorm.DB) *FileDao {
	db.AutoMigrate(&File{})
	return &FileDao{db: db}
}

type FileDao struct {
	db *gorm.DB
}

func (d *FileDao) Create(ctx context.Context, file *File) error {
	return d.db.Model(&File{}).Create(file).Error
}

func (d *FileDao) Get(ctx context.Context, id uint) (*File, error) {
	var file File
	err := d.db.Model(&File{}).Where("id = ?", id).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (d *FileDao) GetList(ctx context.Context, skip int64, limit int64) ([]*File, error) {
	files := make([]*File, 0)
	err := d.db.Model(&File{}).Offset(int(skip)).Limit(int(limit)).Order("created_at desc").Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (d *FileDao) GetAllCount(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.Model(&File{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *FileDao) GetListByIDList(ctx context.Context, idList []uint) ([]*File, error) {
	files := make([]*File, 0)
	err := d.db.Model(&File{}).Where("id IN ?", idList).Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (d *FileDao) Delete(ctx context.Context, id uint) error {
	res := d.db.Model(&File{}).Where("id = ?", id).Delete(&File{})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("删除文件失败")
	}
	return nil
}

func (d *FileDao) DeleteMany(ctx context.Context, idList []uint) error {
	res := d.db.Model(&File{}).Where("id IN ?", idList).Delete(&File{})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected != int64(len(idList)) {
		return errors.New("删除文件失败")
	}
	return nil
}
