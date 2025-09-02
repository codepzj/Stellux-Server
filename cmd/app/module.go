package app

import (
	"fmt"
	"log"

	"github.com/codepzj/Stellux-Server/conf"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	engine *gin.Engine
	cfg    *conf.Config
}

func NewHttpServer(engine *gin.Engine, cfg *conf.Config) *HttpServer {
	return &HttpServer{
		engine: engine,
		cfg:    cfg,
	}
}

func (s *HttpServer) Start() {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	if err := s.engine.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
