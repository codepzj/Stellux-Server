package global

import (
	"flag"
)

func init() {
	// 解析命令行参数
	flag.Parse()
	// 初始化logger
	InitLogger(*Mode)
	// 初始化viper
	InitViper(*Config)
}
