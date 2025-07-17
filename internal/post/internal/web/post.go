package web

import (
	"encoding/json"

	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/codepzj/stellux/server/internal/post/internal/service"
	"github.com/codepzj/stellux/server/internal/setting"
	"github.com/gin-gonic/gin"
)

func NewPostHandler(serv service.IPostService, settingService setting.Service) *PostHandler {
	return &PostHandler{
		serv:           serv,
		settingService: settingService,
	}
}

type PostHandler struct {
	serv           service.IPostService
	settingService setting.Service
}

func (h *PostHandler) RegisterGinRoutes(engine *gin.Engine) {
	adminGroup := engine.Group("/admin-api/post")
	{
		adminGroup.GET("draft/list", apiwrap.WrapWithJson(h.AdminGetDraftDetailPostList))
		adminGroup.GET("bin/list", apiwrap.WrapWithJson(h.AdminGetBinDetailPostList))
		adminGroup.POST("create", apiwrap.WrapWithJson(h.AdminCreatePost))
		adminGroup.PUT("update", apiwrap.WrapWithJson(h.AdminUpdatePost))
		adminGroup.PUT("update/publish-status", apiwrap.WrapWithJson(h.AdminUpdatePostPublishStatus))
		adminGroup.PUT("restore/:id", apiwrap.WrapWithUri(h.AdminRestorePost))
		adminGroup.PUT("restore/batch", apiwrap.WrapWithJson(h.AdminRestorePostBatch))
		adminGroup.DELETE("soft-delete/:id", apiwrap.WrapWithUri(h.AdminSoftDeletePost))
		adminGroup.DELETE("soft-delete/batch", apiwrap.WrapWithJson(h.AdminSoftDeletePostBatch))
		adminGroup.DELETE("delete/:id", apiwrap.WrapWithUri(h.AdminDeletePost))
		adminGroup.DELETE("delete/batch", apiwrap.WrapWithJson(h.AdminDeletePostBatch))
	}
	postGroup := engine.Group("/post")
	{
		postGroup.GET("/list", apiwrap.WrapWithQuery(h.GetPublishPostList)) // 获取发布文章列表
		postGroup.GET("/:id", apiwrap.WrapWithUri(h.GetPostById))           // 获取文章详情
		postGroup.GET("/search", apiwrap.Wrap(h.GetPostByKeyWord))          // 搜索文章
		postGroup.GET("/sitemap", apiwrap.Wrap(h.GetSiteMap))               // 获取站点地图
		postGroup.GET("/all", apiwrap.Wrap(h.GetAllPublishPost))            // 获取所有发布文章
	}
}

func (h *PostHandler) AdminCreatePost(c *gin.Context, postReq PostDto) *apiwrap.Response[any] {
	post := h.PostDTOToDomain(postReq)
	err := h.serv.AdminCreatePost(c, post)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("创建文章成功")
}

func (h *PostHandler) AdminUpdatePost(c *gin.Context, postUpdateReq PostUpdateDto) *apiwrap.Response[any] {
	postUpdate := h.PostUpdateDTOToDomain(postUpdateReq)
	err := h.serv.AdminUpdatePost(c, postUpdate)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新文章成功")
}

func (h *PostHandler) AdminUpdatePostPublishStatus(c *gin.Context, postPublishStatusRequest PostPublishStatusRequest) *apiwrap.Response[any] {
	err := h.serv.AdminUpdatePostPublishStatus(c, apiwrap.ConvertBsonID(postPublishStatusRequest.ID).ToObjectID(), *postPublishStatusRequest.IsPublish)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新文章发布状态成功")
}

func (h *PostHandler) AdminRestorePost(c *gin.Context, postIDRequest PostIdRequest) *apiwrap.Response[any] {
	err := h.serv.AdminRestorePost(c, apiwrap.ConvertBsonID(postIDRequest.Id).ToObjectID())
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("恢复文章成功")
}

func (h *PostHandler) AdminRestorePostBatch(c *gin.Context, postIDListRequest PostIDListRequest) *apiwrap.Response[any] {
	err := h.serv.AdminRestorePostBatch(c, apiwrap.ToObjectIDList(apiwrap.ConvertBsonIDList(postIDListRequest.IDList)))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("批量恢复文章成功")
}

func (h *PostHandler) AdminSoftDeletePost(c *gin.Context, postIDRequest PostIdRequest) *apiwrap.Response[any] {
	err := h.serv.AdminSoftDeletePost(c, apiwrap.ConvertBsonID(postIDRequest.Id).ToObjectID())
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("软删除文章成功")
}

