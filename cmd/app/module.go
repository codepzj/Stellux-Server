package app

import (
	"log"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	engine *gin.Engine
}

func NewHttpServer(engine *gin.Engine) *HttpServer {
	return &HttpServer{
		engine: engine,
	}
}

func (s *HttpServer) Start() {
	if err := s.engine.Run(":9001"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
