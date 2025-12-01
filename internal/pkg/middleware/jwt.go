package middleware

import (
	"strings"

	"github.com/codepzj/Stellux-Server/internal/pkg/logger"
	"github.com/codepzj/Stellux-Server/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		access_token := ctx.Request.Header.Get("Authorization")
		logger.Debug("用户携带的token", logger.WithString("access_token", access_token))
		// 若非GET请求的token为空
		if access_token == "" || !strings.HasPrefix(access_token, "Bearer ") {
			ctx.AbortWithStatusJSON(200, gin.H{"code": 401, "error": "未携带access_token"})
			return
		}
		claims, err := utils.ParseToken(strings.TrimPrefix(access_token, "Bearer "))
		if err != nil {
			logger.Error("解析token失败", logger.WithError(err))
			ctx.AbortWithStatusJSON(200, gin.H{"code": 401, "error": "access_token已过期"})
			return
		}
		ctx.Set("userId", claims.ID)
		ctx.Next()
	}
}
