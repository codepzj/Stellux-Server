package infra

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/codepzj/Stellux-Server/conf"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewPgsql(cfg *conf.Config) *gorm.DB {
	newLogger := logger.New(
		log.New(getLogOutput(cfg), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Duration(cfg.Pgsql.SlowSqlThreshold) * time.Second, // 慢查询阈值，单位秒
			LogLevel:      parseLogLevel(cfg),                                      // 日志级别
		},
	)
	// 初始化gorm
	db, err := gorm.Open(postgres.New(
		postgres.Config{DSN: cfg.Pgsql.Dsn, PreferSimpleProtocol: true}),
		&gorm.Config{Logger: newLogger, NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}},
	)
	if err != nil {
		slog.Error("连接Pgsql失败", "error", err)
		panic(err.Error())
	}
	return db
}

func getLogOutput(cfg *conf.Config) *os.File {
	switch cfg.Pgsql.LogType {
	case 0:
		return os.Stdout
	case 1:
		f, err := os.OpenFile(cfg.Pgsql.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			slog.Error("打开日志文件失败", "error", err)
		}
		return f
	}
	return os.Stdout
}

func parseLogLevel(cfg *conf.Config) logger.LogLevel {
	switch cfg.Pgsql.LogLevel {
	case "silent":
		return logger.Silent
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	}
	return logger.Info
}
