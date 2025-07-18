package web

import (
	"fmt"
	"net/http"

	"github.com/codepzj/stellux/server/internal/document_content/internal/domain"
	"github.com/codepzj/stellux/server/internal/document_content/internal/service"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/codepzj/stellux/server/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	adminDocumentContentGroup := engine.Group("/admin/documentContent")
	{
		adminDocumentContentGroup.Use(middleware.JWT())
		adminDocumentContentGroup.POST("/create", apiwrap.WrapWithJson(h.CreateDocumentContent))                 // 管理员创建文档内容
		adminDocumentContentGroup.GET("/:id", apiwrap.Wrap(h.FindDocumentContentById))                           // 管理员查询特定Id的文档内容
		adminDocumentContentGroup.DELETE("/:id", apiwrap.Wrap(h.DeleteDocumentContentById))                      // 管理员删除特定Id的文档内容
		adminDocumentContentGroup.GET("/getAllDocByParentId", apiwrap.Wrap(h.FindDocumentContentByParentId))     // 管理员查询特定父级Id的所有子文档内容
		adminDocumentContentGroup.GET("/getAllDocByDocumentId", apiwrap.Wrap(h.FindDocumentContentByDocumentId)) // 管理员查询特定文档Id的所有子文档内容
		adminDocumentContentGroup.PUT("/update", apiwrap.WrapWithJson(h.UpdateDocumentContentById))                 // 管理员更新特定Id的文档内容
	}

}

// CreateDocumentContent 管理员创建文档内容
func (h *DocumentContentHandler) CreateDocumentContent(c *gin.Context, dto CreateDocumentContentDto) *apiwrap.Response[any] {
	documentId, _ := bson.ObjectIDFromHex(dto.DocumentId)
	parentId, _ := bson.ObjectIDFromHex(dto.ParentId)

	id, err := h.serv.CreateDocumentContent(c, domain.DocumentContent{
		DocumentId:  documentId,
		Title:       dto.Title,
		Content:     dto.Content,
		Description: dto.Description,
		Version:     dto.Version,
		Alias:       dto.Alias,
		ParentId:    parentId,
		IsDir:       dto.IsDir,
		Sort:        dto.Sort,
	})
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg(fmt.Sprintf("创建文档内容成功, 文档Id:%s", id.Hex()))
}

// FindDocumentContentById 管理员查询特定Id的文档内容
func (h *DocumentContentHandler) FindDocumentContentById(c *gin.Context) *apiwrap.Response[any] {
	documentId := c.Param("id")
	if documentId == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id不能为空")
	}

	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}

	doc, err := h.serv.FindDocumentContentById(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](DocumentContentVO{
		Id:           doc.DocumentId.Hex(),
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
		DocumentId:   doc.DocumentId.Hex(),
		Title:        doc.Title,
		Content:      doc.Content,
		Description:  doc.Description,
		Version:      doc.Version,
		Alias:        doc.Alias,
		ParentId:     doc.ParentId.Hex(),
		IsDir:        doc.IsDir,
		Sort:         doc.Sort,
		LikeCount:    doc.LikeCount,
		DislikeCount: doc.DislikeCount,
		CommentCount: doc.CommentCount,
	}, fmt.Sprintf("查询文档内容成功, 文档Id:%s", documentId))
}

// DeleteDocumentContentById 管理员删除特定Id的文档内容
func (h *DocumentContentHandler) DeleteDocumentContentById(c *gin.Context) *apiwrap.Response[any] {
	documentId := c.Param("id")
	if documentId == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id不能为空")
	}
	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}
	err = h.serv.DeleteDocumentContentById(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg(fmt.Sprintf("删除文档内容成功, 文档Id:%s", documentId))
}

// FindDocumentContentByParentId 管理员根据父级Id查询所有子文档内容
func (h *DocumentContentHandler) FindDocumentContentByParentId(c *gin.Context) *apiwrap.Response[any] {
	parentId := c.Query("parentId")
	if parentId == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "parentId不能为空")
	}
	objId, err := bson.ObjectIDFromHex(parentId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "parentId格式错误")
	}
	docs, err := h.serv.FindDocumentContentByParentId(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	docsVO := make([]DocumentContentVO, len(docs))
	for i, doc := range docs {
		docsVO[i] = DocumentContentVO{
			Id:           doc.DocumentId.Hex(),
			CreatedAt:    doc.CreatedAt,
			UpdatedAt:    doc.UpdatedAt,
			DocumentId:   doc.DocumentId.Hex(),
			Title:        doc.Title,
			Content:      doc.Content,
			Description:  doc.Description,
			Version:      doc.Version,
			Alias:        doc.Alias,
			ParentId:     doc.ParentId.Hex(),
			IsDir:        doc.IsDir,
			Sort:         doc.Sort,
			LikeCount:    doc.LikeCount,
			DislikeCount: doc.DislikeCount,
			CommentCount: doc.CommentCount,
		}
	}
	return apiwrap.SuccessWithDetail[any](docsVO, fmt.Sprintf("查询文档内容成功, 父级Id:%s", parentId))
}

// FindDocumentContentByDocumentId 管理员根据文档Id查询所有子文档内容
func (h *DocumentContentHandler) FindDocumentContentByDocumentId(c *gin.Context) *apiwrap.Response[any] {
	documentId := c.Query("documentId")
	if documentId == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "documentId不能为空")
	}
	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "documentId格式错误")
	}
	docs, err := h.serv.FindDocumentContentByDocumentId(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	docsVO := make([]DocumentContentVO, len(docs))
	for i, doc := range docs {
		docsVO[i] = DocumentContentVO{
			Id:           doc.DocumentId.Hex(),
			CreatedAt:    doc.CreatedAt,
			UpdatedAt:    doc.UpdatedAt,
			DocumentId:   doc.DocumentId.Hex(),
			Title:        doc.Title,
			Content:      doc.Content,
			Description:  doc.Description,
			Version:      doc.Version,
			Alias:        doc.Alias,
			ParentId:     doc.ParentId.Hex(),
			IsDir:        doc.IsDir,
			Sort:         doc.Sort,
			LikeCount:    doc.LikeCount,
			DislikeCount: doc.DislikeCount,
			CommentCount: doc.CommentCount,
		}
	}
	return apiwrap.SuccessWithDetail[any](docsVO, fmt.Sprintf("查询文档内容成功, 文档Id:%s", documentId))
}

// UpdateDocumentContentById 管理员更新特定Id的文档内容
func (h *DocumentContentHandler) UpdateDocumentContentById(c *gin.Context, dto UpdateDocumentContentDto) *apiwrap.Response[any] {
	objId, _ := bson.ObjectIDFromHex(dto.Id)
	parentId, _ := bson.ObjectIDFromHex(dto.ParentId)
	documentId, _ := bson.ObjectIDFromHex(dto.DocumentId)

	err := h.serv.UpdateDocumentContentById(c, objId, domain.DocumentContent{
		DocumentId:  documentId,
		Title:       dto.Title,
		Content:     dto.Content,
		Description: dto.Description,
		Version:     dto.Version,
		Alias:       dto.Alias,
		ParentId:    parentId,
		IsDir:       dto.IsDir,
		Sort:        dto.Sort,
	})
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg(fmt.Sprintf("更新文档内容成功, 文档Id:%s", dto.Id))
}
