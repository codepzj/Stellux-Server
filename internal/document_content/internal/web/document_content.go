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
	adminDocumentContentGroup := engine.Group("/admin-api/document-content")
	{
		adminDocumentContentGroup.Use(middleware.JWT())
		adminDocumentContentGroup.POST("/create", apiwrap.WrapWithJson(h.CreateDocumentContent))           // 管理员创建文档内容
		adminDocumentContentGroup.GET("/:id", apiwrap.Wrap(h.FindDocumentContentById))                     // 管理员查询特定Id的文档内容
		adminDocumentContentGroup.DELETE("/:id", apiwrap.Wrap(h.DeleteDocumentContentById))                // 管理员删除特定Id的文档内容
		adminDocumentContentGroup.PUT("/soft-delete/:id", apiwrap.Wrap(h.SoftDeleteDocumentContentById))   // 管理员软删除特定Id的文档内容
		adminDocumentContentGroup.PUT("/restore/:id", apiwrap.Wrap(h.RestoreDocumentContentById))          // 管理员恢复特定Id的文档内容
		adminDocumentContentGroup.GET("/all/parent-id", apiwrap.Wrap(h.FindDocumentContentByParentId))     // 管理员查询特定父级Id的所有子文档内容
		adminDocumentContentGroup.GET("/all/document-id", apiwrap.Wrap(h.FindDocumentContentByDocumentId)) // 管理员查询特定文档Id的所有子文档内容
		adminDocumentContentGroup.PUT("/update", apiwrap.WrapWithJson(h.UpdateDocumentContentById))        // 管理员更新特定Id的文档内容
		adminDocumentContentGroup.GET("/list", apiwrap.WrapWithQuery(h.GetDocumentContentList))            // 管理员获取文档内容列表
		adminDocumentContentGroup.GET("/search", apiwrap.Wrap(h.SearchDocumentContent))                    // 管理员搜索文档内容
		adminDocumentContentGroup.PUT("/update-like/:id", apiwrap.Wrap(h.UpdateLikeCount))                 // 管理员更新点赞数
		adminDocumentContentGroup.PUT("/update-dislike/:id", apiwrap.Wrap(h.UpdateDislikeCount))           // 管理员更新反对数
		adminDocumentContentGroup.PUT("/update-comment/:id", apiwrap.Wrap(h.UpdateCommentCount))           // 管理员更新评论数
		adminDocumentContentGroup.POST("/delete-list", apiwrap.Wrap(h.DeleteDocumentContentList))          // 管理员批量删除文档内容
	}

	// 公开API
	publicDocumentContentGroup := engine.Group("/document-content")
	{
		publicDocumentContentGroup.GET("/:id", apiwrap.Wrap(h.FindPublicDocumentContentById))                     // 公开查询特定Id的文档内容
		publicDocumentContentGroup.GET("/all/parent-id", apiwrap.Wrap(h.FindPublicDocumentContentByParentId))     // 公开查询特定父级Id的所有子文档内容
		publicDocumentContentGroup.GET("/all/document-id", apiwrap.Wrap(h.FindPublicDocumentContentByDocumentId)) // 公开查询特定文档Id的所有子文档内容
		publicDocumentContentGroup.GET("/list", apiwrap.WrapWithQuery(h.GetPublicDocumentContentList))            // 公开获取文档内容列表
		publicDocumentContentGroup.GET("/search", apiwrap.Wrap(h.SearchPublicDocumentContent))                    // 公开搜索文档内容
	}
}

