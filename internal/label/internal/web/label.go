package web

import (
	"github.com/codepzj/Stellux-Server/internal/label/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/label/internal/service"
	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	labelGroup := engine.Group("/label")
	{
		labelGroup.GET("/:id", apiwrap.Wrap(h.GetByID))                                  // 根据id获取标签
		labelGroup.GET("/list", apiwrap.WrapWithQuery(h.QueryLabelList))                 // 分页查询标签
		labelGroup.GET("/all", apiwrap.Wrap(h.QueryAllByType))                           // 获取所有标签
		labelGroup.GET("/categories/count", apiwrap.Wrap(h.QueryCategoryLabelWithCount)) // 获取分类标签及其文章数量
		labelGroup.GET("/tags/count", apiwrap.Wrap(h.QueryTagsLabelWithCount))           // 获取标签及其文章数量
	}
	adminGroup := engine.Group("/admin-api/label")
	{
		adminGroup.POST("/create", apiwrap.WrapWithJson(h.AdminCreate)) // 创建标签
		adminGroup.PUT("/edit", apiwrap.WrapWithJson(h.AdminUpdate))    // 更新标签
		adminGroup.DELETE("/delete/:id", apiwrap.Wrap(h.AdminDelete))   // 删除标签
	}
}

// AdminCreate 创建标签
func (h *LabelHandler) AdminCreate(c *gin.Context, label *LabelRequest) (*apiwrap.Response[any], error) {
	err := h.serv.CreateLabel(c, &domain.Label{
		LabelType: label.LabelType,
		Name:      label.Name,
	})
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	return apiwrap.SuccessWithMsg("标签创建成功"), nil
}

// AdminUpdate 更新标签
func (h *LabelHandler) AdminUpdate(c *gin.Context, label *LabelRequest) (*apiwrap.Response[any], error) {
	objId, err := bson.ObjectIDFromHex(label.ID)
	if err != nil {
		return nil, apiwrap.NewBadRequest("id格式错误")
	}
	err = h.serv.UpdateLabel(c, label.ID, &domain.Label{
		Id:        objId,
		LabelType: label.LabelType,
		Name:      label.Name,
	})
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	return apiwrap.SuccessWithMsg("标签更新成功"), nil
}

// AdminDelete 删除标签
func (h *LabelHandler) AdminDelete(c *gin.Context) (*apiwrap.Response[any], error) {
	id := c.Param("id")
	err := h.serv.DeleteLabel(c, id)
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	return apiwrap.SuccessWithMsg("标签删除成功"), nil
}

// GetByID 根据id获取标签
func (h *LabelHandler) GetByID(c *gin.Context) (*apiwrap.Response[any], error) {
	id := c.Param("id")
	label, err := h.serv.GetLabelById(c, id)
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.LabelDomainToVO(label), "标签获取成功"), nil
}

// QueryLabelList 分页查询标签
func (h *LabelHandler) QueryLabelList(c *gin.Context, page *Page) (*apiwrap.Response[any], error) {
	labels, count, err := h.serv.QueryLabelList(c, page.LabelType, page.PageNo, page.PageSize)
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	labelVOList := h.DomainToVOList(labels)
	pageVO := apiwrap.ToPageVO(page.PageNo, page.PageSize, count, labelVOList)
	return apiwrap.SuccessWithDetail[any](pageVO, "标签列表获取成功"), nil
}

// QueryAllByType 获取所有标签
func (h *LabelHandler) QueryAllByType(c *gin.Context) (*apiwrap.Response[any], error) {
	labelType := c.Query("label_type")
	if labelType == "" {
		return nil, apiwrap.NewBadRequest("标签类型不能为空")
	}
	labels, err := h.serv.GetAllLabelsByType(c, labelType)
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DomainToVOList(labels), "标签列表获取成功"), nil
}

// QueryCategoryLabelWithCount 获取分类标签及其文章数量
func (h *LabelHandler) QueryCategoryLabelWithCount(c *gin.Context) (*apiwrap.Response[any], error) {

	labels, err := h.serv.GetAllLabelsWithCount(c)
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	labelVOList := h.LabelWithCountDomainToVOList(labels)
	return apiwrap.SuccessWithDetail[any](labelVOList, "分类标签列表获取成功"), nil
}

// QueryTagsLabelWithCount 获取标签及其文章数量
func (h *LabelHandler) QueryTagsLabelWithCount(c *gin.Context) (*apiwrap.Response[any], error) {
	labels, err := h.serv.GetAllTagsLabelWithCount(c)
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	labelVOList := h.LabelWithCountDomainToVOList(labels)
	return apiwrap.SuccessWithDetail[any](labelVOList, "标签列表获取成功"), nil
}

func (h *LabelHandler) LabelDTOToDomain(label *LabelRequest) *domain.Label {
	objId, _ := bson.ObjectIDFromHex(label.ID)
	return &domain.Label{
		Id:        objId,
		LabelType: label.LabelType,
		Name:      label.Name,
	}
}

func (h *LabelHandler) LabelDomainToVO(label *domain.Label) *LabelVO {
	return &LabelVO{
		ID:        label.Id.Hex(),
		LabelType: label.LabelType,
		Name:      label.Name,
	}
}

func (h *LabelHandler) DomainToVOList(labels []*domain.Label) []*LabelVO {
	return lo.Map(labels, func(label *domain.Label, _ int) *LabelVO {
		return h.LabelDomainToVO(label)
	})
}

func (h *LabelHandler) LabelWithCountDomainToVOList(labels []*domain.LabelPostCount) []*LabelWithCountVO {
	return lo.Map(labels, func(label *domain.LabelPostCount, _ int) *LabelWithCountVO {
		return &LabelWithCountVO{
			ID:        label.Id.Hex(),
			LabelType: label.LabelType,
			Name:      label.Name,
			Count:     label.Count,
		}
	})
}
