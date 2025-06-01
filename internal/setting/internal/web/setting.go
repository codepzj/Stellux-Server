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
		settingGroup.GET("/basic", apiwrap.Wrap(h.GetBasicSetting))
		settingGroup.GET("/seo", apiwrap.Wrap(h.GetSeoSetting))
		settingGroup.GET("/blog", apiwrap.Wrap(h.GetBlogSetting))
		settingGroup.GET("/about", apiwrap.Wrap(h.GetAboutSetting))
	}
	adminSettingGroup := engine.Group("/admin-api/setting")
	{
		adminSettingGroup.POST("/upsert/basic", apiwrap.WrapWithBody(h.AdminUpsertBasicSetting))
		adminSettingGroup.POST("/upsert/seo", apiwrap.WrapWithBody(h.AdminUpsertSeoSetting))
		adminSettingGroup.POST("/upsert/blog", apiwrap.WrapWithBody(h.AdminUpsertBlogSetting))
		adminSettingGroup.POST("/upsert/about", apiwrap.WrapWithBody(h.AdminUpsertAboutSetting))
	}
}

func (h *SettingHandler) AdminUpsertBasicSetting(c *gin.Context, req BasicSettingRequest) *apiwrap.Response[any] {
	err := h.serv.AdminUpsertSetting(c, h.BasicSettingRequestToDomain(req))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新基础设置成功")
}

func (h *SettingHandler) AdminUpsertSeoSetting(c *gin.Context, req SeoSettingRequest) *apiwrap.Response[any] {
	err := h.serv.AdminUpsertSetting(c, h.SeoSettingRequestToDomain(req))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新SEO设置成功")
}

func (h *SettingHandler) AdminUpsertBlogSetting(c *gin.Context, req BlogSettingRequest) *apiwrap.Response[any] {
	err := h.serv.AdminUpsertSetting(c, h.BlogSettingRequestToDomain(req))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新博客设置成功")
}

func (h *SettingHandler) AdminUpsertAboutSetting(c *gin.Context, req AboutSettingRequest) *apiwrap.Response[any] {

	err := h.serv.AdminUpsertSetting(c, h.AboutSettingRequestToDomain(req))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新关于设置成功")
}

func (h *SettingHandler) GetBlogSetting(c *gin.Context) *apiwrap.Response[any] {
	setting, err := h.serv.GetSetting(c, "blog_setting")
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	if setting.Value == nil {
		return apiwrap.SuccessWithDetail[any](BlogSettingVO{}, "获取博客设置成功")
	}
	return apiwrap.SuccessWithDetail(setting.Value, "获取博客设置成功")
}

func (h *SettingHandler) GetSeoSetting(c *gin.Context) *apiwrap.Response[any] {
	setting, err := h.serv.GetSetting(c, "seo_setting")
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	if setting.Value == nil {
		return apiwrap.SuccessWithDetail[any](SeoSettingVO{}, "获取SEO设置成功")
	}
	return apiwrap.SuccessWithDetail(setting.Value, "获取SEO设置成功")
}

func (h *SettingHandler) GetBasicSetting(c *gin.Context) *apiwrap.Response[any] {
	setting, err := h.serv.GetSetting(c, "basic_setting")
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	if setting.Value == nil {
		return apiwrap.SuccessWithDetail[any](BasicSettingVO{}, "获取基础设置成功")
	}
	return apiwrap.SuccessWithDetail(setting.Value, "获取基础设置成功")
}

func (h *SettingHandler) GetAboutSetting(c *gin.Context) *apiwrap.Response[any] {
	setting, err := h.serv.GetSetting(c, "about_setting")
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	if setting.Value == nil {
		return apiwrap.SuccessWithDetail[any](AboutSettingVO{}, "获取关于设置成功")
	}
	return apiwrap.SuccessWithDetail(setting.Value, "获取关于设置成功")
}

func (h *SettingHandler) BasicSettingRequestToDomain(req BasicSettingRequest) domain.Setting {
	return domain.Setting{
		Key: "basic_setting",
		Value: map[string]any{
			"site_title":    req.SiteTitle,
			"site_subtitle": req.SiteSubtitle,
			"site_favicon":  req.SiteFavicon,
		},
	}
}

func (h *SettingHandler) SeoSettingRequestToDomain(req SeoSettingRequest) domain.Setting {
	return domain.Setting{
		Key: "seo_setting",
		Value: map[string]any{
			"site_author":      req.SiteAuthor,
			"site_url":         req.SiteUrl,
			"site_description": req.SiteDescription,
			"site_keywords":    req.SiteKeywords,
			"robots":           req.Robots,
			"og_image":         req.OgImage,
			"og_type":          req.OgType,
			"twitter_card":     req.TwitterCard,
			"twitter_site":     req.TwitterSite,
		},
	}
}	

func (h *SettingHandler) BlogSettingRequestToDomain(req BlogSettingRequest) domain.Setting {
	return domain.Setting{
		Key: "blog_setting",
		Value: map[string]any{
			"blog_avatar":    req.BlogAvatar,
			"blog_title":     req.BlogTitle,
			"blog_subtitle":  req.BlogSubtitle,
		},
	}
}

func (h *SettingHandler) AboutSettingRequestToDomain(req AboutSettingRequest) domain.Setting {
	return domain.Setting{
		Key: "about_setting",
		Value: map[string]any{
			"author":        req.Author,
			"avatar_url":    req.AvatarUrl,
			"left_tags":     req.LeftTags,
			"right_tags":    req.RightTags,
			"know_me":       req.KnowMe,
			"github_username": req.GithubUsername,
		},
	}
}