package web

import (
	"net/http"

	"github.com/codepzj/stellux/server/internal/comment/internal/domain"
	"github.com/codepzj/stellux/server/internal/comment/internal/service"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewCommentHandler(serv service.ICommentService) *CommentHandler {
	return &CommentHandler{
		serv: serv,
	}
}

type CommentHandler struct {
	serv service.ICommentService
}

func (h *CommentHandler) RegisterGinRoutes(engine *gin.Engine) {
	commentGroup := engine.Group("/comment")
	{
		commentGroup.POST("/create", apiwrap.WrapWithJson(h.Create))
		commentGroup.GET("/list/:post_id", apiwrap.Wrap(h.GetListByPostId))
	}
	adminCommentGroup := engine.Group("/admin-api/comment")
	{
		adminCommentGroup.POST("/edit", apiwrap.WrapWithJson(h.AdminEdit))
		adminCommentGroup.POST("/delete", apiwrap.Wrap(h.AdminDelete))
	}
}

// Create 创建评论
func (h *CommentHandler) Create(c *gin.Context, req CommentRequest) *apiwrap.Response[any] {
	postObjId, _ := bson.ObjectIDFromHex(req.PostId)
	commentObjId, _ := bson.ObjectIDFromHex(req.CommentId)
	comment := &domain.Comment{
		Content:   req.Content,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Url:       req.Url,
		Avatar:    req.Avatar,
		PostId:    postObjId,
		CommentId: commentObjId,
	}

	err := h.serv.Create(c, comment)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, "创建评论失败")
	}
	return apiwrap.SuccessWithMsg("创建评论成功")
}

// GetListByPostId 根据帖子id获取评论列表
func (h *CommentHandler) GetListByPostId(c *gin.Context) *apiwrap.Response[any] {
	postId := c.Param("post_id")
	postObjId, _ := bson.ObjectIDFromHex(postId)
	comments, err := h.serv.GetListByPostId(c, postObjId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, "获取评论列表失败")
	}
	vo := make([]CommentVO, len(comments))
	for i, comment := range comments {
		vo[i] = CommentVO{
			Id:        comment.Id.Hex(),
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Nickname:  comment.Nickname,
			Email:     comment.Email,
			Url:       comment.Url,
			Avatar:    comment.Avatar,
			Content:   comment.Content,
			PostId:    comment.PostId.Hex(),
			CommentId: comment.CommentId.Hex(),
		}
	}
	return apiwrap.SuccessWithDetail[any](vo, "获取评论列表成功")
}

// AdminEdit 管理员编辑评论
func (h *CommentHandler) AdminEdit(c *gin.Context, req CommentEditRequest) *apiwrap.Response[any] {
	objId, _ := bson.ObjectIDFromHex(req.Id)
	comment := &domain.Comment{
		Id:       objId,
		Content:  req.Content,
		Nickname: req.Nickname,
		Email:    req.Email,
		Url:      req.Url,
		Avatar:   req.Avatar,
	}

	if err := h.serv.AdminEdit(c, comment); err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, "编辑评论失败")
	}
	return apiwrap.SuccessWithMsg("编辑评论成功")
}

func (h *CommentHandler) AdminDelete(c *gin.Context) *apiwrap.Response[any] {
	id := c.Param("id")
	objId, _ := bson.ObjectIDFromHex(id)
	if err := h.serv.AdminDelete(c, objId); err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, "删除评论失败")
	}
	return apiwrap.SuccessWithMsg("删除评论成功")
}
