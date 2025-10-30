package middleware

import (
	"strings"

	"github.com/codepzj/Stellux-Server/pkg/utils"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		access_token := ctx.Request.Header.Get("Authorization")
		// 若非GET请求的token为空
		if access_token == "" || !strings.HasPrefix(access_token, "Bearer ") {
			ctx.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "未携带access_token"})
			return
		}
		claims, err := utils.ParseToken(strings.TrimPrefix(access_token, "Bearer "))
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "access_token已过期"})
			return
		}
		ctx.Set("userId", claims.ID)
		ctx.Next()
	}
}
