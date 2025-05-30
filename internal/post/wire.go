//go:build wireinject

package post

import (
	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/codepzj/stellux/server/internal/post/internal/repository"
	"github.com/codepzj/stellux/server/internal/post/internal/repository/dao"
	"github.com/codepzj/stellux/server/internal/post/internal/service"
	"github.com/codepzj/stellux/server/internal/post/internal/web"
	"github.com/google/wire"
)

var PostProviders = wire.NewSet(web.NewPostHandler, service.NewPostService, repository.NewPostRepository, dao.NewPostDao,
	wire.Bind(new(service.IPostService), new(*service.PostService)),
	wire.Bind(new(repository.IPostRepository), new(*repository.PostRepository)),
	wire.Bind(new(dao.IPostDao), new(*dao.PostDao)))

func InitPostModule(mongoDB *mongox.Database) *Module {
	panic(wire.Build(
		PostProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
