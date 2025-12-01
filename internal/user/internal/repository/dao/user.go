package dao

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty"`
	Username  string        `bson:"username"`
	Password  string        `bson:"password"`
	Nickname  string        `bson:"nickname"`
	RoleId    int           `bson:"role_id"`
	Avatar    string        `bson:"avatar"`
	Email     string        `bson:"email"`
}

type UserUpdate struct {
	Nickname string `bson:"nickname,omitempty"`
	Avatar   string `bson:"avatar,omitempty"`
	Email    string `bson:"email,omitempty"`
}

type IUserDao interface {
	Create(ctx context.Context, user *User) (bson.ObjectID, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, id bson.ObjectID, user *User) error
	UpdatePassword(ctx context.Context, id bson.ObjectID, password string) error
	Delete(ctx context.Context, id bson.ObjectID) error
	FindByCondition(ctx context.Context, skip, limit int64, sort bson.M) ([]*User, int64, error)
	GetByID(ctx context.Context, id bson.ObjectID) (*User, error)
}

var _ IUserDao = (*UserDao)(nil)

func NewUserDao(db *mongo.Database) *UserDao {
	return &UserDao{coll: db.Collection("user")}
}

type UserDao struct {
	coll *mongo.Collection
}

func (d *UserDao) Create(ctx context.Context, user *User) (bson.ObjectID, error) {
	user.ID = bson.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	res, err := d.coll.InsertOne(ctx, user)
	if err != nil {
		return bson.ObjectID{}, err
	}
	if res.InsertedID == nil {
		return bson.ObjectID{}, errors.Wrap(err, "新增用户失败")
	}
	return res.InsertedID.(bson.ObjectID), nil
}

func (d *UserDao) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := d.coll.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDao) Update(ctx context.Context, id bson.ObjectID, user *User) error {
	update := bson.M{
		"$set": bson.M{
			"nickname":   user.Nickname,
			"avatar":     user.Avatar,
			"email":      user.Email,
			"updated_at": time.Now(),
		},
	}
	res, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("更新用户失败")
	}
	return nil
}

func (d *UserDao) UpdatePassword(ctx context.Context, id bson.ObjectID, password string) error {
	update := bson.M{
		"$set": bson.M{
			"password":   password,
			"updated_at": time.Now(),
		},
	}
	res, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("更新密码失败")
	}
	return nil
}

func (d *UserDao) Delete(ctx context.Context, id bson.ObjectID) error {
	res, err := d.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("删除用户失败")
	}
	return nil
}

func (d *UserDao) FindByCondition(ctx context.Context, skip, limit int64, sort bson.M) ([]*User, int64, error) {
	count, err := d.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(sort)
	cursor, err := d.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (d *UserDao) GetByID(ctx context.Context, id bson.ObjectID) (*User, error) {
	var user User
	err := d.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDao) UserToUpdate(user *User) *UserUpdate {
	return &UserUpdate{
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Email:    user.Email,
	}
}
