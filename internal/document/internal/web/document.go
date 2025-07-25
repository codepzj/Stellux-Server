package web

import (
	"fmt"
	"net/http"

	"github.com/codepzj/stellux/server/internal/document/internal/domain"
	docService "github.com/codepzj/stellux/server/internal/document/internal/service"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/codepzj/stellux/server/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewDocumentHandler(serv docService.IDocumentService) *DocumentHandler {
	return &DocumentHandler{
		serv: serv,
	}
}

type DocumentHandler struct {
	serv docService.IDocumentService
}

func (h *DocumentHandler) RegisterGinRoutes(engine *gin.Engine) {
	adminDocumentGroup := engine.Group("/admin-api/document")
	{
		adminDocumentGroup.Use(middleware.JWT())
		adminDocumentGroup.POST("/create", apiwrap.WrapWithJson(h.AdminCreateDocument))     // 管理员创建文档
		adminDocumentGroup.GET("/find", apiwrap.Wrap(h.AdminFindDocument))                  // 管理员查询特定Id的文档
		adminDocumentGroup.PUT("/update", apiwrap.WrapWithJson(h.AdminUpdateDocument))      // 管理员更新文档
		adminDocumentGroup.DELETE("/delete/:id", apiwrap.Wrap(h.AdminDeleteDocument))       // 管理员删除文档
		adminDocumentGroup.PUT("/soft-delete/:id", apiwrap.Wrap(h.AdminSoftDeleteDocument)) // 管理员软删除文档
		adminDocumentGroup.PUT("/restore/:id", apiwrap.Wrap(h.AdminRestoreDocument))        // 管理员恢复文档
		adminDocumentGroup.GET("/find-by-alias", apiwrap.Wrap(h.AdminFindDocumentByAlias))  // 管理员根据别名查询文档
		adminDocumentGroup.GET("/list", apiwrap.WrapWithQuery(h.AdminGetDocumentList))      // 管理员获取文档列表
	}

	// 公开API
	documentGroup := engine.Group("/document")
	{
		documentGroup.GET("/:id", apiwrap.Wrap(h.GetDocument))                  // 获取根文档
		documentGroup.GET("/all", apiwrap.Wrap(h.GetAllPublicDocument))         // 获取所有公开文档
		documentGroup.GET("/find", apiwrap.Wrap(h.FindDocument))                // 公开查询特定Id的文档
		documentGroup.GET("/alias/:alias", apiwrap.Wrap(h.FindDocumentByAlias)) // 公开根据别名查询文档
		documentGroup.GET("/list", apiwrap.WrapWithQuery(h.GetDocumentList))    // 公开获取文档列表
	}
}

// AdminCreateDocument 管理员创建文档
func (h *DocumentHandler) AdminCreateDocument(c *gin.Context, req DocumentCreateRequest) *apiwrap.Response[any] {
	id, err := h.serv.CreateDocument(c, domain.Document{
		Title:       req.Title,
		Description: req.Description,
		Thumbnail:   req.Thumbnail,
		Alias:       req.Alias,
		Sort:        req.Sort,
		IsPublic:    req.IsPublic,
	})
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg(fmt.Sprintf("创建文档成功, Id:%s", id.Hex()))
}

// AdminFindDocument 管理员查询特定Id的文档
func (h *DocumentHandler) AdminFindDocument(c *gin.Context) *apiwrap.Response[any] {
	id := c.Query("id")

	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}

	doc, err := h.serv.FindDocumentById(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, fmt.Sprintf("查询文档失败, Id:%s, err:%s", id, err.Error()))
	}
	docVO := DocumentVO{
		Id:          doc.Id.Hex(),
		Title:       doc.Title,
		Description: doc.Description,
		Thumbnail:   doc.Thumbnail,
		Alias:       doc.Alias,
		Sort:        doc.Sort,
		IsPublic:    doc.IsPublic,
		IsDeleted:   doc.IsDeleted,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}
	return apiwrap.SuccessWithDetail[any](docVO, "查询文档成功")
}

