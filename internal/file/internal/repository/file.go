package repository

import (
	"context"

	"github.com/codepzj/stellux/server/internal/file/internal/domain"
	"github.com/codepzj/stellux/server/internal/file/internal/repository/dao"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IFileRepository interface {
	CreateMany(ctx context.Context, files []*domain.File) error
	Get(ctx context.Context, id bson.ObjectID) (*domain.File, error)
	GetList(ctx context.Context, page *apiwrap.Page) ([]*domain.File, int64, error)
	GetListByIDList(ctx context.Context, idList []bson.ObjectID) ([]*domain.File, error)
	Delete(ctx context.Context, id bson.ObjectID) error
	DeleteMany(ctx context.Context, idList []bson.ObjectID) error
}

var _ IFileRepository = (*FileRepository)(nil)

func NewFileRepository(dao dao.IFileDao) *FileRepository {
	return &FileRepository{dao: dao}
}

type FileRepository struct {
	dao dao.IFileDao
}

func (r *FileRepository) CreateMany(ctx context.Context, files []*domain.File) error {
	return r.dao.CreateMany(ctx, r.FileDomainToDaoList(files))
}

func (r *FileRepository) Get(ctx context.Context, id bson.ObjectID) (*domain.File, error) {
	file, err := r.dao.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.FileDaoToDomain(file), nil
}

func (r *FileRepository) GetList(ctx context.Context, page *apiwrap.Page) ([]*domain.File, int64, error) {
	limit := page.PageSize
	skip := (page.PageNo - 1) * page.PageSize
	files, count, err := r.dao.GetList(ctx, skip, limit)
	if err != nil {
		return nil, 0, err
	}
	return r.FileDaoToDomainList(files), count, nil
}

func (r *FileRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	return r.dao.Delete(ctx, id)
}

func (r *FileRepository) DeleteMany(ctx context.Context, idList []bson.ObjectID) error {
	return r.dao.DeleteMany(ctx, idList)
}

func (r *FileRepository) GetListByIDList(ctx context.Context, idList []bson.ObjectID) ([]*domain.File, error) {
	files, err := r.dao.GetListByIDList(ctx, idList)
	if err != nil {
		return nil, err
	}
	return r.FileDaoToDomainList(files), nil
}

func (r *FileRepository) FileDomainToDao(file *domain.File) *dao.File {
	return &dao.File{
		FileName: file.FileName,
		Url:      file.Url,
		Dst:      file.Dst,
	}
}

func (r *FileRepository) FileDomainToDaoList(files []*domain.File) []*dao.File {
	return lo.Map(files, func(file *domain.File, _ int) *dao.File {
		return r.FileDomainToDao(file)
	})
}

func (r *FileRepository) FileDaoToDomain(file *dao.File) *domain.File {
	return &domain.File{
		ID:       file.ID,
		FileName: file.FileName,
		Url:      file.Url,
		Dst:      file.Dst,
	}
}

func (r *FileRepository) FileDaoToDomainList(files []*dao.File) []*domain.File {
	return lo.Map(files, func(file *dao.File, _ int) *domain.File {
		return r.FileDaoToDomain(file)
	})
}
