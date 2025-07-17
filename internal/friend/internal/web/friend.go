package web

import (
	"net/http"

	"github.com/codepzj/stellux/server/internal/friend/internal/domain"
	"github.com/codepzj/stellux/server/internal/friend/internal/service"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/gin-gonic/gin"
)

func NewFriendHandler(serv service.IFriendService) *FriendHandler {
	return &FriendHandler{
		serv: serv,
	}
}

type FriendHandler struct {
	serv service.IFriendService
}

func (h *FriendHandler) RegisterGinRoutes(engine *gin.Engine) {
	friendGroup := engine.Group("/friend")
	{
		friendGroup.GET("/list", apiwrap.Wrap(h.FindFriendList))
	}
	adminGroup := engine.Group("/admin-api")
	{
		adminGroup.POST("/friend/create", apiwrap.WrapWithJson(h.CreateFriend))
		adminGroup.GET("/friend/all", apiwrap.Wrap(h.FindAllFriends))
		adminGroup.PUT("/friend/update", apiwrap.WrapWithJson(h.UpdateFriend))
		adminGroup.DELETE("/friend/delete/:id", apiwrap.Wrap(h.DeleteFriend))
	}

}

// FindFriendList 获取友链列表
func (h *FriendHandler) FindFriendList(c *gin.Context) *apiwrap.Response[any] {
	friends, err := h.serv.FindFriendList(c)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](FriendDomainToShowVOList(friends), "获取友链列表成功")
}

// CreateFriend 创建友链
func (h *FriendHandler) CreateFriend(c *gin.Context, friend *FriendRequest) *apiwrap.Response[any] {
	err := h.serv.CreateFriend(c, &domain.Friend{
		Name:        friend.Name,
		Description: friend.Description,
		SiteUrl:     friend.SiteUrl,
		AvatarUrl:   friend.AvatarUrl,
		WebsiteType: friend.WebsiteType,
	})
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.Success()
}

// FindAllFriends 获取所有友链
func (h *FriendHandler) FindAllFriends(c *gin.Context) *apiwrap.Response[any] {
	friends, err := h.serv.FindAllFriends(c)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](FriendDomainToVOList(friends), "获取朋友列表成功")
}

// UpdateFriend 更新友链
func (h *FriendHandler) UpdateFriend(c *gin.Context, friend *FriendUpdateRequest) *apiwrap.Response[any] {
	err := h.serv.UpdateFriend(c, apiwrap.ConvertBsonID(friend.ID).ToObjectID(), &domain.Friend{
		Name:        friend.Name,
		Description: friend.Description,
		SiteUrl:     friend.SiteUrl,
		AvatarUrl:   friend.AvatarUrl,
		WebsiteType: friend.WebsiteType,
		IsActive:    friend.IsActive,
	})
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新好友成功")
}

// DeleteFriend 删除友链
func (h *FriendHandler) DeleteFriend(c *gin.Context) *apiwrap.Response[any] {
	err := h.serv.DeleteFriend(c, apiwrap.ConvertBsonID(c.Param("id")).ToObjectID())
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("删除好友成功")
}
