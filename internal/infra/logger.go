package infra

import (
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/codepzj/Stellux-Server/conf"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(cfg *conf.Config) {
	if err := os.MkdirAll("log", os.ModePerm); err != nil {
		log.Fatalf("创建日志目录失败: %v", err)
		return
	}

	loggerOpts := lumberjack.Logger{
		Filename:   cfg.Server.Log.File,
		MaxSize:    10,   // 最大10MB
		MaxBackups: 3,    // 最多保留3个备份
		MaxAge:     30,   // 最多保留30天
		Compress:   true, // 是否压缩
	}

	// 固定东八区（CST）
	cst := time.FixedZone("CST", 8*3600)

	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     parseLevel(cfg.Server.Log.Level),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time().In(cst)
				a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05")) // 标准化时间格式
			}
			return a
		},
	}

	// 日志双写
	multiWriter := io.MultiWriter(&loggerOpts, os.Stdout)
	handler := slog.NewJSONHandler(multiWriter, opts)

	slog.SetDefault(slog.New(handler))
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
