package document_content

import (
	"github.com/codepzj/Stellux-Server/internal/document_content/internal/service"
	"github.com/codepzj/Stellux-Server/internal/document_content/internal/web"
)

type (
	Handler = web.DocumentContentHandler
	Service = service.IDocumentContentService
	Module  struct {
		Svc Service
		Hdl *Handler
	}
)
