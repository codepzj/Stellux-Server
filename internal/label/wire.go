//go:build wireinject

package label

import (
	"gorm.io/gorm"
	"github.com/codepzj/Stellux-Server/internal/label/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/label/internal/repository/dao"
	"github.com/codepzj/Stellux-Server/internal/label/internal/service"
	"github.com/codepzj/Stellux-Server/internal/label/internal/web"
	"github.com/google/wire"
)

var LabelProviders = wire.NewSet(web.NewLabelHandler, service.NewLabelService, repository.NewLabelRepository, dao.NewLabelDao,
	wire.Bind(new(service.ILabelService), new(*service.LabelService)),
	wire.Bind(new(repository.ILabelRepository), new(*repository.LabelRepository)),
	wire.Bind(new(dao.ILabelDao), new(*dao.LabelDao)))

func InitLabelModule(db *gorm.DB) *Module {
	panic(wire.Build(
		LabelProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
