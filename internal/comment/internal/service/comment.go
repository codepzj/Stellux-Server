package service

import (
	"github.com/codepzj/Stellux-Server/internal/comment/internal/repository"
)

type ICommentService interface {
}

var _ ICommentService = (*CommentService)(nil)

func NewCommentService(repo repository.ICommentRepository) *CommentService {
	return &CommentService{
		repo: repo,
	}
}

type CommentService struct {
	repo repository.ICommentRepository
}
