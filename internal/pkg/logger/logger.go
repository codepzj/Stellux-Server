package logger

import (
	"io"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

type Env string

const (
	Dev  Env = "dev"
	Prod Env = "prod"
)

type Option struct {
	Env              Env
	Level            string
	FullLogFilename  string
	ErrorLogFilename string
	MaxSize          int
	MaxBackups       int
	MaxAge           int
	Compress         bool
}

func NewLogger(option *Option) {
	var encoder zapcore.Encoder
	if option.Env == Dev {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	} else {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	coreFull := zapcore.NewCore(encoder, getFullLogWriter(option), parseLogLevel(option.Level))
	coreError := zapcore.NewCore(encoder, getErrorLogWriter(option), zapcore.ErrorLevel)

	core := zapcore.NewTee(coreFull, coreError)
	logger = zap.New(core)
	log.Printf("Logger Mode: %s\n", option.Env)
	log.Println("Logger init success...")
}

func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	}
	return zapcore.InfoLevel
}

// 全量日志（控制台 + 文件）
func getFullLogWriter(option *Option) zapcore.WriteSyncer {
	fullLogger := &lumberjack.Logger{
		Filename:   option.FullLogFilename, // 日志文件路径
		MaxSize:    option.MaxSize,         // 单个文件最大 10MB
		MaxBackups: option.MaxBackups,      // 保留 5 个旧文件
		MaxAge:     option.MaxAge,          // 保留 30 天
		Compress:   option.Compress,        // 启用压缩
	}
	ws := io.MultiWriter(fullLogger, os.Stdout)
	return zapcore.AddSync(ws)
}

// 错误日志（仅文件）
func getErrorLogWriter(option *Option) zapcore.WriteSyncer {
	errLogger := &lumberjack.Logger{
		Filename:   option.ErrorLogFilename, // 错误日志文件路径
		MaxSize:    option.MaxSize,          // 单个文件最大 10MB
		MaxBackups: option.MaxBackups,       // 保留 5 个旧文件
		MaxAge:     option.MaxAge,           // 保留 30 天
		Compress:   option.Compress,         // 启用压缩
	}
	return zapcore.AddSync(errLogger)
}

func GetLogger() *zap.Logger {
	return logger
}

func Debug(msg string, fields ...zap.Field) {
	if logger == nil {
		log.Printf("[DEBUG] %s (logger not initialized)\n", msg)
		return
	}
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	if logger == nil {
		log.Printf("[INFO] %s (logger not initialized)\n", msg)
		return
	}
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	if logger == nil {
		log.Printf("[WARN] %s (logger not initialized)\n", msg)
		return
	}
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	if logger == nil {
		log.Printf("[ERROR] %s (logger not initialized)\n", msg)
		return
	}
	logger.Error(msg, fields...)
}

// Fatal logs a message at FatalLevel and then calls os.Exit(1)
func Fatal(msg string, fields ...zap.Field) {
	if logger == nil {
		log.Fatalf("[FATAL] %s (logger not initialized)\n", msg)
		return
	}
	logger.Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	if logger != nil {
		return logger.Sync()
	}
	return nil
}

// WithContext creates a new logger with context fields
func WithContext(ctx string) *zap.Logger {
	if logger == nil {
		log.Println("[WARN] logger not initialized, returning nil")
		return nil
	}
	return logger.With(zap.String("context", ctx))
}

// WithError logs with error field
func WithError(err error) zap.Field {
	return zap.Error(err)
}

// WithString logs with string field
func WithString(key, value string) zap.Field {
	return zap.String(key, value)
}

// WithInt logs with int field
func WithInt(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// WithAny logs with any type field
func WithAny(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}
