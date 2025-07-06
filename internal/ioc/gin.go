package ioc

import (
	"github.com/codepzj/stellux/server/internal/document"
	"github.com/codepzj/stellux/server/internal/file"
	"github.com/codepzj/stellux/server/internal/friend"
	"github.com/codepzj/stellux/server/internal/label"
	"github.com/codepzj/stellux/server/internal/post"
	"github.com/codepzj/stellux/server/internal/setting"
	"github.com/codepzj/stellux/server/internal/user"

	"github.com/gin-gonic/gin"
)

// NewGin 初始化gin服务器
func NewGin(userHdl *user.Handler, postHdl *post.Handler, labelHdl *label.Handler, fileHdl *file.Handler, documentHdl *document.Handler, settingHdl *setting.Handler, friendHdl *friend.Handler, middleware []gin.HandlerFunc) *gin.Engine {
	router := gin.Default()

	// 中间件
	router.Use(middleware...)

	// 初始化路由
	{
		router.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"msg": "stellux后端服务正常运行中",
			})
		})
		userHdl.RegisterGinRoutes(router)
		postHdl.RegisterGinRoutes(router)
		labelHdl.RegisterGinRoutes(router)
		fileHdl.RegisterGinRoutes(router)
		documentHdl.RegisterGinRoutes(router)
		settingHdl.RegisterGinRoutes(router)
		friendHdl.RegisterGinRoutes(router)
	}

	return router
}
