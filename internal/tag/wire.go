//go:build wireinject

package tag

import (
	"gorm.io/gorm"
	"github.com/codepzj/Stellux-Server/internal/tag/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/tag/internal/repository/dao"
	"github.com/codepzj/Stellux-Server/internal/tag/internal/service"
	"github.com/codepzj/Stellux-Server/internal/tag/internal/web"
	"github.com/google/wire"
)

var TagProviders = wire.NewSet(web.NewTagHandler, service.NewTagService, repository.NewTagRepository, dao.NewTagDao,
	wire.Bind(new(service.ITagService), new(*service.TagService)),
	wire.Bind(new(repository.ITagRepository), new(*repository.TagRepository)),
	wire.Bind(new(dao.ITagDao), new(*dao.TagDao)))
	
func InitTagModule(db *gorm.DB) *Module {
	panic(wire.Build(
		TagProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
