package web

import (
	"github.com/codepzj/stellux/server/internal/label/internal/domain"
	"github.com/codepzj/stellux/server/internal/label/internal/service"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func NewLabelHandler(serv service.ILabelService) *LabelHandler {
	return &LabelHandler{
		serv: serv,
	}
}

type LabelHandler struct {
	serv service.ILabelService
}

func (h *LabelHandler) RegisterGinRoutes(engine *gin.Engine) {
	adminGroup := engine.Group("/admin-api/label")
	{
		adminGroup.POST("/create", apiwrap.WrapWithBody(h.AdminCreate))
		adminGroup.PUT("/edit", apiwrap.WrapWithBody(h.AdminUpdate))
		adminGroup.DELETE("/delete/:id", apiwrap.Wrap(h.AdminDelete))
	}
	labelGroup := engine.Group("/label")
	{
		labelGroup.GET("/:id", apiwrap.Wrap(h.GetByID))
		labelGroup.GET("/list", apiwrap.WrapWithBody(h.QueryLabelList))
		labelGroup.GET("/all", apiwrap.Wrap(h.QueryAllByType))
	}
}

func (h *LabelHandler) AdminCreate(c *gin.Context, label *LabelRequest) *apiwrap.Response[any] {
	err := h.serv.Create(c, &domain.Label{
		LabelType: label.LabelType,
		Name:      label.Name,
	})
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("标签创建成功")
}

func (h *LabelHandler) AdminUpdate(c *gin.Context, label *LabelRequest) *apiwrap.Response[any] {
	err := h.serv.Update(c, label.ID, h.LabelDTOToDomain(label))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("标签更新成功")
}

func (h *LabelHandler) AdminDelete(c *gin.Context) *apiwrap.Response[any] {
	id := c.Param("id")
	err := h.serv.Delete(c, id)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("标签删除成功")
}

func (h *LabelHandler) GetByID(c *gin.Context) *apiwrap.Response[any] {
	id := c.Param("id")
	label, err := h.serv.Get(c, id)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.LabelDomainToVO(label), "标签获取成功")
}

func (h *LabelHandler) QueryLabelList(c *gin.Context, page *Page) *apiwrap.Response[any] {
	labels, count, err := h.serv.QueryLabelList(c, page.LabelType, page.PageNo, page.PageSize)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	labelVOList := h.DomainToVOList(labels)
	pageVO := apiwrap.ToPageVO(page.PageNo, page.PageSize, count, labelVOList)
	return apiwrap.SuccessWithDetail[any](pageVO, "标签列表获取成功")
}

func (h *LabelHandler) QueryAllByType(c *gin.Context) *apiwrap.Response[any] {
	labelType := c.Query("label_type")
	if labelType == "" {
		return apiwrap.FailWithMsg(apiwrap.RuquestBadRequest, "标签类型不能为空")
	}
	labels, err := h.serv.QueryAllByType(c, labelType)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DomainToVOList(labels), "标签列表获取成功")
}

func (h *LabelHandler) LabelDTOToDomain(label *LabelRequest) *domain.Label {
	return &domain.Label{
		ID:        apiwrap.ConvertBsonID(label.ID),
		LabelType: label.LabelType,
		Name:      label.Name,
	}
}

func (h *LabelHandler) LabelDomainToVO(label *domain.Label) *LabelVO {
	return &LabelVO{
		ID:        label.ID.Hex(),
		LabelType: label.LabelType,
		Name:      label.Name,
	}
}

func (h *LabelHandler) DomainToVOList(labels []*domain.Label) []*LabelVO {
	return lo.Map(labels, func(label *domain.Label, _ int) *LabelVO {
		return h.LabelDomainToVO(label)
	})
}
