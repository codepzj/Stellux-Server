package document

import (
	"github.com/codepzj/stellux/server/internal/document/internal/service"
	"github.com/codepzj/stellux/server/internal/document/internal/web"
	"github.com/codepzj/stellux/server/internal/setting"
)

type (
	Handler = web.DocumentHandler
	Service = service.IDocumentService
	SettingService = setting.Service
	Module   struct {
		Svc Service
		Hdl *Handler
	}
)