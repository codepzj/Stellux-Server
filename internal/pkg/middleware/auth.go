package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/codepzj/stellux/server/global"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 过滤非管理员接口
		if !strings.HasPrefix(ctx.Request.RequestURI, "/admin-api") {
			ctx.Next()
			return
		}
		userId := ctx.GetString("userId")
		requestURI := ctx.Request.RequestURI
		method := ctx.Request.Method

		ok, err := global.Enforcer.Enforce(userId, requestURI, method)
		if err != nil {
			slog.Error("权限校验失败", "error", err)
			ctx.AbortWithStatusJSON(http.StatusOK, global.LoadPermissionFailed)
			return
		}
		if ok {
			slog.Debug("允许访问", "userId", userId, "requestURI", requestURI, "method", method)
			ctx.Next()
		} else {
			slog.Debug("禁止访问", "userId", userId, "requestURI", requestURI, "method", method)
			ctx.AbortWithStatusJSON(http.StatusOK, global.PermissionDenied)
		}
	}
}
