package document

import (
	"github.com/codepzj/stellux/server/internal/document/internal/service"
	"github.com/codepzj/stellux/server/internal/document/internal/web"
)

type (
	Handler = web.DocumentHandler
	Service = service.IDocumentService
	Module   struct {
		Svc Service
		Hdl *Handler
	}
)