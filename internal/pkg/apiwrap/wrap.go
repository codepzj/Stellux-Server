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

type ErrorResponse struct {
	Code int    `json:"code"`          // 错误码
	Msg  string `json:"msg,omitempty"` // 错误信息
}

func respond[T any](code int, data T, msg string) *Response[T] {
	return &Response[T]{Code: code, Data: data, Msg: msg}
}

func Success() *Response[any] {
	return respond[any](0, nil, "操作成功")
}

func SuccessWithMsg(msg string) *Response[any] {
	return respond[any](0, nil, msg)
}

func SuccessWithDetail[T any](data T, msg string) *Response[T] {
	return respond(0, data, msg)
}

// 请求失败可以指定code业务状态码
func Fail(code int) *Response[any] {
	return respond[any](code, nil, "操作失败")
}

func FailWithMsg(code int, msg string) *Response[any] {
	return respond[any](code, nil, msg)
}

// 通用包装器
func Wrap[T any](fn func(ctx *gin.Context) (*Response[T], error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err := fn(ctx)
		if err != nil {
			ctx.JSON(resp.Code, resp)
			return
		}
		ctx.JSON(200, resp)
	}
}

// 绑定uri请求体[PATCH, PUT, DELETE]
func WrapWithUri[R any, T any](fn func(ctx *gin.Context, req R) (*Response[T], error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := fn(ctx, req)
		if err != nil {
			ctx.JSON(resp.Code, resp)
			return
		}
		ctx.JSON(200, resp)
	}
}

// 绑定query请求体[GET]
func WrapWithQuery[R any, T any](fn func(ctx *gin.Context, req R) (*Response[T], error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(400, FailWithMsg(400, err.Error()))
			return
		}
		resp, err := fn(ctx, req)
		if err != nil {
			ctx.JSON(resp.Code, resp)
			return
		}
		ctx.JSON(200, resp)
	}
}

// 绑定json请求体[POST]
func WrapWithJson[R any, T any](fn func(ctx *gin.Context, req R) (*Response[T], error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R

		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, FailWithMsg(400, err.Error()))
			return
		}

		if len(bodyBytes) == 0 {
			ctx.JSON(http.StatusBadRequest, FailWithMsg(400, "请求体为空"))
			return
		}

		// 恢复body内容, 供后续绑定
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, FailWithMsg(400, err.Error()))
			return
		}

		resp, err := fn(ctx, req)
		if err != nil {
			ctx.JSON(resp.Code, resp)
			return
		}
		ctx.JSON(200, resp)
	}
}
