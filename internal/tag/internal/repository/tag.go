package repository

import (
	"github.com/codepzj/Stellux-Server/internal/tag/internal/repository/dao"
)

type ITagRepository interface {
}

var _ ITagRepository = (*TagRepository)(nil)

func NewTagRepository(dao dao.ITagDao) *TagRepository {
	return &TagRepository{dao: dao}
}

type TagRepository struct {
	dao dao.ITagDao
}
