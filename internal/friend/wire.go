//go:build wireinject

package friend

import (
    "github.com/chenmingyong0423/go-mongox/v2"
	"github.com/codepzj/stellux/server/internal/friend/internal/repository"
	"github.com/codepzj/stellux/server/internal/friend/internal/repository/dao"
	"github.com/codepzj/stellux/server/internal/friend/internal/service"
	"github.com/codepzj/stellux/server/internal/friend/internal/web"
	"github.com/google/wire"
)

var FriendProviders = wire.NewSet(web.NewFriendHandler, service.NewFriendService, repository.NewFriendRepository, dao.NewFriendDao,
	wire.Bind(new(service.IFriendService), new(*service.FriendService)),
	wire.Bind(new(repository.IFriendRepository), new(*repository.FriendRepository)),
	wire.Bind(new(dao.IFriendDao), new(*dao.FriendDao)))
	
func InitFriendModule(mongoDB *mongox.Database) *Module {
	panic(wire.Build(
		FriendProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
