package global

import (
	"errors"
	"flag"

	"github.com/casbin/casbin/v2"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/gin-gonic/gin"
)

// 全局变量
var (
	Enforcer *casbin.Enforcer
)

// 命令行参数
var (
	Mode    = flag.String("mode", "development", "运行模式,eg: development/production")
	Config  = flag.String("config", "config/stellux.development.yaml", "配置文件路径,eg: config/stellux.development.yaml")
)

// 错误
var (
	ErrDocumentNotPublic = errors.New("文档不是公共文档")
)


// 常量
var (
	AccessTokenNotFound  = gin.H{"code": apiwrap.RequestAccessTokenNotFound, "msg": "未携带access_token"}
	AccessTokenExpired   = gin.H{"code": apiwrap.RequestAccessTokenExpired, "msg": "access_token已过期"}
	LoadPermissionFailed = gin.H{"code": apiwrap.RequestLoadPermissionFailed, "msg": "权限加载失败"}
	PermissionDenied     = gin.H{"code": apiwrap.RequestPermissionDenied, "msg": "权限不足,禁止访问"}
)
