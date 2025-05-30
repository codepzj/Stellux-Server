package dao

import (
	"context"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	mongox.Model `bson:",inline"`
	Username     string `bson:"username"`
	Password     string `bson:"password"`
	Nickname     string `bson:"nickname"`
	RoleId       int    `bson:"role_id"`
	Avatar       string `bson:"avatar"`
	Email        string `bson:"email"`
}

type UserUpdate struct {
	Nickname string `bson:"nickname,omitempty"`
	Avatar   string `bson:"avatar,omitempty"`
	Email    string `bson:"email,omitempty"`
}

type IUserDao interface {
	Create(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, id bson.ObjectID, user *User) error
	UpdatePassword(ctx context.Context, id bson.ObjectID, password string) error
	Delete(ctx context.Context, id bson.ObjectID) error
	FindByCondition(ctx context.Context, findOptions *options.FindOptionsBuilder) ([]*User, int64, error)
	GetByID(ctx context.Context, id bson.ObjectID) (*User, error)
}

var _ IUserDao = (*UserDao)(nil)

func NewUserDao(db *mongox.Database) *UserDao {
	return &UserDao{coll: mongox.NewCollection[User](db, "user")}
}

type UserDao struct {
	coll *mongox.Collection[User]
}

func (d *UserDao) Create(ctx context.Context, user *User) error {
	res, err := d.coll.Creator().InsertOne(ctx, user)
	if err != nil {
		return err
	}
	if res.InsertedID == nil {
		return errors.Wrap(err, "新增用户失败")
	}
	return nil
}

func (d *UserDao) GetByUsername(ctx context.Context, username string) (*User, error) {
	return d.coll.Finder().Filter(bson.M{"username": username}).FindOne(ctx)
}

func (d *UserDao) Update(ctx context.Context, id bson.ObjectID, user *User) error {
	res, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.SetFields(d.UserToUpdate(user))).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("更新用户失败")
	}
	return nil
}

func (d *UserDao) UpdatePassword(ctx context.Context, id bson.ObjectID, password string) error {
	res, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.Set("password", password)).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("更新密码失败")
	}
	return nil
}

func (d *UserDao) Delete(ctx context.Context, id bson.ObjectID) error {
	res, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("删除用户失败")
	}
	return nil
}

func (d *UserDao) FindByCondition(ctx context.Context, findOptions *options.FindOptionsBuilder) ([]*User, int64, error) {
	count, err := d.coll.Finder().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	users, err := d.coll.Finder().Find(ctx, findOptions)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (d *UserDao) GetByID(ctx context.Context, id bson.ObjectID) (*User, error) {
	return d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
}

func (d *UserDao) UserToUpdate(user *User) *UserUpdate {
	return &UserUpdate{
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Email:    user.Email,
	}
}
