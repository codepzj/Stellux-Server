package web

import (
	"fmt"
	"time"

	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/codepzj/stellux/server/internal/pkg/utils"
	"github.com/codepzj/stellux/server/internal/user/internal/domain"
	"github.com/codepzj/stellux/server/internal/user/internal/service"
	"github.com/gin-gonic/gin"
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
		userGroup.GET("/refresh", apiwrap.Wrap(h.RefreshToken))
		userGroup.POST("/login", apiwrap.WrapWithBody(h.Login))
	}
	adminGroup := engine.Group("/admin-api/user")
	{
		adminGroup.POST("/create", apiwrap.WrapWithBody(h.AdminCreateUser))
		adminGroup.PUT("/update", apiwrap.WrapWithBody(h.AdminUpdateUser))
		adminGroup.PUT("/update-password", apiwrap.WrapWithBody(h.AdminUpdatePassword))
		adminGroup.DELETE("/delete/:id", apiwrap.Wrap(h.AdminDeleteUser))
		adminGroup.GET("/list", apiwrap.WrapWithBody(h.AdminGetUserList))
		adminGroup.GET("/info", apiwrap.Wrap(h.AdminGetUserInfo))
	}
}

func (h *UserHandler) Login(c *gin.Context, userRequest LoginRequest) *apiwrap.Response[any] {
	user := domain.User{
		Username: userRequest.Username,
		Password: userRequest.Password,
	}
	exist, id := h.serv.CheckUserExist(c, &user)
	if !exist {
		return apiwrap.FailWithMsg(apiwrap.RuquestBadRequest, "用户名或密码错误")
	}
	accessToken, err := utils.GenerateAccessToken(id)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	refreshToken, err := utils.GenerateRefreshToken(id)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	loginVO := LoginVO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return apiwrap.SuccessWithDetail[any](loginVO, "登录成功")
}

func (h *UserHandler) RefreshToken(c *gin.Context) *apiwrap.Response[any] {
	// 校验refresh_token是否有效
	refreshToken := c.Query("refresh_token")
	fmt.Println(refreshToken)
	if refreshToken == "" {
		return apiwrap.FailWithMsg(apiwrap.RuquestBadRequest, "refresh_token不能为空")
	}
	claims, err := utils.ParseToken(refreshToken)
	if err != nil {
		fmt.Println(err.Error())
		return apiwrap.FailWithMsg(apiwrap.RequestRefreshTokenExpired, "refresh_token已过期,请重新登录")
	}
	accessToken, err := utils.GenerateAccessToken(claims.ID)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	// 若refresh_token临近过期，则生成新的refresh_token
	if claims.ExpiresAt.Before(time.Now().Add(time.Hour * 24 * 7)) {
		var err error
		refreshToken, err = utils.GenerateRefreshToken(claims.ID)
		if err != nil {
			return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
		}
	}
	return apiwrap.SuccessWithDetail[any](LoginVO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, "刷新token成功")
}

func (h *UserHandler) AdminCreateUser(c *gin.Context, createUserRequest CreateUserRequest) *apiwrap.Response[any] {
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
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("创建用户成功")
}

func (h *UserHandler) AdminUpdateUser(c *gin.Context, updateUserRequest UpdateUserRequest) *apiwrap.Response[any] {
	user := domain.User{
		ID:       apiwrap.ConvertBsonID(updateUserRequest.ID),
		Nickname: updateUserRequest.Nickname,
		Avatar:   updateUserRequest.Avatar,
		Email:    updateUserRequest.Email,
	}
	err := h.serv.AdminUpdate(c, &user)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新用户成功")
}

func (h *UserHandler) AdminUpdatePassword(c *gin.Context, updatePasswordRequest UpdatePasswordRequest) *apiwrap.Response[any] {
	err := h.serv.AdminUpdatePassword(c, updatePasswordRequest.ID, updatePasswordRequest.OldPassword, updatePasswordRequest.NewPassword)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新密码成功")
}

func (h *UserHandler) AdminDeleteUser(c *gin.Context) *apiwrap.Response[any] {
	id := c.Param("id")
	err := h.serv.AdminDelete(c, id)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("删除用户成功")
}

func (h *UserHandler) AdminGetUserList(c *gin.Context, page apiwrap.Page) *apiwrap.Response[any] {
	users, count, err := h.serv.GetUserList(c, &domain.Page{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	})
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](apiwrap.ToPageVO(page.PageNo, page.PageSize, count, h.UserDomainToVOList(users)), "获取用户列表成功")
}

func (h *UserHandler) AdminGetUserInfo(c *gin.Context) *apiwrap.Response[any] {
	id := c.GetString("userId")
	user, err := h.serv.GetUserInfo(c, id)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.UserDomainToVO(user), "获取用户信息成功")
}
