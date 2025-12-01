package web

import (
	"fmt"

	"github.com/codepzj/Stellux-Server/internal/document_content/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/document_content/internal/service"
	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"
	"github.com/codepzj/Stellux-Server/internal/pkg/middleware"
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
		adminDocumentContentGroup.POST("/create", apiwrap.WrapWithJson(h.CreateDocumentContent))         // 管理员创建文档内容
		adminDocumentContentGroup.GET("/:id", apiwrap.Wrap(h.FindDocumentContentById))                   // 管理员查询特定Id的文档内容
		adminDocumentContentGroup.DELETE("/:id", apiwrap.Wrap(h.DeleteDocumentContentById))              // 管理员删除特定Id的文档内容
		adminDocumentContentGroup.PUT("/soft-delete/:id", apiwrap.Wrap(h.SoftDeleteDocumentContentById)) // 管理员软删除特定Id的文档内容
		adminDocumentContentGroup.PUT("/restore/:id", apiwrap.Wrap(h.RestoreDocumentContentById))        // 管理员恢复特定Id的文档内容
		adminDocumentContentGroup.GET("/all/parent-id", apiwrap.Wrap(h.FindDocumentContentByParentId))   // 管理员查询特定父级Id的所有子文档内容
		adminDocumentContentGroup.GET("/all", apiwrap.Wrap(h.FindDocumentContentByDocumentId))           // 管理员查询特定文档Id的所有子文档内容
		adminDocumentContentGroup.PUT("/update", apiwrap.WrapWithJson(h.UpdateDocumentContentById))      // 管理员更新特定Id的文档内容
		adminDocumentContentGroup.GET("/list", apiwrap.WrapWithQuery(h.GetDocumentContentList))          // 管理员获取文档内容列表
		adminDocumentContentGroup.GET("/search", apiwrap.Wrap(h.SearchDocumentContent))                  // 管理员搜索文档内容
		adminDocumentContentGroup.POST("/delete-list", apiwrap.Wrap(h.DeleteDocumentContentList))        // 管理员批量删除文档内容
	}

	// 公开API
	documentContentGroup := engine.Group("/document-content")
	{
		documentContentGroup.GET("/:id", apiwrap.Wrap(h.FindPublicDocumentContentById))         // 公开查询特定Id的文档内容
		documentContentGroup.GET("/all", apiwrap.Wrap(h.FindPublicDocumentContentByDocumentId)) // 公开查询特定文档Id的所有子文档内容
		// documentContentGroup.GET("/search", apiwrap.Wrap(h.SearchPublicDocumentContent))                          // 公开搜索文档内容
		documentContentGroup.GET("/by-root-and-alias", apiwrap.Wrap(h.FindPublicDocumentContentByRootIdAndAlias)) // 公开根据根文档ID和别名查询文档内容
	}
}

// CreateDocumentContent 管理员创建文档内容
func (h *DocumentContentHandler) CreateDocumentContent(c *gin.Context, dto CreateDocumentContentRequest) (int, string, any) {
	documentId, _ := bson.ObjectIDFromHex(dto.DocumentId)
	parentId, _ := bson.ObjectIDFromHex(dto.ParentId)
	if !dto.IsDir && dto.Alias == "" {
		return 400, "alias不能为空", nil
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
		return 500, err.Error(), nil
	}
	return 200, fmt.Sprintf("创建文档内容成功, 文档Id:%s", id.Hex()), nil
}

// FindDocumentContentById 管理员查询特定Id的文档内容
func (h *DocumentContentHandler) FindDocumentContentById(c *gin.Context) (int, string, any) {
	documentId := c.Param("id")
	if documentId == "" {
		return 400, "id不能为空", nil
	}

	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return 400, "id格式错误", nil
	}

	doc, err := h.serv.FindDocumentContentById(c, objId)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, fmt.Sprintf("查询文档内容成功, 文档Id:%s", documentId), h.DocumentContentDomainToVO(doc)
}

// DeleteDocumentContentById 管理员删除特定Id的文档内容
func (h *DocumentContentHandler) DeleteDocumentContentById(c *gin.Context) (int, string, any) {
	documentId := c.Param("id")
	if documentId == "" {
		return 400, "id不能为空", nil
	}
	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return 400, "id格式错误", nil
	}
	err = h.serv.DeleteDocumentContentById(c, objId)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, fmt.Sprintf("删除文档内容成功, 文档Id:%s", documentId), nil
}

