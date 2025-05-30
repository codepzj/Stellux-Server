package dao

import (
	"context"
	"errors"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type File struct {
	mongox.Model `bson:",inline"`
	FileName     string
	Url          string
	Dst          string
}

type IFileDao interface {
	CreateMany(ctx context.Context, files []*File) error
	Get(ctx context.Context, id bson.ObjectID) (*File, error)
	GetList(ctx context.Context, skip int64, limit int64) ([]*File, int64, error)
	GetListByIDList(ctx context.Context, idList []bson.ObjectID) ([]*File, error)

	Delete(ctx context.Context, id bson.ObjectID) error
	DeleteMany(ctx context.Context, idList []bson.ObjectID) error
}

var _ IFileDao = (*FileDao)(nil)

func NewFileDao(db *mongox.Database) *FileDao {
	return &FileDao{coll: mongox.NewCollection[File](db, "file")}
}

type FileDao struct {
	coll *mongox.Collection[File]
}

func (d *FileDao) CreateMany(ctx context.Context, files []*File) error {
	insertResult, err := d.coll.Creator().InsertMany(ctx, files)
	if err != nil {
		return err
	}
	if len(insertResult.InsertedIDs) != len(files) {
		return errors.New("文件信息与保存数量不匹配")
	}
	return nil
}

func (d *FileDao) Get(ctx context.Context, id bson.ObjectID) (*File, error) {
	return d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
}

func (d *FileDao) GetList(ctx context.Context, skip int64, limit int64) ([]*File, int64, error) {
	count, err := d.coll.Finder().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	files, err := d.coll.Finder().Skip(skip).Limit(limit).Sort(bson.M{"created_at": -1}).Find(ctx)
	if err != nil {
		return nil, 0, err
	}
	return files, count, nil
}

func (d *FileDao) GetListByIDList(ctx context.Context, idList []bson.ObjectID) ([]*File, error) {
	return d.coll.Finder().Filter(query.In("_id", idList...)).Find(ctx)
}

func (d *FileDao) Delete(ctx context.Context, id bson.ObjectID) error {
	deleteResult, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return errors.New("删除文件失败")
	}
	return nil
}

func (d *FileDao) DeleteMany(ctx context.Context, idList []bson.ObjectID) error {
	deleteResult, err := d.coll.Deleter().Filter(query.In("_id", idList...)).DeleteMany(ctx)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount != int64(len(idList)) {
		return errors.New("删除文件失败")
	}
	return nil
}
