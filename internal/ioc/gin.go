package ioc

import (
	"github.com/codepzj/Stellux-Server/internal/document"
	"github.com/codepzj/Stellux-Server/internal/document_content"
	"github.com/codepzj/Stellux-Server/internal/file"
	"github.com/codepzj/Stellux-Server/internal/friend"
	"github.com/codepzj/Stellux-Server/internal/label"
	"github.com/codepzj/Stellux-Server/internal/post"
	"github.com/codepzj/Stellux-Server/internal/user"

	"github.com/gin-gonic/gin"
)

// NewGin 初始化gin服务器
func NewGin(userHdl *user.Handler, postHdl *post.Handler, labelHdl *label.Handler, fileHdl *file.Handler, documentHdl *document.Handler, documentContentHdl *document_content.Handler, friendHdl *friend.Handler, middleware []gin.HandlerFunc) *gin.Engine {
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
		documentContentHdl.RegisterGinRoutes(router)
		friendHdl.RegisterGinRoutes(router)
	}

	return router
}
