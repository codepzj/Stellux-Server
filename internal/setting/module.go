package setting

import (
	"github.com/codepzj/stellux/server/internal/setting/internal/service"
	"github.com/codepzj/stellux/server/internal/setting/internal/web"
)

type (
	Handler = web.SettingHandler
	Service = service.ISettingService
	Module  struct {
		Svc Service
		Hdl *Handler
	}
)
