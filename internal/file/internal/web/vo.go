package web

import (
	"github.com/codepzj/stellux/server/internal/file/internal/domain"
	"github.com/samber/lo"
)

type FileVO struct {
	Id       string `json:"id"`
	FileName string `json:"file_name"`
	Url      string `json:"url"`
}

func (h *FileHandler) FileDomainToVO(file *domain.File) *FileVO {
	return &FileVO{
		Id:       file.ID.Hex(),
		FileName: file.FileName,
		Url:      file.Url,
	}
}

func (h *FileHandler) FileDomainToVOList(files []*domain.File) []*FileVO {
	return lo.Map(files, func(file *domain.File, _ int) *FileVO {
		return h.FileDomainToVO(file)
	})
}
