//go:build wireinject
// +build wireinject

package main

import (
	"github.com/codepzj/Stellux-Server/conf"
	"github.com/codepzj/Stellux-Server/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/server"
	"github.com/codepzj/Stellux-Server/internal/service"
	"github.com/google/wire"
)

// wireApp 初始化应用
func wireApp(confServer *conf.Server, confMysql *conf.Mysql, confLog *conf.Log) (*server.HttpServer, error) {
	panic(wire.Build(
		server.ProviderSet,
		service.ProviderSet,
		domain.ProviderSet,
	))
}
