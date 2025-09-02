//go:build wireinject
// +build wireinject

package app

import (
	"github.com/codepzj/Stellux-Server/conf"
	"github.com/codepzj/Stellux-Server/internal/document"
	"github.com/codepzj/Stellux-Server/internal/document_content"
	"github.com/codepzj/Stellux-Server/internal/infra"
	"github.com/codepzj/Stellux-Server/internal/ioc"
	"github.com/codepzj/Stellux-Server/internal/user"

	"github.com/codepzj/Stellux-Server/internal/file"
	"github.com/codepzj/Stellux-Server/internal/friend"
	"github.com/codepzj/Stellux-Server/internal/label"
	"github.com/codepzj/Stellux-Server/internal/post"
	"github.com/google/wire"
)

// 基础设施
var InfraProvider = wire.NewSet(
	infra.NewMongoDB,
)

// 控制反转
var IocProvider = wire.NewSet(
	ioc.InitMiddleWare,
	ioc.NewGin,
)

func InitApp(cfg *conf.Config) *HttpServer {
	wire.Build(
		InfraProvider,
		IocProvider,

		user.InitUserModule,
		wire.FieldsOf(new(*user.Module), "Hdl"),

		post.InitPostModule,
		wire.FieldsOf(new(*post.Module), "Hdl"),

		label.InitLabelModule,
		wire.FieldsOf(new(*label.Module), "Hdl"),

		file.InitFileModule,
		wire.FieldsOf(new(*file.Module), "Hdl"),

		document.InitDocumentModule,
		wire.FieldsOf(new(*document.Module), "Hdl"),

		document_content.InitDocumentContentModule,
		wire.FieldsOf(new(*document_content.Module), "Hdl"),

		friend.InitFriendModule,
		wire.FieldsOf(new(*friend.Module), "Hdl"),

		NewHttpServer,
	)
	return nil
}
