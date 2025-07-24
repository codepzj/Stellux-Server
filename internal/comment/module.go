package comment

import (
	"github.com/codepzj/stellux/server/internal/comment/internal/service"
	"github.com/codepzj/stellux/server/internal/comment/internal/web"
)

type (
	Handler = web.CommentHandler
	Service = service.ICommentService
	Module  struct {
		Svc Service
		Hdl *Handler
	}
)
