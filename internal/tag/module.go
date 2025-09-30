package tag

import (
	"github.com/codepzj/Stellux-Server/internal/tag/internal/service"
	"github.com/codepzj/Stellux-Server/internal/tag/internal/web"
)

type (
	Handler = web.TagHandler
	Service = service.ITagService
	Module   struct {
		Svc Service
		Hdl *Handler
	}
)