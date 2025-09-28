package dao

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"username;type:varchar(20);uniqueIndex;not null"`
	Password string `gorm:"password"`
	Nickname string `gorm:"nickname"`
	RoleId   int    `gorm:"role_id"`
	Avatar   string `gorm:"avatar"`
	Email    string `gorm:"email"`
}

func (User) TableName() string {
	return "user"
}

type UserUpdate struct {
	Nickname string `gorm:"nickname"`
	Avatar   string `gorm:"avatar"`
	Email    string `gorm:"email"`
}

func (UserUpdate) TableName() string {
	return "user"
}

type IUserDao interface {
	Create(ctx context.Context, user *User) error
	Find(ctx context.Context, filter func(db *gorm.DB) *gorm.DB) ([]*User, error)
	GetAllCount(ctx context.Context) (int64, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, id uint, user *User) error
	UpdatePassword(ctx context.Context, id uint, password string) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*User, error)
}

var _ IUserDao = (*UserDao)(nil)

func NewUserDao(db *gorm.DB) *UserDao {
	db.AutoMigrate(&User{})
	return &UserDao{db: db}
}

type UserDao struct {
	db *gorm.DB
}

func (d *UserDao) Create(ctx context.Context, user *User) error {
	err := d.db.Model(&User{}).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *UserDao) Find(ctx context.Context, filter func(db *gorm.DB) *gorm.DB) ([]*User, error) {
	var users []*User
	err := d.db.Model(&User{}).Scopes(filter).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (d *UserDao) GetAllCount(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.Model(&User{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *UserDao) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := d.db.Model(&User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDao) Update(ctx context.Context, id uint, user *User) error {
	res := d.db.Model(&User{}).Where("id = ?", id).Updates(user)
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("更新用户失败")
	}
	return nil
}

func (d *UserDao) UpdatePassword(ctx context.Context, id uint, password string) error {
	res := d.db.Model(&User{}).Where("id = ?", id).Update("password", password)
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("更新用户密码失败")
	}
	return nil
}

func (d *UserDao) Delete(ctx context.Context, id uint) error {
	res := d.db.Model(&User{}).Where("id = ?", id).Delete(&User{})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("删除用户失败")
	}
	return nil
}

func (d *UserDao) GetByID(ctx context.Context, id uint) (*User, error) {
	var user User
	err := d.db.Model(&User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
