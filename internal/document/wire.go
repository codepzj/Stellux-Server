//go:build wireinject

package document

import (
    "github.com/chenmingyong0423/go-mongox/v2"
	"github.com/codepzj/stellux/server/internal/document/internal/repository"
	"github.com/codepzj/stellux/server/internal/document/internal/repository/dao"
	"github.com/codepzj/stellux/server/internal/document/internal/service"
	"github.com/codepzj/stellux/server/internal/document/internal/web"
	"github.com/codepzj/stellux/server/internal/setting"
	"github.com/google/wire"
)

var DocumentProviders = wire.NewSet(web.NewDocumentHandler, service.NewDocumentService, repository.NewDocumentRepository, dao.NewDocumentDao,
	wire.Bind(new(service.IDocumentService), new(*service.DocumentService)),
	wire.Bind(new(repository.IDocumentRepository), new(*repository.DocumentRepository)),
	wire.Bind(new(dao.IDocumentDao), new(*dao.DocumentDao)))
	
func InitDocumentModule(mongoDB *mongox.Database, settingModule *setting.Module) *Module {
	panic(wire.Build(
		DocumentProviders,
		wire.FieldsOf(new(*setting.Module), "Svc"),
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
