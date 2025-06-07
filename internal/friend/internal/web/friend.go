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
	}

}

func (h *FriendHandler) CreateFriend(c *gin.Context, friend *FriendRequest) *apiwrap.Response[any] {
	err := h.serv.Create(c, &domain.Friend{
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
	friends, err := h.serv.FindAll(c)
	if err != nil {
		return apiwrap.FailWithMsg(http.StatusInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](FriendDomainToVOList(friends), "获取朋友列表成功")
}