package document

import (
	"github.com/codepzj/Stellux-Server/internal/document/internal/service"
	"github.com/codepzj/Stellux-Server/internal/document/internal/web"
)

type (
	Handler = web.DocumentHandler
	Service = service.IDocumentService
	Module  struct {
		Svc Service
		Hdl *Handler
	}
)
