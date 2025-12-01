package dao

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type File struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty"`
	FileName  string        `bson:"file_name"`
	Url       string        `bson:"url"`
	Dst       string        `bson:"dst"`
}

type IFileDao interface {
	Create(ctx context.Context, file *File) error
	Get(ctx context.Context, id bson.ObjectID) (*File, error)
	GetList(ctx context.Context, skip int64, limit int64) ([]*File, int64, error)
	GetListByIDList(ctx context.Context, idList []bson.ObjectID) ([]*File, error)

	Delete(ctx context.Context, id bson.ObjectID) error
	DeleteMany(ctx context.Context, idList []bson.ObjectID) error
}

var _ IFileDao = (*FileDao)(nil)

func NewFileDao(db *mongo.Database) *FileDao {
	return &FileDao{coll: db.Collection("file")}
}

type FileDao struct {
	coll *mongo.Collection
}

func (d *FileDao) Create(ctx context.Context, file *File) error {
	file.ID = bson.NewObjectID()
	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()
	result, err := d.coll.InsertOne(ctx, file)
	if err != nil {
		return err
	}
	if result.InsertedID == nil {
		return errors.New("保存文件失败")
	}
	return nil
}

func (d *FileDao) Get(ctx context.Context, id bson.ObjectID) (*File, error) {
	var file File
	err := d.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (d *FileDao) GetList(ctx context.Context, skip int64, limit int64) ([]*File, int64, error) {
	count, err := d.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"created_at": -1})
	cursor, err := d.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var files []*File
	if err = cursor.All(ctx, &files); err != nil {
		return nil, 0, err
	}
	return files, count, nil
}

func (d *FileDao) GetListByIDList(ctx context.Context, idList []bson.ObjectID) ([]*File, error) {
	cursor, err := d.coll.Find(ctx, bson.M{"_id": bson.M{"$in": idList}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []*File
	if err = cursor.All(ctx, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (d *FileDao) Delete(ctx context.Context, id bson.ObjectID) error {
	deleteResult, err := d.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return errors.New("删除文件失败")
	}
	return nil
}

func (d *FileDao) DeleteMany(ctx context.Context, idList []bson.ObjectID) error {
	deleteResult, err := d.coll.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": idList}})
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount != int64(len(idList)) {
		return errors.New("删除文件失败")
	}
	return nil
}
