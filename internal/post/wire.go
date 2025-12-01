//go:build wireinject

package post

import (
	"github.com/codepzj/Stellux-Server/internal/post/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/post/internal/repository/dao"
	"github.com/codepzj/Stellux-Server/internal/post/internal/service"
	"github.com/codepzj/Stellux-Server/internal/post/internal/web"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var PostProviders = wire.NewSet(web.NewPostHandler, service.NewPostService, repository.NewPostRepository, dao.NewPostDao,
	wire.Bind(new(service.IPostService), new(*service.PostService)),
	wire.Bind(new(repository.IPostRepository), new(*repository.PostRepository)),
	wire.Bind(new(dao.IPostDao), new(*dao.PostDao)))

func InitPostModule(mongoDB *mongo.Database) *Module {
	panic(wire.Build(
		PostProviders,
		wire.Struct(new(Module), "Hdl"),
	))
}
