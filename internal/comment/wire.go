//go:build wireinject

package comment

import (
    "github.com/chenmingyong0423/go-mongox/v2"
	"github.com/codepzj/stellux/server/internal/comment/internal/repository"
	"github.com/codepzj/stellux/server/internal/comment/internal/repository/dao"
	"github.com/codepzj/stellux/server/internal/comment/internal/service"
	"github.com/codepzj/stellux/server/internal/comment/internal/web"
	"github.com/google/wire"
)

var CommentProviders = wire.NewSet(web.NewCommentHandler, service.NewCommentService, repository.NewCommentRepository, dao.NewCommentDao,
	wire.Bind(new(service.ICommentService), new(*service.CommentService)),
	wire.Bind(new(repository.ICommentRepository), new(*repository.CommentRepository)),
	wire.Bind(new(dao.ICommentDao), new(*dao.CommentDao)))
	
func InitCommentModule(mongoDB *mongox.Database) *Module {
	panic(wire.Build(
		CommentProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
