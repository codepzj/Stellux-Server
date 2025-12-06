package config

import (
	"github.com/codepzj/Stellux-Server/internal/config/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/config/internal/repository/dao"
	"github.com/codepzj/Stellux-Server/internal/config/internal/service"
	"github.com/codepzj/Stellux-Server/internal/config/internal/web"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type (
	Handler    = web.ConfigHandler
	Service    = service.IConfigService
	Repository = repository.IConfigRepository
	Dao        = dao.IConfigDao
	Module     struct {
		Svc Service
		Hdl *Handler
	}
)

func New(db *mongo.Database) *Module {
	configDao := dao.NewConfigDao(db)
	configRepository := repository.NewConfigRepository(configDao)
	configService := service.NewConfigService(configRepository)
	configHandler := web.NewConfigHandler(configService)

	return &Module{
		Svc: configService,
		Hdl: configHandler,
	}
}
