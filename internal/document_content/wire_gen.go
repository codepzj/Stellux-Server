// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package document_content

import (
	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/codepzj/stellux/server/internal/document_content/internal/repository"
	"github.com/codepzj/stellux/server/internal/document_content/internal/repository/dao"
	"github.com/codepzj/stellux/server/internal/document_content/internal/service"
	"github.com/codepzj/stellux/server/internal/document_content/internal/web"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitDocumentContentModule(mongoDB *mongox.Database) *Module {
	documentContentDao := dao.NewDocumentContentDao(mongoDB)
	documentContentRepository := repository.NewDocumentContentRepository(documentContentDao)
	documentContentService := service.NewDocumentContentService(documentContentRepository)
	documentContentHandler := web.NewDocumentContentHandler(documentContentService)
	module := &Module{
		Svc: documentContentService,
		Hdl: documentContentHandler,
	}
	return module
}

// wire.go:

var DocumentContentProviders = wire.NewSet(web.NewDocumentContentHandler, service.NewDocumentContentService, repository.NewDocumentContentRepository, dao.NewDocumentContentDao, wire.Bind(new(service.IDocumentContentService), new(*service.DocumentContentService)), wire.Bind(new(repository.IDocumentContentRepository), new(*repository.DocumentContentRepository)), wire.Bind(new(dao.IDocumentContentDao), new(*dao.DocumentContentDao)))
