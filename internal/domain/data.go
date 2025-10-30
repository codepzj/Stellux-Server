package domain

import (
	"github.com/codepzj/Stellux-Server/internal/domain/infra"
	"github.com/codepzj/Stellux-Server/internal/domain/query"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(infra.NewLogger, infra.NewDB, NewData, NewPostRepo)

type Data struct {
	// 注入数据库、缓存、日志、配置等
	query *query.Query
	log   *zap.SugaredLogger
}

func NewData(db *gorm.DB, log *zap.SugaredLogger) *Data {
	return &Data{query: query.Use(db), log: log}
}
