package web

import (
	"net/http"

	"github.com/codepzj/stellux/server/internal/document_content/internal/domain"
	"github.com/codepzj/stellux/server/internal/document_content/internal/service"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/gin-gonic/gin"
)

func NewDocumentContentHandler(serv service.IDocumentContentService) *DocumentContentHandler {
	return &DocumentContentHandler{
		serv: serv,
	}
}

type DocumentContentHandler struct {
	serv service.IDocumentContentService
}

func (h *DocumentContentHandler) RegisterGinRoutes(engine *gin.Engine) {
	documentContentGroup := engine.Group("/documentContent")
	documentContentGroup.POST("/create", apiwrap.WrapWithJson(h.CreateDocumentContent))
	documentContentGroup.GET("/findByDocumentID", h.FindDocumentContentByDocumentID)
}

func (h *DocumentContentHandler) CreateDocumentContent(c *gin.Context, dto *CreateDocumentContentDto) *apiwrap.Response[any] {
	err := h.serv.CreateDocumentContent(c, &domain.DocumentContent{
		DocumentId:   dto.DocumentId.ToObjectID(),
		Title:        dto.Title,
		Content:      dto.Content,
		Version:      dto.Version,
		Alias:        dto.Alias,
		ParentID:     dto.ParentID.ToObjectID(),
		IsDir:        dto.IsDir,
		LikeCount:    dto.LikeCount,
		DislikeCount: dto.DislikeCount,
		CommentCount: dto.CommentCount,
	})
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("创建文档内容成功")
}

func (h *DocumentContentHandler) FindDocumentContentByDocumentID(c *gin.Context) {
	documentID := apiwrap.ConvertBsonID(c.Query("documentId"))

	docs, err := h.serv.FindDocumentContentByDocumentID(c, documentID.ToObjectID())
	if err != nil {
		apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, docs)
}
