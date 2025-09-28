package infra

import (
	"github.com/codepzj/Stellux-Server/conf"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(cfg *conf.Config) {
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Fatalf("创建日志目录失败: %v", err)
		return
	}

	loggerOpts := lumberjack.Logger{
		Filename:   cfg.Log.Filename,
		MaxSize:    cfg.Log.MaxSize,    // 最大10MB
		MaxBackups: cfg.Log.MaxBackups, // 最多保留3个备份
		MaxAge:     cfg.Log.MaxAge,     // 最多保留30天
		Compress:   cfg.Log.Compress,   // 是否压缩
	}

	// 北京时区
	cst, _ := time.LoadLocation("Asia/Shanghai")

	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     parseLevel(cfg.Log.Level),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time().In(cst)
				a.Value = slog.StringValue(t.Format(time.DateTime)) // 标准化时间格式
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
