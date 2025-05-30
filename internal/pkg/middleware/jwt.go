package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/codepzj/stellux/server/global"
	"github.com/codepzj/stellux/server/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 过滤非管理员接口
		if !strings.HasPrefix(ctx.Request.RequestURI, "/admin-api") {
			ctx.Next()
			return
		}

		access_token := ctx.Request.Header.Get("Authorization")
		slog.Debug("用户携带的token", "access_token", access_token)
		// 若非GET请求的token为空
		if access_token == "" || !strings.HasPrefix(access_token, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusOK, global.AccessTokenNotFound)
			return
		}
		claims, err := utils.ParseToken(strings.TrimPrefix(access_token, "Bearer "))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, global.AccessTokenExpired)
			return
		}
		slog.Debug("解析后的用户ID", "userId", claims.ID)
		ctx.Set("userId", claims.ID)
		ctx.Next()
	}
}
