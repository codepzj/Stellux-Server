package apiwrap

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Code int    `json:"code"`
	Data T      `json:"data,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

const (
	RuquestSuccess              = iota // 请求成功
	RuquestBadRequest                  // 请求参数错误
	RuquestUnauthorized                // 未授权
	RuquestForbidden                   // 禁止访问
	RuquestNotFound                    // 未找到
	RuquestInternalServerError         // 服务器错误
	RequestAccessTokenExpired          // access_token已过期
	RequestAccessTokenNotFound         // access_token未找到
	RequestRefreshTokenExpired         // refresh_token已过期
	RequestRefreshTokenNotFound        // refresh_token未找到
	RequestPermissionDenied            // 权限不足
	RequestLoadPermissionFailed        // 权限加载失败
	RequestDocumentNotPublic           // 文档不是公共文档
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

// 通用包装器
func Wrap[T any](fn func(ctx *gin.Context) T) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 无论是否成功，都返回响应状态码200
		ctx.JSON(http.StatusOK, fn(ctx))
	}
}

// 绑定uri请求体[PATCH, PUT, DELETE]
func WrapWithUri[R any, T any](fn func(ctx *gin.Context, req R) T) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, fn(ctx, req))
	}
}

// 绑定query请求体[GET]
func WrapWithQuery[R any, T any](fn func(ctx *gin.Context, req R) T) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, FailWithMsg(RuquestBadRequest, err.Error()))
			return
		}
		// 无论是否成功，都返回响应状态码200
		ctx.JSON(http.StatusOK, fn(ctx, req))
	}
}

// 绑定json请求体[POST]
func WrapWithJson[R any, T any](fn func(ctx *gin.Context, req R) T) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R

		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, FailWithMsg(RuquestBadRequest, err.Error()))
			return
		}

		if len(bodyBytes) == 0 {
			ctx.JSON(http.StatusBadRequest, FailWithMsg(RuquestBadRequest, "请求体为空"))
			return
		}

		// 恢复body内容, 供后续绑定
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, FailWithMsg(RuquestBadRequest, err.Error()))
			return
		}
		// 无论是否成功，都返回响应状态码200
		ctx.JSON(http.StatusOK, fn(ctx, req))
	}
}
