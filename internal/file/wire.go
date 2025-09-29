//go:build wireinject

package file

import (
	"gorm.io/gorm"
	"github.com/codepzj/Stellux-Server/internal/file/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/file/internal/repository/dao"
	"github.com/codepzj/Stellux-Server/internal/file/internal/service"
	"github.com/codepzj/Stellux-Server/internal/file/internal/web"
	"github.com/google/wire"
)

var FileProviders = wire.NewSet(web.NewFileHandler, service.NewFileService, repository.NewFileRepository, dao.NewFileDao,
	wire.Bind(new(service.IFileService), new(*service.FileService)),
	wire.Bind(new(repository.IFileRepository), new(*repository.FileRepository)),
	wire.Bind(new(dao.IFileDao), new(*dao.FileDao)))

func InitFileModule(db *gorm.DB) *Module {
	panic(wire.Build(
		FileProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
