//go:build wireinject

package file

import (
	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/codepzj/stellux/server/internal/file/internal/repository"
	"github.com/codepzj/stellux/server/internal/file/internal/repository/dao"
	"github.com/codepzj/stellux/server/internal/file/internal/service"
	"github.com/codepzj/stellux/server/internal/file/internal/web"
	"github.com/google/wire"
)

var FileProviders = wire.NewSet(web.NewFileHandler, service.NewFileService, repository.NewFileRepository, dao.NewFileDao,
	wire.Bind(new(service.IFileService), new(*service.FileService)),
	wire.Bind(new(repository.IFileRepository), new(*repository.FileRepository)),
	wire.Bind(new(dao.IFileDao), new(*dao.FileDao)))

func InitFileModule(mongoDB *mongox.Database) *Module {
	panic(wire.Build(
		FileProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
