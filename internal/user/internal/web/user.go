package web

import (
	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"
	"github.com/codepzj/Stellux-Server/internal/pkg/middleware"
	"github.com/codepzj/Stellux-Server/internal/pkg/utils"
	"github.com/codepzj/Stellux-Server/internal/user/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/user/internal/service"
	"github.com/codepzj/gokit/str"
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
		userGroup.POST("/login", apiwrap.WrapWithJson(h.Login))
		userGroup.GET("/mockCreate", apiwrap.Wrap(h.MockCreateUser))
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

func (h *UserHandler) MockCreateUser(c *gin.Context) (*apiwrap.Response[any], error) {
	err := h.serv.AdminCreate(c, &domain.User{
		Username: "admin",
		Password: "admin",
		Nickname: "admin",
		RoleId:   1,
		Avatar:   "https://avatars.githubusercontent.com/u/22550712?v=4",
		Email:    "admin@admin.com",
	})
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	return apiwrap.SuccessWithMsg("创建用户成功"), nil
}

func (h *UserHandler) Login(c *gin.Context, userRequest LoginRequest) (*apiwrap.Response[any], error) {
	user, err := h.serv.VerifyUser(c, &domain.User{
		Username: userRequest.Username,
		Password: userRequest.Password,
	})
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	loginVO := LoginVO{
		AccessToken: accessToken,
	}
	return apiwrap.SuccessWithDetail[any](loginVO, "登录成功"), nil
}

func (h *UserHandler) AdminCreateUser(c *gin.Context, createUserRequest CreateUserRequest) (*apiwrap.Response[any], error) {
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
		return nil, apiwrap.NewInternalError(err.Error())
	}
	return apiwrap.SuccessWithMsg("创建用户成功"), nil
}

func (h *UserHandler) AdminUpdateUser(c *gin.Context, updateUserRequest UpdateUserRequest) (*apiwrap.Response[any], error) {
	user := domain.User{
		ID:       updateUserRequest.ID,
		Nickname: updateUserRequest.Nickname,
		Avatar:   updateUserRequest.Avatar,
		Email:    updateUserRequest.Email,
	}
	err := h.serv.AdminUpdate(c, &user)
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	return apiwrap.SuccessWithMsg("更新用户成功"), nil
}

func (h *UserHandler) AdminUpdatePassword(c *gin.Context, updatePasswordRequest UpdatePasswordRequest) (*apiwrap.Response[any], error) {
	err := h.serv.AdminUpdatePassword(c, updatePasswordRequest.ID, updatePasswordRequest.OldPassword, updatePasswordRequest.NewPassword)
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	return apiwrap.SuccessWithMsg("更新密码成功"), nil
}

func (h *UserHandler) AdminDeleteUser(c *gin.Context) (*apiwrap.Response[any], error) {
	id, err := str.ToUint(c.Param("id"))
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	err = h.serv.AdminDelete(c, id)
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	return apiwrap.SuccessWithMsg("删除用户成功"), nil
}

func (h *UserHandler) AdminGetUserList(c *gin.Context, page apiwrap.Page) (*apiwrap.Response[any], error) {
	users, count, err := h.serv.GetUserList(c, &domain.Page{
		PageNo:   page.PageNo,
		PageSize: page.PageSize,
	})
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	return apiwrap.SuccessWithDetail[any](apiwrap.ToPageVO(page.PageNo, page.PageSize, count, h.UserDomainToVOList(users)), "获取用户列表成功"), nil
}

func (h *UserHandler) AdminGetUserInfo(c *gin.Context) (*apiwrap.Response[any], error) {
	id := c.GetUint("userId")
	user, err := h.serv.GetUserInfo(c, id)
	if err != nil {
		return nil, apiwrap.NewInternalError(err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.UserDomainToVO(user), "获取用户信息成功"), nil
}
