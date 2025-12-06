//go:build wireinject

package config

import (
	"github.com/codepzj/Stellux-Server/internal/config/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/config/internal/repository/dao"
	"github.com/codepzj/Stellux-Server/internal/config/internal/service"
	"github.com/codepzj/Stellux-Server/internal/config/internal/web"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ConfigProviders = wire.NewSet(web.NewConfigHandler, service.NewConfigService, repository.NewConfigRepository, dao.NewConfigDao,
	wire.Bind(new(service.IConfigService), new(*service.ConfigService)),
	wire.Bind(new(repository.IConfigRepository), new(*repository.ConfigRepository)),
	wire.Bind(new(dao.IConfigDao), new(*dao.ConfigDao)))

func InitConfigModule(mongoDB *mongo.Database) *Module {
	panic(wire.Build(
		ConfigProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