// CreateDocumentContent 管理员创建文档内容
func (h *DocumentContentHandler) CreateDocumentContent(c *gin.Context, dto CreateDocumentContentRequest) *apiwrap.Response[any] {
	documentId, _ := bson.ObjectIDFromHex(dto.DocumentId)

	// 处理ParentId，如果为空则使用documentId作为父级ID
	var parentId bson.ObjectID
	if dto.ParentId == "" {
		parentId = documentId // 根目录的父级ID就是documentId
	} else {
		parentId, _ = bson.ObjectIDFromHex(dto.ParentId)
	}

	id, err := h.serv.CreateDocumentContent(c, domain.DocumentContent{
		DocumentId:  documentId,
		Title:       dto.Title,
		Content:     dto.Content,
		Description: dto.Description,
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
	return apiwrap.SuccessWithDetail[any](h.DocumentContentDomainToVO(doc), fmt.Sprintf("查询文档内容成功, 文档Id:%s", documentId))
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

// SoftDeleteDocumentContentById 管理员软删除特定Id的文档内容
func (h *DocumentContentHandler) SoftDeleteDocumentContentById(c *gin.Context) *apiwrap.Response[any] {
	documentId := c.Param("id")
	if documentId == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id不能为空")
	}
	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}
	err = h.serv.SoftDeleteDocumentContentById(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg(fmt.Sprintf("软删除文档内容成功, 文档Id:%s", documentId))
}

// RestoreDocumentContentById 管理员恢复特定Id的文档内容
func (h *DocumentContentHandler) RestoreDocumentContentById(c *gin.Context) *apiwrap.Response[any] {
	documentId := c.Param("id")
	if documentId == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id不能为空")
	}
	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}
	err = h.serv.RestoreDocumentContentById(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg(fmt.Sprintf("恢复文档内容成功, 文档Id:%s", documentId))
}

// FindDocumentContentByParentId 管理员根据父级Id查询所有子文档内容
func (h *DocumentContentHandler) FindDocumentContentByParentId(c *gin.Context) *apiwrap.Response[any] {
	parentId := c.Query("parent_id")
	if parentId == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "parent_id不能为空")
	}
	objId, err := bson.ObjectIDFromHex(parentId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "parentId格式错误")
	}
	docs, err := h.serv.FindDocumentContentByParentId(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentContentDomainToVOList(docs), fmt.Sprintf("查询文档内容成功, 父级Id:%s", parentId))
}

// FindDocumentContentByDocumentId 管理员根据文档Id查询所有子文档内容
func (h *DocumentContentHandler) FindDocumentContentByDocumentId(c *gin.Context) *apiwrap.Response[any] {
	documentId := c.Query("document_id")
	if documentId == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "document_id不能为空")
	}
	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "documentId格式错误")
	}
	docs, err := h.serv.FindDocumentContentByDocumentId(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentContentDomainToVOList(docs), fmt.Sprintf("查询文档内容成功, 文档Id:%s", documentId))
}

// UpdateDocumentContentById 管理员更新特定Id的文档内容
func (h *DocumentContentHandler) UpdateDocumentContentById(c *gin.Context, dto UpdateDocumentContentRequest) *apiwrap.Response[any] {
	objId, _ := bson.ObjectIDFromHex(dto.Id)
	documentId, _ := bson.ObjectIDFromHex(dto.DocumentId)

	// 处理ParentId，如果为空则使用documentId作为父级ID
	var parentId bson.ObjectID
	if dto.ParentId == "" {
		parentId = documentId // 根目录的父级ID就是documentId
	} else {
		parentId, _ = bson.ObjectIDFromHex(dto.ParentId)
	}

	err := h.serv.UpdateDocumentContentById(c, objId, domain.DocumentContent{
		DocumentId:  documentId,
		Title:       dto.Title,
		Content:     dto.Content,
		Description: dto.Description,
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

// GetDocumentContentList 管理员获取文档内容列表
func (h *DocumentContentHandler) GetDocumentContentList(c *gin.Context, page apiwrap.Page) *apiwrap.Response[any] {
	docs, count, err := h.serv.GetDocumentContentList(c, &domain.Page{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	})
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}

	docsVO := h.DocumentContentDomainToVOList(docs)
	// 转换为指针切片
	docsVOPtr := make([]*DocumentContentVO, len(docsVO))
	for i := range docsVO {
		docsVOPtr[i] = &docsVO[i]
	}

	return apiwrap.SuccessWithDetail[any](apiwrap.ToPageVO(page.PageNo, page.PageSize, count, docsVOPtr), "获取文档内容列表成功")
}

// SearchDocumentContent 管理员搜索文档内容
func (h *DocumentContentHandler) SearchDocumentContent(c *gin.Context) *apiwrap.Response[any] {
	keyword := c.Query("keyword")
	if keyword == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "keyword不能为空")
	}

	docs, err := h.serv.SearchDocumentContent(c, keyword)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentContentDomainToVOList(docs), "搜索文档内容成功")
}

// UpdateLikeCount 管理员更新点赞数
func (h *DocumentContentHandler) UpdateLikeCount(c *gin.Context) *apiwrap.Response[any] {
	id := c.Param("id")
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}

	err = h.serv.UpdateLikeCount(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新点赞数成功")
}

// UpdateDislikeCount 管理员更新反对数
func (h *DocumentContentHandler) UpdateDislikeCount(c *gin.Context) *apiwrap.Response[any] {
	id := c.Param("id")
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}

	err = h.serv.UpdateDislikeCount(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新反对数成功")
}

