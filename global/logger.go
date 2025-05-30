package global

import (
	"io"
	"log/slog"
	"os"
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
