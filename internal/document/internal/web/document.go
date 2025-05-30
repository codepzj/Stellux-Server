package web

import (
	"github.com/codepzj/stellux/server/global"
	"github.com/codepzj/stellux/server/internal/document/internal/domain"
	"github.com/codepzj/stellux/server/internal/document/internal/service"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func NewDocumentHandler(serv service.IDocumentService) *DocumentHandler {
	return &DocumentHandler{
		serv: serv,
	}
}

type DocumentHandler struct {
	serv service.IDocumentService
}

func (h *DocumentHandler) RegisterGinRoutes(engine *gin.Engine) {
	documentGroup := engine.Group("/document")
	{
		documentGroup.GET("/public", apiwrap.Wrap(h.GetAllPublicDoc))
		documentGroup.GET("/tree", apiwrap.Wrap(h.GetDocumentTreeByID))
		documentGroup.GET("/:id", apiwrap.Wrap(h.GetDocumentByID))
		documentGroup.GET("/root/:id", apiwrap.Wrap(h.GetRootDocumentByID))
	}
	adminGroup := engine.Group("/admin-api/document")
	{
		adminGroup.GET("/tree", apiwrap.Wrap(h.AdminGetDocumentTreeByID))
		adminGroup.GET("/root/:id", apiwrap.Wrap(h.AdminGetRootDocumentByID))
		adminGroup.GET("/:id", apiwrap.Wrap(h.AdminGetDocumentByID))
		adminGroup.GET("/list", apiwrap.Wrap(h.AdminGetAllRootDoc))
		adminGroup.GET("/parent-list", apiwrap.Wrap(h.AdminGetAllParentDoc))
		adminGroup.POST("/create-root", apiwrap.WrapWithBody(h.AdminCreateRootDocument))
		adminGroup.POST("/create", apiwrap.WrapWithBody(h.AdminCreateDocument))
		adminGroup.PUT("/edit-root", apiwrap.WrapWithBody(h.AdminEditRootDocument))
		adminGroup.PUT("/save", apiwrap.WrapWithBody(h.AdminSaveDocument))
		adminGroup.PUT("/rename", apiwrap.WrapWithBody(h.AdminRenameDocument))
		adminGroup.DELETE("/delete", apiwrap.WrapWithBody(h.AdminDeleteDocument))
		adminGroup.DELETE("/delete-list", apiwrap.WrapWithBody(h.AdminDeleteDocumentList))
		adminGroup.DELETE("/delete-root/:id", apiwrap.Wrap(h.AdminDeleteRootDocument))
	}
}

func (h *DocumentHandler) GetDocumentTreeByID(ctx *gin.Context) *apiwrap.Response[any] {
	documentID := ctx.Query("document_id")
	documentList, err := h.serv.FindAllPublicByDocumentID(ctx, apiwrap.ConvertBsonID(documentID))
	if err != nil {
		if errors.Is(err, global.ErrDocumentNotPublic) {
			return apiwrap.FailWithMsg(apiwrap.RequestDocumentNotPublic, err.Error())
		}
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentDomainListToTreeVOList(documentList), "获取公共文档目录树成功")
}

// 获取所有公共根文档
func (h *DocumentHandler) GetAllPublicDoc(ctx *gin.Context) *apiwrap.Response[any] {
	documentList, err := h.serv.FindAllPublic(ctx)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentDomainListToRootVOList(documentList), "获取公共根文档列表成功")
}

