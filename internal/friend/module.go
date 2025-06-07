package friend

import (
	"github.com/codepzj/stellux/server/internal/friend/internal/service"
	"github.com/codepzj/stellux/server/internal/friend/internal/web"
)

type (
	Handler = web.FriendHandler
	Service = service.IFriendService
	Module   struct {
		Svc Service
		Hdl *Handler
	}
)