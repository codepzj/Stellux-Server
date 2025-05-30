package label

import (
	"github.com/codepzj/stellux/server/internal/label/internal/domain"
	"github.com/codepzj/stellux/server/internal/label/internal/service"
	"github.com/codepzj/stellux/server/internal/label/internal/web"
)

type (
	Handler = web.LabelHandler
	Service = service.ILabelService
	Domain  = domain.Label
	Module  struct {
		Svc Service
		Hdl *Handler
	}
)
