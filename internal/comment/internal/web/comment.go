package web

import (
	"github.com/codepzj/Stellux-Server/internal/comment/internal/service"
	"github.com/gin-gonic/gin"
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

}
