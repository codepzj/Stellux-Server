package file

import (
	"github.com/codepzj/Stellux-Server/internal/file/internal/service"
	"github.com/codepzj/Stellux-Server/internal/file/internal/web"
)

type (
	Handler = web.FileHandler
	Service = service.IFileService
	Module  struct {
		Svc Service
		Hdl *Handler
	}
)