// UpdateCommentCount 管理员更新评论数
func (h *DocumentContentHandler) UpdateCommentCount(c *gin.Context) *apiwrap.Response[any] {
	id := c.Param("id")
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}

	err = h.serv.UpdateCommentCount(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新评论数成功")
}

// FindPublicDocumentContentById 公开查询特定Id的文档内容
func (h *DocumentContentHandler) FindPublicDocumentContentById(c *gin.Context) *apiwrap.Response[any] {
	documentId := c.Param("id")
	if documentId == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id不能为空")
	}

	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}

	doc, err := h.serv.FindPublicDocumentContentById(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentContentDomainToVO(doc), fmt.Sprintf("查询文档内容成功, 文档Id:%s", documentId))
}

// FindPublicDocumentContentByParentId 公开根据父级Id查询所有子文档内容
func (h *DocumentContentHandler) FindPublicDocumentContentByParentId(c *gin.Context) *apiwrap.Response[any] {
	parentId := c.Query("parent_id")
	if parentId == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "parent_id不能为空")
	}
	objId, err := bson.ObjectIDFromHex(parentId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "parentId格式错误")
	}
	docs, err := h.serv.FindPublicDocumentContentByParentId(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentContentDomainToVOList(docs), fmt.Sprintf("查询文档内容成功, 父级Id:%s", parentId))
}

// FindPublicDocumentContentByDocumentId 公开根据文档Id查询所有子文档内容
func (h *DocumentContentHandler) FindPublicDocumentContentByDocumentId(c *gin.Context) *apiwrap.Response[any] {
	documentId := c.Query("document_id")
	if documentId == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "document_id不能为空")
	}
	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "documentId格式错误")
	}
	docs, err := h.serv.FindPublicDocumentContentByDocumentId(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentContentDomainToVOList(docs), fmt.Sprintf("查询文档内容成功, 文档Id:%s", documentId))
}

// GetPublicDocumentContentList 公开获取文档内容列表
func (h *DocumentContentHandler) GetPublicDocumentContentList(c *gin.Context, page apiwrap.Page) *apiwrap.Response[any] {
	docs, count, err := h.serv.GetPublicDocumentContentList(c, &domain.Page{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	})
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}

	docsVO := h.DocumentContentDomainToVOList(docs)
	// 转换为指针切片
	docsVOPtr := make([]*DocumentContentVO, len(docsVO))
	for i := range docsVO {
		docsVOPtr[i] = &docsVO[i]
	}

	return apiwrap.SuccessWithDetail[any](apiwrap.ToPageVO(page.PageNo, page.PageSize, count, docsVOPtr), "获取文档内容列表成功")
}

// SearchPublicDocumentContent 公开搜索文档内容
func (h *DocumentContentHandler) SearchPublicDocumentContent(c *gin.Context) *apiwrap.Response[any] {
	keyword := c.Query("keyword")
	if keyword == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "keyword不能为空")
	}

	docs, err := h.serv.SearchPublicDocumentContent(c, keyword)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentContentDomainToVOList(docs), "搜索文档内容成功")
}

// DocumentContentDomainToVO 将domain对象转换为VO
func (h *DocumentContentHandler) DocumentContentDomainToVO(doc domain.DocumentContent) DocumentContentVO {
	return DocumentContentVO{
		Id:           doc.Id.Hex(),
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
		DocumentId:   doc.DocumentId.Hex(),
		Title:        doc.Title,
		Content:      doc.Content,
		Description:  doc.Description,
		Alias:        doc.Alias,
		ParentId:     doc.ParentId.Hex(),
		IsDir:        doc.IsDir,
		Sort:         doc.Sort,
		LikeCount:    doc.LikeCount,
		DislikeCount: doc.DislikeCount,
		CommentCount: doc.CommentCount,
	}
}

// DocumentContentDomainToVOList 将domain对象转换为VO列表
func (h *DocumentContentHandler) DocumentContentDomainToVOList(docs []domain.DocumentContent) []DocumentContentVO {
	docsVO := make([]DocumentContentVO, len(docs))
	for i, doc := range docs {
		docsVO[i] = h.DocumentContentDomainToVO(doc)
	}
	return docsVO
}

// DeleteDocumentContentList 批量删除文档内容
func (h *DocumentContentHandler) DeleteDocumentContentList(c *gin.Context) *apiwrap.Response[any] {
	var req struct {
		DocumentIdList []string `json:"document_id_list"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "参数错误")
	}
	if len(req.DocumentIdList) == 0 {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "document_id_list不能为空")
	}
	if err := h.serv.DeleteDocumentContentList(c, req.DocumentIdList); err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("批量删除成功")
}
