package server

import (
	"fmt"
	"log"

	"github.com/codepzj/Stellux-Server/conf"
	"github.com/codepzj/Stellux-Server/internal/service"
	"github.com/codepzj/Stellux-Server/pkg/apiwrap"
	"github.com/codepzj/Stellux-Server/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	cfg       *conf.Server
	ginEngine *GinEngine
}

type GinEngine struct {
	engine      *gin.Engine
	postService *service.PostService
}

func NewGin(postService *service.PostService) *GinEngine {
	engine := gin.Default()
	// 中间件
	middlewares := []gin.HandlerFunc{
		gin.Recovery(),
		middleware.Cors(),
	}
	engine.Use(middlewares...)

	// 文章路由
	postRouter := engine.Group("/post")
	{
		postRouter.POST("/create", apiwrap.WrapWithJson(postService.CreatePost))
		postRouter.POST("/update", apiwrap.WrapWithJson(postService.UpdatePost))
		postRouter.POST("/delete", apiwrap.WrapWithUri(postService.DeletePost))
		postRouter.GET("/get", apiwrap.WrapWithUri(postService.GetPost))
	}
	return &GinEngine{
		engine:      engine,
		postService: postService,
	}
}

func (g *GinEngine) Run(port int) error {
	return g.engine.Run(fmt.Sprintf(":%d", port))
}

func NewHttpServer(ginEngine *GinEngine, cfg *conf.Server) *HttpServer {
	return &HttpServer{
		cfg:       cfg,
		ginEngine: ginEngine,
	}
}

func (s *HttpServer) Start() {
	// 启动服务器
	if err := s.ginEngine.Run(s.cfg.Port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
