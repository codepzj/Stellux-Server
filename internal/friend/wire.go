//go:build wireinject

package friend

import (
	"github.com/codepzj/Stellux-Server/internal/friend/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/friend/internal/repository/dao"
	"github.com/codepzj/Stellux-Server/internal/friend/internal/service"
	"github.com/codepzj/Stellux-Server/internal/friend/internal/web"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var FriendProviders = wire.NewSet(web.NewFriendHandler, service.NewFriendService, repository.NewFriendRepository, dao.NewFriendDao,
	wire.Bind(new(service.IFriendService), new(*service.FriendService)),
	wire.Bind(new(repository.IFriendRepository), new(*repository.FriendRepository)),
	wire.Bind(new(dao.IFriendDao), new(*dao.FriendDao)))

func InitFriendModule(mongoDB *mongo.Database) *Module {
	panic(wire.Build(
		FriendProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
