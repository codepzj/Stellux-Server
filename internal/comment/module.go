package comment

import (
	"github.com/codepzj/Stellux-Server/internal/comment/internal/service"
	"github.com/codepzj/Stellux-Server/internal/comment/internal/web"
)

type (
	Handler = web.CommentHandler
	Service = service.ICommentService
	Module  struct {
		Svc Service
		Hdl *Handler
	}
)