// SoftDeleteDocumentContentById 管理员软删除特定Id的文档内容
func (h *DocumentContentHandler) SoftDeleteDocumentContentById(c *gin.Context) (int, string, any) {
	documentId := c.Param("id")
	if documentId == "" {
		return 400, "id不能为空", nil
	}
	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return 400, "id格式错误", nil
	}
	err = h.serv.SoftDeleteDocumentContentById(c, objId)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, fmt.Sprintf("软删除文档内容成功, 文档Id:%s", documentId), nil
}

// RestoreDocumentContentById 管理员恢复特定Id的文档内容
func (h *DocumentContentHandler) RestoreDocumentContentById(c *gin.Context) (int, string, any) {
	documentId := c.Param("id")
	if documentId == "" {
		return 400, "id不能为空", nil
	}
	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return 400, "id格式错误", nil
	}
	err = h.serv.RestoreDocumentContentById(c, objId)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, fmt.Sprintf("恢复文档内容成功, 文档Id:%s", documentId), nil
}

// FindDocumentContentByParentId 管理员根据父级Id查询所有子文档内容
func (h *DocumentContentHandler) FindDocumentContentByParentId(c *gin.Context) (int, string, any) {
	parentId := c.Query("parent_id")
	if parentId == "" {
		return 400, "parent_id不能为空", nil
	}
	objId, err := bson.ObjectIDFromHex(parentId)
	if err != nil {
		return 400, "parentId格式错误", nil
	}
	docs, err := h.serv.FindDocumentContentByParentId(c, objId)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, fmt.Sprintf("查询文档内容成功, 父级Id:%s", parentId), h.DocumentContentDomainToVOList(docs)
}

// FindDocumentContentByDocumentId 管理员根据文档Id查询所有子文档内容
func (h *DocumentContentHandler) FindDocumentContentByDocumentId(c *gin.Context) (int, string, any) {
	documentId := c.Query("document_id")
	if documentId == "" {
		return 400, "document_id不能为空", nil
	}
	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return 400, "documentId格式错误", nil
	}
	docs, err := h.serv.FindDocumentContentByDocumentId(c, objId)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, fmt.Sprintf("查询文档内容成功, 文档Id:%s", documentId), h.DocumentContentDomainToVOList(docs)
}

// UpdateDocumentContentById 管理员更新特定Id的文档内容
func (h *DocumentContentHandler) UpdateDocumentContentById(c *gin.Context, dto UpdateDocumentContentRequest) (int, string, any) {
	objId, _ := bson.ObjectIDFromHex(dto.Id)
	documentId, _ := bson.ObjectIDFromHex(dto.DocumentId)

	if !dto.IsDir && dto.Alias == "" {
		return 400, "alias不能为空", nil
	}

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
		return 500, err.Error(), nil
	}
	return 200, fmt.Sprintf("更新文档内容成功, 文档Id:%s", dto.Id), nil
}

// GetDocumentContentList 管理员获取文档内容列表
func (h *DocumentContentHandler) GetDocumentContentList(c *gin.Context, page apiwrap.Page) (int, string, any) {
	// 从查询参数获取documentId
	documentIdStr := c.Query("document_id")
	documentId, err := bson.ObjectIDFromHex(documentIdStr)
	if err != nil {
		return 400, "document_id格式错误", nil
	}

	docs, count, err := h.serv.GetDocumentContentList(c, &page, documentId)
	if err != nil {
		return 500, err.Error(), nil
	}

	docsVO := h.DocumentContentDomainToVOList(docs)
	// 转换为指针切片
	docsVOPtr := make([]*DocumentContentVO, len(docsVO))
	for i := range docsVO {
		docsVOPtr[i] = &docsVO[i]
	}

	return 200, "获取文档内容列表成功", apiwrap.ToPageVO(page.PageNo, page.PageSize, count, docsVOPtr)
}

// SearchDocumentContent 管理员搜索文档内容
func (h *DocumentContentHandler) SearchDocumentContent(c *gin.Context) (int, string, any) {
	keyword := c.Query("keyword")
	if keyword == "" {
		return 400, "keyword不能为空", nil
	}

	docs, err := h.serv.SearchDocumentContent(c, keyword)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, "搜索文档内容成功", h.DocumentContentDomainToVOList(docs)
}

