package ioc

import (
	"github.com/codepzj/stellux/server/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func InitMiddleWare() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		gin.Recovery(),
		middleware.Cors(),
		middleware.JWT(),
		middleware.Auth(),
		middleware.CacheControlMiddleware(),
	}
}
