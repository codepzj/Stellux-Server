package apiwrap

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Code int    `json:"code"`
	Data T      `json:"data,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

var (
	RuquestSuccess              = 0  // 请求成功
	RuquestBadRequest           = 1  // 请求参数错误
	RuquestUnauthorized         = 2  // 未授权
	RuquestForbidden            = 3  // 禁止访问
	RuquestNotFound             = 4  // 未找到
	RuquestInternalServerError  = 5  // 服务器错误
	RequestAccessTokenExpired   = 6  // access_token已过期
	RequestAccessTokenNotFound  = 7  // access_token未找到
	RequestRefreshTokenExpired  = 8  // refresh_token已过期
	RequestRefreshTokenNotFound = 9  // refresh_token未找到
	RequestPermissionDenied     = 10 // 权限不足
	RequestLoadPermissionFailed = 11 // 权限加载失败
	RequestDocumentNotPublic    = 12 // 文档不是公共文档
)

func respond[T any](code int, data T, msg string) *Response[T] {
	return &Response[T]{Code: code, Data: data, Msg: msg}
}

func Success() *Response[any] {
	return respond[any](RuquestSuccess, nil, "操作成功")
}

func SuccessWithMsg(msg string) *Response[any] {
	return respond[any](RuquestSuccess, nil, msg)
}

func SuccessWithDetail[T any](data T, msg string) *Response[T] {
	return respond(RuquestSuccess, data, msg)
}

// 请求失败可以指定code业务状态码
func Fail(code int) *Response[any] {
	return respond[any](code, nil, "操作失败")
}

func FailWithMsg(code int, msg string) *Response[any] {
	return respond[any](code, nil, msg)
}

func Wrap[T any](fn func(ctx *gin.Context) T) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 无论是否成功，都返回响应状态码200
		ctx.JSON(http.StatusOK, fn(ctx))
	}
}

func WrapWithUri[R any, T any](fn func(ctx *gin.Context, req R) T) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, FailWithMsg(RuquestBadRequest, err.Error()))
			return
		}
		// 无论是否成功，都返回响应状态码200
		ctx.JSON(http.StatusOK, fn(ctx, req))
	}
}

func WrapWithBody[R any, T any](fn func(ctx *gin.Context, req R) T) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, FailWithMsg(RuquestBadRequest, err.Error()))
			return
		}
		// 无论是否成功，都返回响应状态码200
		ctx.JSON(http.StatusOK, fn(ctx, req))
	}
}
