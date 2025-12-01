//go:build wireinject

package comment

import (
	"github.com/codepzj/Stellux-Server/internal/comment/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/comment/internal/repository/dao"
	"github.com/codepzj/Stellux-Server/internal/comment/internal/service"
	"github.com/codepzj/Stellux-Server/internal/comment/internal/web"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var CommentProviders = wire.NewSet(web.NewCommentHandler, service.NewCommentService, repository.NewCommentRepository, dao.NewCommentDao,
	wire.Bind(new(service.ICommentService), new(*service.CommentService)),
	wire.Bind(new(repository.ICommentRepository), new(*repository.CommentRepository)),
	wire.Bind(new(dao.ICommentDao), new(*dao.CommentDao)))

func InitCommentModule(mongoDB *mongo.Database) *Module {
	panic(wire.Build(
		CommentProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
