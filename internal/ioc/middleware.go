package ioc

import (
	"github.com/codepzj/Stellux-Server/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func InitMiddleWare() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		gin.Recovery(),
		middleware.Cors(),
		middleware.CacheControlMiddleware(),
	}
}