// AdminUpdateDocument 管理员更新文档
func (h *DocumentHandler) AdminUpdateDocument(c *gin.Context, req DocumentUpdateRequest) *apiwrap.Response[any] {
	objId, err := bson.ObjectIDFromHex(req.Id)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}

	err = h.serv.UpdateDocumentById(c, objId, domain.Document{
		Title:       req.Title,
		Description: req.Description,
		Thumbnail:   req.Thumbnail,
		Alias:       req.Alias,
		Sort:        req.Sort,
		IsPublic:    req.IsPublic,
	})
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新文档成功")
}

// AdminDeleteDocument 管理员删除文档
func (h *DocumentHandler) AdminDeleteDocument(c *gin.Context) *apiwrap.Response[any] {
	id := c.Param("id")
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}

	err = h.serv.DeleteDocumentById(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("删除文档成功")
}

// AdminSoftDeleteDocument 管理员软删除文档
func (h *DocumentHandler) AdminSoftDeleteDocument(c *gin.Context) *apiwrap.Response[any] {
	id := c.Param("id")
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}

	err = h.serv.SoftDeleteDocumentById(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("软删除文档成功")
}

// AdminRestoreDocument 管理员恢复文档
func (h *DocumentHandler) AdminRestoreDocument(c *gin.Context) *apiwrap.Response[any] {
	id := c.Param("id")
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}

	err = h.serv.RestoreDocumentById(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("恢复文档成功")
}

// AdminFindDocumentByAlias 管理员根据别名查询文档
func (h *DocumentHandler) AdminFindDocumentByAlias(c *gin.Context) *apiwrap.Response[any] {
	alias := c.Query("alias")
	if alias == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "alias不能为空")
	}

	doc, err := h.serv.FindDocumentByAlias(c, alias, bson.D{})
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	docVO := DocumentVO{
		Id:          doc.Id.Hex(),
		Title:       doc.Title,
		Description: doc.Description,
		Thumbnail:   doc.Thumbnail,
		Alias:       doc.Alias,
		Sort:        doc.Sort,
		IsPublic:    doc.IsPublic,
		IsDeleted:   doc.IsDeleted,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}
	return apiwrap.SuccessWithDetail[any](docVO, "查询文档成功")
}

// AdminGetDocumentList 管理员获取文档列表
func (h *DocumentHandler) AdminGetDocumentList(c *gin.Context, page apiwrap.Page) *apiwrap.Response[any] {
	docs, count, err := h.serv.GetDocumentList(c, &domain.Page{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	})
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}

	docsVO := h.DocumentDomainToVOList(docs)
	// 转换为指针切片
	docsVOPtr := make([]*DocumentVO, len(docsVO))
	for i := range docsVO {
		docsVOPtr[i] = &docsVO[i]
	}

	return apiwrap.SuccessWithDetail[any](apiwrap.ToPageVO(page.PageNo, page.PageSize, count, docsVOPtr), "获取文档列表成功")
}

// FindDocument 公开查询特定Id的文档
func (h *DocumentHandler) FindDocument(c *gin.Context) *apiwrap.Response[any] {
	id := c.Query("id")

	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}

	doc, err := h.serv.FindDocumentById(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, fmt.Sprintf("查询文档失败, Id:%s, err:%s", id, err.Error()))
	}

	// 检查文档是否公开
	if !doc.IsPublic {
		return apiwrap.FailWithMsg(http.StatusForbidden, "文档未公开")
	}

	docVO := DocumentVO{
		Id:          doc.Id.Hex(),
		Title:       doc.Title,
		Description: doc.Description,
		Alias:       doc.Alias,
		Sort:        doc.Sort,
		IsPublic:    doc.IsPublic,
		IsDeleted:   doc.IsDeleted,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}
	return apiwrap.SuccessWithDetail[any](docVO, "查询文档成功")
}

