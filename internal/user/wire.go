//go:build wireinject

package user

import (
	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/codepzj/stellux/server/internal/user/internal/repository"
	"github.com/codepzj/stellux/server/internal/user/internal/repository/dao"
	"github.com/codepzj/stellux/server/internal/user/internal/service"
	"github.com/codepzj/stellux/server/internal/user/internal/web"
	"github.com/google/wire"
)

var UserProviders = wire.NewSet(web.NewUserHandler, service.NewUserService, repository.NewUserRepository, dao.NewUserDao,
	wire.Bind(new(service.IUserService), new(*service.UserService)),
	wire.Bind(new(repository.IUserRepository), new(*repository.UserRepository)),
	wire.Bind(new(dao.IUserDao), new(*dao.UserDao)))

func InitUserModule(mongoDB *mongox.Database) *Module {
	panic(wire.Build(
		UserProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
