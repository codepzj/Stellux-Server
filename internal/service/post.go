package service

import (
	"github.com/codepzj/Stellux-Server/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/domain/model"
	"github.com/codepzj/Stellux-Server/internal/dto/req"
	"github.com/codepzj/Stellux-Server/pkg/apiwrap"
	"github.com/gin-gonic/gin"
)

type PostService struct {
	repo *domain.PostRepo
}

func NewPostService(repo *domain.PostRepo) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (s *PostService) CreatePost(ctx *gin.Context, req req.CreatePostReq) (*apiwrap.Response[any], error) {
	post := &model.Post{
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		Alias_:      req.Alias,
		CategoryID:  req.CategoryID,
		IsPublish:   req.IsPublish,
		IsTop:       req.IsTop,
		Thumbnail:   req.Thumbnail,
	}
	if err := s.repo.CreatePost(ctx, post); err != nil {
		return nil, err
	}
	return apiwrap.Success(), nil
}

func (s *PostService) UpdatePost(ctx *gin.Context, req req.UpdatePostReq) (*apiwrap.Response[any], error) {
	post := &model.Post{
		ID:          req.ID,
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		Alias_:      req.Alias,
		CategoryID:  req.CategoryID,
		IsPublish:   req.IsPublish,
		IsTop:       req.IsTop,
		Thumbnail:   req.Thumbnail,
	}
	if err := s.repo.UpdatePost(ctx, post); err != nil {
		return nil, err
	}
	return apiwrap.Success(), nil
}

func (s *PostService) DeletePost(ctx *gin.Context, id int64) (*apiwrap.Response[any], error) {
	if err := s.repo.DeletePost(ctx, id); err != nil {
		return nil, err
	}
	return apiwrap.Success(), nil
}

func (s *PostService) GetPost(ctx *gin.Context, id int64) (*apiwrap.Response[any], error) {
	post, err := s.repo.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessWithDetail[any](post, "获取文章成功"), nil
}
