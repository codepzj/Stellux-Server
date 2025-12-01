//go:build wireinject

package document_content

import (
	"github.com/codepzj/Stellux-Server/internal/document_content/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/document_content/internal/repository/dao"
	"github.com/codepzj/Stellux-Server/internal/document_content/internal/service"
	"github.com/codepzj/Stellux-Server/internal/document_content/internal/web"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var DocumentContentProviders = wire.NewSet(web.NewDocumentContentHandler, service.NewDocumentContentService, repository.NewDocumentContentRepository, dao.NewDocumentContentDao,
	wire.Bind(new(service.IDocumentContentService), new(*service.DocumentContentService)),
	wire.Bind(new(repository.IDocumentContentRepository), new(*repository.DocumentContentRepository)),
	wire.Bind(new(dao.IDocumentContentDao), new(*dao.DocumentContentDao)))

func InitDocumentContentModule(mongoDB *mongo.Database) *Module {
	panic(wire.Build(
		DocumentContentProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
