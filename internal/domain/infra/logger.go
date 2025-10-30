package infra

import (
	"io"
	"os"

	"github.com/codepzj/Stellux-Server/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	_ = os.MkdirAll("./logs", 0755) // 创建日志目录
}

func NewLogger(cfg *conf.Log) *zap.SugaredLogger {
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	coreFull := zapcore.NewCore(encoder, getFullLogWriter(cfg), zapcore.InfoLevel)
	coreError := zapcore.NewCore(encoder, getErrorLogWriter(cfg), zapcore.ErrorLevel)

	core := zapcore.NewTee(coreFull, coreError)
	logger := zap.New(core)
	sugarLogger := logger.Sugar()
	return sugarLogger
}

// 全量日志（控制台 + 文件）
func getFullLogWriter(cfg *conf.Log) zapcore.WriteSyncer {
	fullLogger := &lumberjack.Logger{
		Filename:   cfg.FullLogFilename, // 日志文件路径
		MaxSize:    cfg.MaxSize,         // 单个文件最大 10MB
		MaxBackups: cfg.MaxBackups,      // 保留 5 个旧文件
		MaxAge:     cfg.MaxAge,          // 保留 30 天
		Compress:   cfg.Compress,        // 启用压缩
	}
	ws := io.MultiWriter(fullLogger, os.Stdout)
	return zapcore.AddSync(ws)
}

// 错误日志（仅文件）
func getErrorLogWriter(cfg *conf.Log) zapcore.WriteSyncer {
	errLogger := &lumberjack.Logger{
		Filename:   cfg.ErrorLogFilename, // 错误日志文件路径
		MaxSize:    cfg.MaxSize,          // 单个文件最大 10MB
		MaxBackups: cfg.MaxBackups,       // 保留 5 个旧文件
		MaxAge:     cfg.MaxAge,           // 保留 30 天
		Compress:   cfg.Compress,         // 启用压缩
	}
	return zapcore.AddSync(errLogger)
}