// FindDocumentByAlias 公开根据别名查询文档
func (h *DocumentHandler) FindDocumentByAlias(c *gin.Context) *apiwrap.Response[any] {
	alias := c.Param("alias")
	if alias == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "alias不能为空")
	}

	doc, err := h.serv.FindDocumentByAlias(c, alias, bson.D{{Key: "is_public", Value: true}, {Key: "is_deleted", Value: false}})
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	docVO := DocumentVO{
		Id:          doc.Id.Hex(),
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
		Title:       doc.Title,
		Description: doc.Description,
		Thumbnail:   doc.Thumbnail,
		Alias:       doc.Alias,
		Sort:        doc.Sort,
		IsPublic:    doc.IsPublic,
		IsDeleted:   doc.IsDeleted,
	}
	return apiwrap.SuccessWithDetail[any](docVO, "查询文档成功")
}

// GetDocumentList 公开获取文档列表
func (h *DocumentHandler) GetDocumentList(c *gin.Context, page apiwrap.Page) *apiwrap.Response[any] {
	docs, count, err := h.serv.GetPublicDocumentList(c, &domain.Page{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	})
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}

	docsVO := h.DocumentDomainToVOList(docs)
	// 转换为指针切片
	docsVOPtr := make([]*DocumentVO, len(docsVO))
	for i := range docsVO {
		docsVOPtr[i] = &docsVO[i]
	}

	return apiwrap.SuccessWithDetail[any](apiwrap.ToPageVO(page.PageNo, page.PageSize, count, docsVOPtr), "获取文档列表成功")
}

// GetAllPublicDocument 获取所有公开文档
func (h *DocumentHandler) GetAllPublicDocument(c *gin.Context) *apiwrap.Response[any] {
	// 获取所有公开文档，不使用分页
	docs, err := h.serv.GetAllPublicDocuments(c)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}

	// 转换为DocumentVO数组
	docsVO := make([]DocumentVO, len(docs))
	for i, doc := range docs {
		docsVO[i] = DocumentVO{
			Id:          doc.Id.Hex(),
			Title:       doc.Title,
			Description: doc.Description,
			Thumbnail:   doc.Thumbnail,
			Alias:       doc.Alias,
			IsPublic:    doc.IsPublic,
			CreatedAt:   doc.CreatedAt,
			UpdatedAt:   doc.UpdatedAt,
		}
	}

	// 直接返回文档列表，不包装成分页格式
	return apiwrap.SuccessWithDetail[any](docsVO, "获取所有公开文档成功")
}

// GetDocument 获取文档
func (h *DocumentHandler) GetDocument(c *gin.Context) *apiwrap.Response[any] {
	id := c.Param("id")
	if id == "" {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id不能为空")
	}

	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusBadRequest, "id格式错误")
	}

	doc, err := h.serv.FindDocumentById(c, objId)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}

	// 检查文档是否公开
	if !doc.IsPublic {
		return apiwrap.FailWithMsg(http.StatusForbidden, "文档未公开")
	}

	docVO := DocumentVO{
		Id:          doc.Id.Hex(),
		Title:       doc.Title,
		Description: doc.Description,
		Alias:       doc.Alias,
		Thumbnail:   doc.Thumbnail,
		IsPublic:    doc.IsPublic,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}

	return apiwrap.SuccessWithDetail[any](docVO, "获取文档成功")
}

// DocumentDomainToVOList 将domain对象转换为VO列表
func (h *DocumentHandler) DocumentDomainToVOList(docs []domain.Document) []DocumentVO {
	docsVO := make([]DocumentVO, len(docs))
	for i, doc := range docs {
		docsVO[i] = DocumentVO{
			Id:          doc.Id.Hex(),
			Title:       doc.Title,
			Description: doc.Description,
			Thumbnail:   doc.Thumbnail,
			Alias:       doc.Alias,
			Sort:        doc.Sort,
			IsPublic:    doc.IsPublic,
			IsDeleted:   doc.IsDeleted,
			CreatedAt:   doc.CreatedAt,
			UpdatedAt:   doc.UpdatedAt,
		}
	}
	return docsVO
}
