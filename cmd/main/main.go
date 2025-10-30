package main

import (
	"flag"
	"log"

	"github.com/codepzj/Stellux-Server/conf"
)

// 命令行参数
var (
	CfgPath = flag.String("cfg_path", "../../conf/dev.yaml", "配置文件路径,eg: conf/dev.yaml")
)

func main() {
	// 解析命令行参数
	flag.Parse()

	// 根据环境读取配置文件
	config := conf.GetConfig(*CfgPath)

	// 使用 wire 生成的初始化函数
	server, err := wireApp(config.Server, config.Mysql, config.Log)
	if err != nil {
		log.Fatalf("初始化应用失败: %v", err)
	}

	// 启动服务
	server.Start()
}
