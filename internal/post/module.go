package post

import (
	"github.com/codepzj/Stellux-Server/internal/label"
	"github.com/codepzj/Stellux-Server/internal/post/internal/service"
	"github.com/codepzj/Stellux-Server/internal/post/internal/web"
)

type (
	Handler     = web.PostHandler
	Service     = service.IPostService
	LabelDomain = label.Domain
	Module      struct {
		Svc Service
		Hdl *Handler
	}
)