// FindPublicDocumentContentById 公开查询特定Id的文档内容
func (h *DocumentContentHandler) FindPublicDocumentContentById(c *gin.Context) (int, string, any) {
	documentId := c.Param("id")
	if documentId == "" {
		return 400, "id不能为空", nil
	}

	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return 400, "id格式错误", nil
	}

	doc, err := h.serv.FindPublicDocumentContentById(c, objId)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, fmt.Sprintf("查询文档内容成功, 文档Id:%s", documentId), h.DocumentContentDomainToVO(doc)
}

// FindPublicDocumentContentByParentId 公开根据父级Id查询所有子文档内容
func (h *DocumentContentHandler) FindPublicDocumentContentByParentId(c *gin.Context) (int, string, any) {
	parentId := c.Query("parent_id")
	if parentId == "" {
		return 400, "parent_id不能为空", nil
	}
	objId, err := bson.ObjectIDFromHex(parentId)
	if err != nil {
		return 400, "parentId格式错误", nil
	}
	docs, err := h.serv.FindPublicDocumentContentByParentId(c, objId)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, fmt.Sprintf("查询文档内容成功, 父级Id:%s", parentId), h.DocumentContentDomainToVOList(docs)
}

// FindPublicDocumentContentByDocumentId 公开根据文档Id查询所有子文档内容
func (h *DocumentContentHandler) FindPublicDocumentContentByDocumentId(c *gin.Context) (int, string, any) {
	documentId := c.Query("document_id")
	if documentId == "" {
		return 400, "document_id不能为空", nil
	}
	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return 400, "documentId格式错误", nil
	}
	docs, err := h.serv.FindPublicDocumentContentByDocumentId(c, objId)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, fmt.Sprintf("查询文档内容成功, 文档Id:%s", documentId), h.DocumentContentDomainToVOList(docs)
}

// FindPublicDocumentContentByRootIdAndAlias 公开根据根文档ID和别名查询文档内容
func (h *DocumentContentHandler) FindPublicDocumentContentByRootIdAndAlias(c *gin.Context) (int, string, any) {
	documentId := c.Query("document_id")
	if documentId == "" {
		return 400, "document_id不能为空", nil
	}

	alias := c.Query("alias")
	if alias == "" {
		return 400, "alias不能为空", nil
	}

	objId, err := bson.ObjectIDFromHex(documentId)
	if err != nil {
		return 400, "document_id格式错误", nil
	}

	docs, err := h.serv.FindPublicDocumentContentByDocumentId(c, objId)
	if err != nil {
		return 500, err.Error(), nil
	}

	// 在获取的文档列表中查找匹配别名的文档
	var targetDoc domain.DocumentContent
	found := false
	for _, doc := range docs {
		if doc.Alias == alias {
			targetDoc = doc
			found = true
			break
		}
	}

	if !found {
		return 400, "未找到指定别名的文档", nil
	}

	return 200, fmt.Sprintf("查询文档内容成功, 根文档ID:%s, 别名:%s", documentId, alias), h.DocumentContentDomainToVO(targetDoc)
}

// DocumentContentDomainToVO 将domain对象转换为VO
func (h *DocumentContentHandler) DocumentContentDomainToVO(doc domain.DocumentContent) DocumentContentVO {
	return DocumentContentVO{
		Id:          doc.Id.Hex(),
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
		DocumentId:  doc.DocumentId.Hex(),
		Title:       doc.Title,
		Content:     doc.Content,
		Description: doc.Description,
		Alias:       doc.Alias,
		ParentId:    doc.ParentId.Hex(),
		IsDir:       doc.IsDir,
		Sort:        doc.Sort,
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
func (h *DocumentContentHandler) DeleteDocumentContentList(c *gin.Context) (int, string, any) {
	var req struct {
		IdList []string `json:"id_list"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return 400, "参数错误", nil
	}
	if len(req.IdList) == 0 {
		return 400, "id_list不能为空", nil
	}
	if err := h.serv.DeleteDocumentContentList(c, req.IdList); err != nil {
		return 500, err.Error(), nil
	}
	return 200, "批量删除成功", nil
}
