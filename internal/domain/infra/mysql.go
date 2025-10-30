package infra

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/codepzj/Stellux-Server/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewDB(cfg *conf.Mysql) *gorm.DB {
	newLogger := logger.New(
		log.New(getLogOutput(cfg), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Duration(cfg.SlowSqlThreshold) * time.Second, // 慢查询阈值，单位秒
			LogLevel:      parseMysqlLogLevel(cfg),                           // 日志级别
		},
	)
	// 初始化gorm
	db, err := gorm.Open(mysql.New(
		mysql.Config{DSN: cfg.Dsn}),
		&gorm.Config{Logger: newLogger, NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}},
	)
	if err != nil {
		slog.Error("连接MySQL失败", "error", err)
		panic(err)
	}
	return db
}

func getLogOutput(cfg *conf.Mysql) *os.File {
	switch cfg.LogType {
	case 0:
		return os.Stdout
	case 1:
		f, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			slog.Error("打开日志文件失败", "error", err)
		}
		return f
	}
	return os.Stdout
}

func parseMysqlLogLevel(cfg *conf.Mysql) logger.LogLevel {
	switch cfg.LogLevel {
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