func (h *PostHandler) AdminSoftDeletePostBatch(c *gin.Context, postIDListRequest PostIDListRequest) *apiwrap.Response[any] {
	err := h.serv.AdminSoftDeletePostBatch(c, apiwrap.ToObjectIDList(apiwrap.ConvertBsonIDList(postIDListRequest.IDList)))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("批量软删除文章成功")
}

func (h *PostHandler) AdminDeletePost(c *gin.Context, postIDRequest PostIdRequest) *apiwrap.Response[any] {
	err := h.serv.AdminDeletePost(c, apiwrap.ConvertBsonID(postIDRequest.Id).ToObjectID())
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("删除文章成功")
}

func (h *PostHandler) AdminDeletePostBatch(c *gin.Context, postIDListRequest PostIDListRequest) *apiwrap.Response[any] {
	err := h.serv.AdminDeletePostBatch(c, apiwrap.ToObjectIDList(apiwrap.ConvertBsonIDList(postIDListRequest.IDList)))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("批量删除文章成功")
}

// GetPublishPostList 获取发布文章列表
func (h *PostHandler) GetPublishPostList(c *gin.Context, pageReq apiwrap.Page) *apiwrap.Response[any] {
	postDetailList, total, err := h.serv.GetPostList(c, &pageReq, "publish")
	if err != nil {
		return apiwrap.FailWithMsg(500, err.Error())
	}
	postVos := h.PostDetailListToVOList(postDetailList)

	pageVo := apiwrap.ToPageVO(pageReq.PageNo, pageReq.PageSize, total, postVos)
	return apiwrap.SuccessWithDetail[any](pageVo, "获取文章列表成功")
}

// AdminGetDraftDetailPostList 获取草稿箱文章列表
func (h *PostHandler) AdminGetDraftDetailPostList(c *gin.Context, pageReq apiwrap.Page) *apiwrap.Response[any] {
	postDetailList, total, err := h.serv.GetPostList(c, &pageReq, "draft")
	if err != nil {
		return apiwrap.FailWithMsg(500, err.Error())
	}
	postVos := h.PostDetailListToVOList(postDetailList)

	pageVo := apiwrap.ToPageVO(pageReq.PageNo, pageReq.PageSize, total, postVos)
	return apiwrap.SuccessWithDetail[any](pageVo, "获取草稿箱文章列表成功")
}

func (h *PostHandler) AdminGetBinDetailPostList(c *gin.Context, pageReq apiwrap.Page) *apiwrap.Response[any] {
	postDetailList, total, err := h.serv.GetPostList(c, &pageReq, "bin")
	if err != nil {
		return apiwrap.FailWithMsg(500, err.Error())
	}
	postVos := h.PostDetailListToVOList(postDetailList)

	pageVo := apiwrap.ToPageVO(pageReq.PageNo, pageReq.PageSize, total, postVos)
	return apiwrap.SuccessWithDetail[any](pageVo, "获取回收站文章列表成功")
}

// GetPostById 获取文章详情
func (h *PostHandler) GetPostById(c *gin.Context, postIDRequest PostIdRequest) *apiwrap.Response[any] {
	postDetail, err := h.serv.GetPostDetailById(c, apiwrap.ConvertBsonID(postIDRequest.Id).ToObjectID())
	if err != nil {
		return apiwrap.FailWithMsg(500, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.PostDetailToVO(postDetail), "获取文章详情成功")
}

func (h *PostHandler) GetPostByKeyWord(c *gin.Context) *apiwrap.Response[any] {
	keyword := c.Query("keyword")
	postList, err := h.serv.GetPostByKeyWord(c, keyword)
	if err != nil {
		return apiwrap.FailWithMsg(500, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.PostListToVOList(postList), "获取文章成功")
}

func (h *PostHandler) GetAllPublishPost(c *gin.Context) *apiwrap.Response[any] {
	postList, err := h.serv.GetAllPublishPost(c)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.PostListToVOList(postList), "获取所有发布文章成功")
}

func (h *PostHandler) GetSiteMap(c *gin.Context) *apiwrap.Response[any] {
	postList, err := h.serv.GetAllPublishPost(c)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	setting, err := h.settingService.GetSetting(c, "seo_setting")
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	var seoSetting SeoSettingVO
	bj, _ := json.Marshal(setting.Value)
	json.Unmarshal(bj, &seoSetting)
	siteMapVos := h.PostListToSiteMapVOList(postList, seoSetting.SiteUrl)
	return apiwrap.SuccessWithDetail[any](siteMapVos, "获取站点地图成功")
}
