package post

import (
	"github.com/codepzj/stellux/server/internal/label"
	"github.com/codepzj/stellux/server/internal/post/internal/service"
	"github.com/codepzj/stellux/server/internal/post/internal/web"
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
