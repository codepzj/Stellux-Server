package app

import (
	"flag"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

// 命令行参数
var (
	Mode   = flag.String("mode", "development", "运行模式,eg: development/production")
	Config = flag.String("config", "config/stellux.development.yaml", "配置文件路径,eg: config/stellux.development.yaml")
)

func InitLogger(mode string) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{AddSource: false}

	switch mode {
	case "production":
		if err := os.MkdirAll("log", os.ModePerm); err != nil {
			slog.Error("创建日志目录失败", "error", err)
			return
		}
		writer, err := os.OpenFile("log/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			slog.Error("打开日志文件失败", "error", err)
			return
		}
		opts.Level = slog.LevelInfo
		// 日志双写
		multiWriter := io.MultiWriter(writer, os.Stdout)
		handler = slog.NewJSONHandler(multiWriter, opts)
	case "development":
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		panic("mode参数无效")
	}

	slog.SetDefault(slog.New(handler))
}

func InitViper(config string) {
	viper.SetConfigFile(config)
	slog.Info("使用配置文件", "config", config)
	err := viper.ReadInConfig()
	if err != nil {
		panic("viper读取配置文件失败")
	}
}

func init() {
	// 解析命令行参数
	flag.Parse()
	// 初始化logger
	InitLogger(*Mode)
	// 初始化viper
	InitViper(*Config)
}
