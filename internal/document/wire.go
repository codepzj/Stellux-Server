//go:build wireinject

package document

import (
	"github.com/codepzj/Stellux-Server/internal/document/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/document/internal/repository/dao"
	"github.com/codepzj/Stellux-Server/internal/document/internal/service"
	"github.com/codepzj/Stellux-Server/internal/document/internal/web"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var DocumentProviders = wire.NewSet(web.NewDocumentHandler, service.NewDocumentService, repository.NewDocumentRepository, dao.NewDocumentDao,
	wire.Bind(new(service.IDocumentService), new(*service.DocumentService)),
	wire.Bind(new(repository.IDocumentRepository), new(*repository.DocumentRepository)),
	wire.Bind(new(dao.IDocumentDao), new(*dao.DocumentDao)))

func InitDocumentModule(mongoDB *mongo.Database) *Module {
	panic(wire.Build(
		DocumentProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