func (h *DocumentHandler) GetRootDocumentByID(ctx *gin.Context) *apiwrap.Response[any] {

	documentID := ctx.Param("id")
	document, err := h.serv.GetRootDocumentByID(ctx, apiwrap.ConvertBsonID(documentID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentRootDomainToVO(document), "获取根文档成功")
}

func (h *DocumentHandler) GetDocumentByID(ctx *gin.Context) *apiwrap.Response[any] {
	documentID := ctx.Param("id")
	document, err := h.serv.GetDocumentByID(ctx, apiwrap.ConvertBsonID(documentID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentDomainToVO(document), "获取文档成功")
}

// 新增文档
func (h *DocumentHandler) AdminCreateDocument(ctx *gin.Context, documentReq DocumentRequest) *apiwrap.Response[any] {
	err := h.serv.AdminCreate(ctx, h.DocumentRequestToDomain(documentReq))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.Success()
}

// 新增根文档
func (h *DocumentHandler) AdminCreateRootDocument(ctx *gin.Context, documentRootReq DocumentRootRequest) *apiwrap.Response[any] {
	err := h.serv.AdminCreateRoot(ctx, h.DocumentRequestToRootDomain(documentRootReq))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.Success()
}

// 编辑根文档
func (h *DocumentHandler) AdminEditRootDocument(ctx *gin.Context, documentRootReq DocumentRootEditRequest) *apiwrap.Response[any] {
	err := h.serv.AdminEditRootDocumentByID(ctx, apiwrap.ConvertBsonID(documentRootReq.ID), &domain.DocumentRoot{
		Title:        documentRootReq.Title,
		Alias:        documentRootReq.Alias,
		Description:  documentRootReq.Description,
		Thumbnail:    documentRootReq.Thumbnail,
		IsPublic:     documentRootReq.IsPublic,
		DocumentType: documentRootReq.DocumentType,
	})
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.Success()
}

// 删除根文档
func (h *DocumentHandler) AdminDeleteRootDocument(ctx *gin.Context) *apiwrap.Response[any] {
	documentID := ctx.Param("id")
	err := h.serv.AdminDeleteRootDocumentByID(ctx, apiwrap.ConvertBsonID(documentID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.Success()
}

// 保存文档
func (h *DocumentHandler) AdminSaveDocument(ctx *gin.Context, req UpdateDocumentRequest) *apiwrap.Response[any] {
	err := h.serv.AdminUpdateDocumentByID(ctx, apiwrap.ConvertBsonID(req.ID), req.Title, req.Content)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.Success()
}

// 重命名文档
func (h *DocumentHandler) AdminRenameDocument(ctx *gin.Context, req RenameDocumentRequest) *apiwrap.Response[any] {
	err := h.serv.AdminRenameDocumentByID(ctx, apiwrap.ConvertBsonID(req.ID), req.Title)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.Success()
}

// 删除文档
func (h *DocumentHandler) AdminDeleteDocument(ctx *gin.Context, req DeleteDocumentRequest) *apiwrap.Response[any] {
	err := h.serv.AdminDeleteByID(ctx, apiwrap.ConvertBsonID(req.DocumentID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.Success()
}

// 删除多个文档
func (h *DocumentHandler) AdminDeleteDocumentList(ctx *gin.Context, req DeleteDocumentListRequest) *apiwrap.Response[any] {
	err := h.serv.AdminDeleteByIDList(ctx, apiwrap.ConvertBsonIDList(req.DocumentIDList))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.Success()
}

// 获取文档目录树
func (h *DocumentHandler) AdminGetDocumentTreeByID(ctx *gin.Context) *apiwrap.Response[any] {
	documentID := ctx.Query("document_id")
	documentList, err := h.serv.AdminFindAllByDocumentID(ctx, apiwrap.ConvertBsonID(documentID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentDomainListToTreeVOList(documentList), "获取文档目录树成功")
}

// 获取根文档
func (h *DocumentHandler) AdminGetRootDocumentByID(ctx *gin.Context) *apiwrap.Response[any] {
	documentID := ctx.Param("id")
	document, err := h.serv.AdminGetRootDocumentByID(ctx, apiwrap.ConvertBsonID(documentID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentRootDomainToVO(document), "获取根文档成功")
}

// 获取文档
func (h *DocumentHandler) AdminGetDocumentByID(ctx *gin.Context) *apiwrap.Response[any] {
	documentID := ctx.Param("id")
	document, err := h.serv.AdminGetDocumentByID(ctx, apiwrap.ConvertBsonID(documentID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentDomainToVO(document), "获取文档成功")
}

// 获取所有根文档
func (h *DocumentHandler) AdminGetAllRootDoc(ctx *gin.Context) *apiwrap.Response[any] {
	documentList, err := h.serv.AdminFindAllRoot(ctx)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentDomainListToRootVOList(documentList), "获取根文档列表成功")
}

// 获取一个文档的所有子文档(包含非直接子文档)
func (h *DocumentHandler) AdminGetAllParentDoc(ctx *gin.Context) *apiwrap.Response[any] {
	documentID := ctx.Query("document_id")
	documentList, err := h.serv.AdminFindAllParent(ctx, apiwrap.ConvertBsonID(documentID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentDomainListToVOList(documentList), "获取父文档列表成功")
}

// 将文档请求转换为domain
func (h *DocumentHandler) DocumentRequestToDomain(req DocumentRequest) *domain.Document {
	return &domain.Document{
		Title:        req.Title,
		Content:      req.Content,
		DocumentType: req.DocumentType,
		ParentID:     apiwrap.ConvertBsonID(req.ParentID),
		DocumentID:   apiwrap.ConvertBsonID(req.DocumentID),
	}
}

// 将根文档请求转换为domain
func (h *DocumentHandler) DocumentRequestToRootDomain(req DocumentRootRequest) *domain.DocumentRoot {
	return &domain.DocumentRoot{
		Title:        req.Title,
		Alias:        req.Alias,
		Description:  req.Description,
		Thumbnail:    req.Thumbnail,
		IsPublic:     req.IsPublic,
		DocumentType: req.DocumentType,
	}
}

// 将文档domain转换为目录树vo
func (h *DocumentHandler) DocumentDomainToTreeVO(doc *domain.Document) *DocumentTreeVO {
	return &DocumentTreeVO{
		ID:           doc.ID.Hex(),
		CreatedAt:    doc.CreatedAt.String(),
		UpdatedAt:    doc.UpdatedAt.String(),
		Title:        doc.Title,
		DocumentType: doc.DocumentType,
		ParentID:     apiwrap.BsonID(doc.ParentID.Hex()),
		DocumentID:   apiwrap.BsonID(doc.DocumentID.Hex()),
	}
}

// 将文档domain列表转换为目录树
func (h *DocumentHandler) DocumentDomainListToTreeVOList(docList []*domain.Document) []*DocumentTreeVO {
	return lo.Map(docList, func(doc *domain.Document, _ int) *DocumentTreeVO {
		return h.DocumentDomainToTreeVO(doc)
	})
}

// 将根文档domain转换为vo
func (h *DocumentHandler) DocumentRootDomainToVO(doc *domain.DocumentRoot) *DocumentRootVO {
	return &DocumentRootVO{
		ID:           doc.ID.Hex(),
		CreatedAt:    doc.CreatedAt.String(),
		UpdatedAt:    doc.UpdatedAt.String(),
		Title:        doc.Title,
		Alias:        doc.Alias,
		Description:  doc.Description,
		Thumbnail:    doc.Thumbnail,
		IsPublic:     doc.IsPublic,
		DocumentType: doc.DocumentType,
	}
}

// 将根文档domain列表转换为vo
func (h *DocumentHandler) DocumentDomainListToRootVOList(docList []*domain.DocumentRoot) []*DocumentRootVO {
	return lo.Map(docList, func(doc *domain.DocumentRoot, _ int) *DocumentRootVO {
		return h.DocumentRootDomainToVO(doc)
	})
}

// 将文档domain转换为vo
func (h *DocumentHandler) DocumentDomainToVO(doc *domain.Document) *DocumentVO {
	return &DocumentVO{
		ID:           doc.ID.Hex(),
		CreatedAt:    doc.CreatedAt.String(),
		UpdatedAt:    doc.UpdatedAt.String(),
		Title:        doc.Title,
		Content:      doc.Content,
		DocumentType: doc.DocumentType,
		ParentID:     apiwrap.BsonID(doc.ParentID.Hex()),
		DocumentID:   apiwrap.BsonID(doc.DocumentID.Hex()),
	}
}

// 将文档domain列表转换为vo
func (h *DocumentHandler) DocumentDomainListToVOList(docList []*domain.Document) []*DocumentVO {
	return lo.Map(docList, func(doc *domain.Document, _ int) *DocumentVO {
		return h.DocumentDomainToVO(doc)
	})
}
