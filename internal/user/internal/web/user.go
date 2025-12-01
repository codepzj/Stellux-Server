package web

import (
	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"
	"github.com/codepzj/Stellux-Server/internal/pkg/middleware"
	"github.com/codepzj/Stellux-Server/internal/pkg/utils"
	"github.com/codepzj/Stellux-Server/internal/user/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/user/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewUserHandler(serv service.IUserService) *UserHandler {
	return &UserHandler{
		serv: serv,
	}
}

type UserHandler struct {
	serv service.IUserService
}

func (h *UserHandler) RegisterGinRoutes(engine *gin.Engine) {
	userGroup := engine.Group("/user")
	{
		userGroup.POST("/login", apiwrap.WrapWithJson(h.Login))
	}
	adminGroup := engine.Group("/admin-api/user")
	{
		adminGroup.Use(middleware.JWT())
		adminGroup.POST("/create", apiwrap.WrapWithJson(h.AdminCreateUser))
		adminGroup.PUT("/update", apiwrap.WrapWithJson(h.AdminUpdateUser))
		adminGroup.PUT("/update-password", apiwrap.WrapWithJson(h.AdminUpdatePassword))
		adminGroup.DELETE("/delete/:id", apiwrap.Wrap(h.AdminDeleteUser))
		adminGroup.GET("/list", apiwrap.WrapWithQuery(h.AdminGetUserList))
		adminGroup.GET("/info", apiwrap.Wrap(h.AdminGetUserInfo))
	}
}

func (h *UserHandler) Login(c *gin.Context, userRequest LoginRequest) (int, string, any) {
	user := domain.User{
		Username: userRequest.Username,
		Password: userRequest.Password,
	}
	exist, id := h.serv.CheckUserExist(c, &user)
	if !exist {
		return 500, "用户名或密码错误", nil
	}
	accessToken, err := utils.GenerateAccessToken(id)
	if err != nil {
		return 500, err.Error(), nil
	}
	loginVO := LoginVO{
		AccessToken: accessToken,
	}
	return 200, "登录成功", loginVO
}

func (h *UserHandler) AdminCreateUser(c *gin.Context, createUserRequest CreateUserRequest) (int, string, any) {
	user := domain.User{
		Username: createUserRequest.Username,
		Password: createUserRequest.Password,
		Nickname: createUserRequest.Nickname,
		RoleId:   *createUserRequest.RoleId,
		Avatar:   createUserRequest.Avatar,
		Email:    createUserRequest.Email,
	}
	err := h.serv.AdminCreate(c, &user)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, "创建用户成功", nil
}

func (h *UserHandler) AdminUpdateUser(c *gin.Context, updateUserRequest UpdateUserRequest) (int, string, any) {
	objId, err := bson.ObjectIDFromHex(updateUserRequest.ID)
	if err != nil {
		return 500, err.Error(), nil
	}
	user := domain.User{
		ID:       objId,
		Nickname: updateUserRequest.Nickname,
		Avatar:   updateUserRequest.Avatar,
		Email:    updateUserRequest.Email,
	}
	err = h.serv.AdminUpdate(c, &user)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, "更新用户成功", nil
}

func (h *UserHandler) AdminUpdatePassword(c *gin.Context, updatePasswordRequest UpdatePasswordRequest) (int, string, any) {
	err := h.serv.AdminUpdatePassword(c, updatePasswordRequest.ID, updatePasswordRequest.OldPassword, updatePasswordRequest.NewPassword)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, "更新密码成功", nil
}

func (h *UserHandler) AdminDeleteUser(c *gin.Context) (int, string, any) {
	id := c.Param("id")
	err := h.serv.AdminDelete(c, id)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, "删除用户成功", nil
}

func (h *UserHandler) AdminGetUserList(c *gin.Context, page apiwrap.Page) (int, string, any) {
	users, count, err := h.serv.GetUserList(c, &apiwrap.Page{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	})
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, "获取用户列表成功", apiwrap.ToPageVO(page.PageNo, page.PageSize, count, h.UserDomainToVOList(users))}

func (h *UserHandler) AdminGetUserInfo(c *gin.Context) (int, string, any) {
	id := c.GetString("userId")
	user, err := h.serv.GetUserInfo(c, id)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, "获取用户信息成功", h.UserDomainToVO(user)
}
