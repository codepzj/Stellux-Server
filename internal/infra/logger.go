package infra

import (
	"github.com/codepzj/Stellux-Server/conf"
	"github.com/codepzj/Stellux-Server/internal/pkg/logger"
)

func InitLogger(cfg *conf.Config) {
	logger.NewLogger(&logger.Option{
		Env:              logger.Env(cfg.Server.Log.Env),
		Level:            cfg.Server.Log.Level,
		FullLogFilename:  cfg.Server.Log.File,
		ErrorLogFilename: cfg.Server.Log.ErrorFile,
		MaxSize:          cfg.Server.Log.MaxSize,
		MaxBackups:       cfg.Server.Log.MaxBackups,
		MaxAge:           cfg.Server.Log.MaxAge,
		Compress:         cfg.Server.Log.Compress,
	})
}
