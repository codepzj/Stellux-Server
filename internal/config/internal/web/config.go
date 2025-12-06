package web

import (
	"github.com/codepzj/Stellux-Server/internal/config/internal/service"
	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewConfigHandler(serv service.IConfigService) *ConfigHandler {
	return &ConfigHandler{
		serv: serv,
	}
}

type ConfigHandler struct {
	serv service.IConfigService
}

// RegisterGinRoutes 注册路由
func (h *ConfigHandler) RegisterGinRoutes(engine *gin.Engine) {
	// 管理API
	adminGroup := engine.Group("/admin-api/config")
	{
		adminGroup.POST("create", apiwrap.WrapWithJson(h.AdminCreateConfig))
		adminGroup.PUT("update", apiwrap.WrapWithJson(h.AdminUpdateConfig))
		adminGroup.DELETE(":id", apiwrap.WrapWithUri(h.AdminDeleteConfig))
		adminGroup.GET("list", apiwrap.Wrap(h.AdminListConfigs))
		adminGroup.GET(":id", apiwrap.WrapWithUri(h.AdminGetConfigByID))
	}

	// 公开API
	configGroup := engine.Group("/config")
	{
		configGroup.GET("/:type", apiwrap.WrapWithUri(h.GetConfig))
	}
}

// AdminCreateConfig 管理员创建网站配置
func (h *ConfigHandler) AdminCreateConfig(c *gin.Context, req ConfigDto) (int, string, any) {
	config := h.ConfigDtoToDomain(req)
	err := h.serv.CreateConfig(c, config)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, "创建网站配置成功", nil
}

// AdminUpdateConfig 管理员更新网站配置
func (h *ConfigHandler) AdminUpdateConfig(c *gin.Context, req ConfigUpdateDto) (int, string, any) {
	config := h.ConfigUpdateDtoToDomain(req)
	err := h.serv.UpdateConfig(c, config)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, "更新网站配置成功", nil
}

// AdminDeleteConfig 管理员删除网站配置
func (h *ConfigHandler) AdminDeleteConfig(c *gin.Context, req ConfigIdRequest) (int, string, any) {
	objId, err := bson.ObjectIDFromHex(req.ID)
	if err != nil {
		return 400, "ID格式错误", nil
	}

	err = h.serv.DeleteConfig(c, objId)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, "删除网站配置成功", nil
}

// AdminListConfigs 管理员获取网站配置列表
func (h *ConfigHandler) AdminListConfigs(c *gin.Context) (int, string, any) {
	configs, err := h.serv.ListConfigs(c)
	if err != nil {
		return 500, err.Error(), nil
	}

	summaryVOs := h.ConfigListToSummaryVOList(configs)
	return 200, "获取网站配置列表成功", summaryVOs
}

// AdminGetConfigByID 管理员根据ID获取网站配置
func (h *ConfigHandler) AdminGetConfigByID(c *gin.Context, req ConfigIdRequest) (int, string, any) {
	objId, err := bson.ObjectIDFromHex(req.ID)
	if err != nil {
		return 400, "ID格式错误", nil
	}

	config, err := h.serv.GetConfigByID(c, objId)
	if err != nil {
		return 404, "网站配置不存在", nil
	}

	vo := h.ConfigToVO(config)
	return 200, "获取网站配置成功", vo
}

// GetConfig 获取网站配置
func (h *ConfigHandler) GetConfig(c *gin.Context, req ConfigTypeRequest) (int, string, any) {
	config, err := h.serv.GetConfigByType(c, req.Type)
	if err != nil {
		return 404, "网站配置不存在", nil
	}

	vo := h.ConfigToVO(config)
	return 200, "获取网站配置成功", vo
}
