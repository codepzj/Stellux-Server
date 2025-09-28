//go:build wireinject

package user

import (
	"gorm.io/gorm"
	"github.com/codepzj/Stellux-Server/internal/user/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/user/internal/repository/dao"
	"github.com/codepzj/Stellux-Server/internal/user/internal/service"
	"github.com/codepzj/Stellux-Server/internal/user/internal/web"
	"github.com/google/wire"
)

var UserProviders = wire.NewSet(web.NewUserHandler, service.NewUserService, repository.NewUserRepository, dao.NewUserDao,
	wire.Bind(new(service.IUserService), new(*service.UserService)),
	wire.Bind(new(repository.IUserRepository), new(*repository.UserRepository)),
	wire.Bind(new(dao.IUserDao), new(*dao.UserDao)))

func InitUserModule(db *gorm.DB) *Module {
	panic(wire.Build(
		UserProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
