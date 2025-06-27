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
		friendGroup.POST("/create", apiwrap.WrapWithBody(h.CreateFriend))
		friendGroup.GET("/all", apiwrap.Wrap(h.FindAllFriends))
		friendGroup.PUT("/update", apiwrap.WrapWithBody(h.UpdateFriend))
		friendGroup.DELETE("/delete/:id", apiwrap.Wrap(h.DeleteFriend))
	}

}

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

func (h *FriendHandler) FindAllFriends(c *gin.Context) *apiwrap.Response[any] {
	friends, err := h.serv.FindAllFriends(c)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](FriendDomainToVOList(friends), "获取朋友列表成功")
}

func (h *FriendHandler) UpdateFriend(c *gin.Context, friend *FriendUpdateRequest) *apiwrap.Response[any] {
	err := h.serv.UpdateFriend(c, apiwrap.ConvertBsonID(friend.ID), &domain.Friend{
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

func (h *FriendHandler) DeleteFriend(c *gin.Context) *apiwrap.Response[any] {
	err := h.serv.DeleteFriend(c, apiwrap.ConvertBsonID(c.Param("id")))
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("删除好友成功")
}