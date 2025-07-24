//go:build wireinject

package setting

import (
	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/codepzj/stellux/server/internal/setting/internal/repository"
	"github.com/codepzj/stellux/server/internal/setting/internal/repository/dao"
	"github.com/codepzj/stellux/server/internal/setting/internal/service"
	"github.com/codepzj/stellux/server/internal/setting/internal/web"
	"github.com/google/wire"
)

var SettingProviders = wire.NewSet(web.NewSettingHandler, service.NewSettingService, repository.NewSettingRepository, dao.NewSettingDao,
	wire.Bind(new(service.ISettingService), new(*service.SettingService)),
	wire.Bind(new(repository.ISettingRepository), new(*repository.SettingRepository)),
	wire.Bind(new(dao.ISettingDao), new(*dao.SettingDao)))

func InitSettingModule(mongoDB *mongox.Database) *Module {
	panic(wire.Build(
		SettingProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
