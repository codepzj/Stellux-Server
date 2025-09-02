package web

import (
	"github.com/codepzj/Stellux-Server/internal/friend/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/friend/internal/service"
	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
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
func (h *FriendHandler) FindFriendList(c *gin.Context) (*apiwrap.Response[any], error) {
	friends, err := h.serv.FindFriendList(c)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessWithDetail[any](FriendDomainToShowVOList(friends), "获取友链列表成功"), nil
}

// CreateFriend 创建友链
func (h *FriendHandler) CreateFriend(c *gin.Context, friend *FriendRequest) (*apiwrap.Response[any], error) {
	err := h.serv.CreateFriend(c, &domain.Friend{
		Name:        friend.Name,
		Description: friend.Description,
		SiteUrl:     friend.SiteUrl,
		AvatarUrl:   friend.AvatarUrl,
		WebsiteType: friend.WebsiteType,
	})

	if err != nil {
		return nil, err
	}
	return apiwrap.Success(), nil
}

// FindAllFriends 获取所有友链
func (h *FriendHandler) FindAllFriends(c *gin.Context) (*apiwrap.Response[any], error) {
	friends, err := h.serv.FindAllFriends(c)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessWithDetail[any](FriendDomainToVOList(friends), "获取朋友列表成功"), nil
}

// UpdateFriend 更新友链
func (h *FriendHandler) UpdateFriend(c *gin.Context, friend *FriendUpdateRequest) (*apiwrap.Response[any], error) {
	objId, err := bson.ObjectIDFromHex(friend.ID)
	if err != nil {
		return nil, apiwrap.NewBadRequest("id格式错误")
	}
	err = h.serv.UpdateFriend(c, objId, &domain.Friend{
		Name:        friend.Name,
		Description: friend.Description,
		SiteUrl:     friend.SiteUrl,
		AvatarUrl:   friend.AvatarUrl,
		WebsiteType: friend.WebsiteType,
		IsActive:    friend.IsActive,
	})
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessWithMsg("更新好友成功"), nil
}

// DeleteFriend 删除友链
func (h *FriendHandler) DeleteFriend(c *gin.Context) (*apiwrap.Response[any], error) {
	objId, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return nil, apiwrap.NewBadRequest("id格式错误")
	}
	err = h.serv.DeleteFriend(c, objId)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessWithMsg("删除好友成功"), nil
}
