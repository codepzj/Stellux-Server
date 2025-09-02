package friend

import (
	"github.com/codepzj/Stellux-Server/internal/friend/internal/service"
	"github.com/codepzj/Stellux-Server/internal/friend/internal/web"
)

type (
	Handler = web.FriendHandler
	Service = service.IFriendService
	Module  struct {
		Svc Service
		Hdl *Handler
	}
)
