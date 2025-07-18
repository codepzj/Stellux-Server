package web

import (
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/codepzj/stellux/server/internal/setting/internal/domain"
	"github.com/codepzj/stellux/server/internal/setting/internal/service"
	"github.com/gin-gonic/gin"
)

func NewSettingHandler(serv service.ISettingService) *SettingHandler {
	return &SettingHandler{
		serv: serv,
	}
}

type SettingHandler struct {
	serv service.ISettingService
}

func (h *SettingHandler) RegisterGinRoutes(engine *gin.Engine) {
	settingGroup := engine.Group("/setting")
	{
		settingGroup.GET("/site_config", apiwrap.Wrap(h.GetSiteConfigSetting))
	}
	adminSettingGroup := engine.Group("/admin-api/setting")
	{
		adminSettingGroup.POST("/upsert/site_config", apiwrap.WrapWithJson(h.AdminUpsertSiteConfigSetting))
	}
}

func (h *SettingHandler) AdminUpsertSiteConfigSetting(c *gin.Context, req SiteConfigRequest) *apiwrap.Response[any] {
	err := h.serv.AdminUpsertSetting(c, h.SiteConfigRequestToDomain(req))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新网站配置成功")
}

// GetSiteConfigSetting 获取网站配置
func (h *SettingHandler) GetSiteConfigSetting(c *gin.Context) *apiwrap.Response[any] {
	setting, err := h.serv.GetSetting(c, "site_config")

	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.SiteConfigDomainToVO(*setting), "获取网站配置成功")
}

// SiteConfigRequestToDomain 网站配置请求转换为领域对象
func (h *SettingHandler) SiteConfigRequestToDomain(req SiteConfigRequest) domain.SiteSetting {
	return domain.SiteSetting{
		Key: "site_config",
		Value: domain.SiteConfig{
			SiteTitle:       req.SiteTitle,
			SiteSubtitle:    req.SiteSubtitle,
			SiteUrl:         req.SiteUrl,
			SiteFavicon:     req.SiteFavicon,
			SiteAuthor:      req.SiteAuthor,
			SiteAnimateText: req.SiteAnimateText,
			SiteAvatar:      req.SiteAvatar,
			SiteKeywords:    req.SiteKeywords,
			SiteDescription: req.SiteDescription,
			SiteCopyright:   req.SiteCopyright,
			SiteIcp:         req.SiteIcp,
			SiteIcpLink:     req.SiteIcpLink,
			GithubUsername:  req.GithubUsername,
		},
	}
}

// SiteConfigRequestToVO 网站配置请求转换为VO
func (h *SettingHandler) SiteConfigDomainToVO(req domain.SiteSetting) SiteConfigSettingVO {
	return SiteConfigSettingVO{
		SiteTitle:       req.Value.SiteTitle,
		SiteSubtitle:    req.Value.SiteSubtitle,
		SiteUrl:         req.Value.SiteUrl,
		SiteFavicon:     req.Value.SiteFavicon,
		SiteAuthor:      req.Value.SiteAuthor,
		SiteAnimateText: req.Value.SiteAnimateText,
		SiteAvatar:      req.Value.SiteAvatar,
		SiteDescription: req.Value.SiteDescription,
		SiteKeywords:    req.Value.SiteKeywords,
		SiteCopyright:   req.Value.SiteCopyright,
		SiteIcp:         req.Value.SiteIcp,
		SiteIcpLink:     req.Value.SiteIcpLink,
		GithubUsername:  req.Value.GithubUsername,
	}
}
