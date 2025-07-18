// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/codepzj/stellux/server/internal/document"
	"github.com/codepzj/stellux/server/internal/document_content"
	"github.com/codepzj/stellux/server/internal/file"
	"github.com/codepzj/stellux/server/internal/friend"
	"github.com/codepzj/stellux/server/internal/ioc"
	"github.com/codepzj/stellux/server/internal/label"
	"github.com/codepzj/stellux/server/internal/post"
	"github.com/codepzj/stellux/server/internal/setting"
	"github.com/codepzj/stellux/server/internal/user"
)

// Injectors from wire.go:

func InitApp() *HttpServer {
	database := ioc.NewMongoDB()
	module := user.InitUserModule(database)
	userHandler := module.Hdl
	settingModule := setting.InitSettingModule(database)
	postModule := post.InitPostModule(database, settingModule)
	postHandler := postModule.Hdl
	labelModule := label.InitLabelModule(database)
	labelHandler := labelModule.Hdl
	fileModule := file.InitFileModule(database)
	fileHandler := fileModule.Hdl
	documentModule := document.InitDocumentModule(database, settingModule)
	documentHandler := documentModule.Hdl
	settingHandler := settingModule.Hdl
	friendModule := friend.InitFriendModule(database)
	friendHandler := friendModule.Hdl
	document_contentModule := document_content.InitDocumentContentModule(database)
	documentContentHandler := document_contentModule.Hdl
	v := ioc.InitMiddleWare()
	engine := ioc.NewGin(userHandler, postHandler, labelHandler, fileHandler, documentHandler, settingHandler, friendHandler, documentContentHandler, v)
	httpServer := NewHttpServer(engine)
	return httpServer
}
