package main

import (
	"flag"

	"github.com/codepzj/Stellux-Server/cmd/app"
	"github.com/codepzj/Stellux-Server/conf"
	"github.com/codepzj/Stellux-Server/internal/infra"
)

// 命令行参数
var (
	CfgPath = flag.String("cfg_path", "conf/dev.yaml", "配置文件路径,eg: conf/dev.yaml")
)

func main() {
	// 解析命令行参数
	flag.Parse()

	// 根据环境读取配置文件
	config := conf.GetConfig(*CfgPath)

	// 日志
	infra.InitLogger(config)

	infra.NewPgsql(config)

	// 启动服务
	app.InitApp(config).Start()
}
