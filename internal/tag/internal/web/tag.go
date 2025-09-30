package web

import (
	"github.com/codepzj/Stellux-Server/internal/tag/internal/service"
	"github.com/gin-gonic/gin"
)

func NewTagHandler(serv service.ITagService) *TagHandler {
	return &TagHandler{
		serv: serv,
	}
}

type TagHandler struct {
	serv service.ITagService
}

func (h *TagHandler) RegisterGinRoutes(engine *gin.Engine) {

}
