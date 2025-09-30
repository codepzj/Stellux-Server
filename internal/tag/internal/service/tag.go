package service

import (
	"github.com/codepzj/Stellux-Server/internal/tag/internal/repository"
)


type ITagService interface {
}

var _ ITagService = (*TagService)(nil)

func NewTagService(repo repository.ITagRepository) *TagService {
	return &TagService{
		repo: repo,
	}
}

type TagService struct {
	repo repository.ITagRepository
}
