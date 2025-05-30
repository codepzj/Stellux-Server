package user

import (
	"github.com/codepzj/stellux/server/internal/user/internal/service"
	"github.com/codepzj/stellux/server/internal/user/internal/web"
)

type (
	Handler = web.UserHandler
	Service = service.IUserService
	Module  struct {
		Svc Service
		Hdl *Handler
	}
)
