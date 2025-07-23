//go:build wireinject
// +build wireinject

package app

import (
	"github.com/codepzj/stellux/server/internal/comment"
	"github.com/codepzj/stellux/server/internal/document"
	"github.com/codepzj/stellux/server/internal/document_content"
	"github.com/codepzj/stellux/server/internal/friend"
	"github.com/codepzj/stellux/server/internal/ioc"
	"github.com/codepzj/stellux/server/internal/user"

	"github.com/codepzj/stellux/server/internal/file"
	"github.com/codepzj/stellux/server/internal/label"
	"github.com/codepzj/stellux/server/internal/post"
	"github.com/codepzj/stellux/server/internal/setting"
	"github.com/google/wire"
)

func InitApp() *HttpServer {
	wire.Build(
		ioc.InitMiddleWare,
		ioc.NewGin,
		ioc.NewMongoDB,

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

		setting.InitSettingModule,
		wire.FieldsOf(new(*setting.Module), "Hdl"),

		friend.InitFriendModule,
		wire.FieldsOf(new(*friend.Module), "Hdl"),

		comment.InitCommentModule,
		wire.FieldsOf(new(*comment.Module), "Hdl"),

		NewHttpServer,
	)
	return nil
}
