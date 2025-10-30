package apiwrap

import (
	"github.com/gin-gonic/gin"
)

// 通用返回结构
type Response[T any] struct {
	Code int    `json:"code"`
	Data T      `json:"data,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

/* -------- 错误模型 -------- */
// 业务错误（返回 200，但 code != 0）
type BizError struct {
	Code int
	Msg  string
}

func (e *BizError) Error() string { return e.Msg }
func NewBizError(code int, msg string) *BizError {
	return &BizError{Code: code, Msg: msg}
}

// 客户端请求错误（返回 400）
type BadRequest struct {
	Code int
	Msg  string
}

func (e *BadRequest) Error() string { return e.Msg }
func NewBadRequest(msg string) *BadRequest {
	return &BadRequest{Code: -1, Msg: msg}
}

// 服务端内部错误（返回 500）
type InternalError struct {
	Code int
	Msg  string
}

func (e *InternalError) Error() string { return e.Msg }
func NewInternalError(msg string) *InternalError {
	return &InternalError{Code: -1, Msg: msg}
}

// -------- 成功返回 --------

func Success() *Response[any] {
	return &Response[any]{Code: 0, Msg: "操作成功"}
}

func SuccessWithMsg(msg string) *Response[any] {
	return &Response[any]{Code: 0, Msg: msg}
}

func SuccessWithDetail[T any](data T, msg string) *Response[T] {
	return &Response[T]{Code: 0, Data: data, Msg: msg}
}

/* -------- 错误返回 -------- */

func FailWithMsg(msg string) *Response[any] {
	return &Response[any]{Code: -1, Msg: msg}
}

/* -------- 包装器 -------- */
func Wrap[T any](fn func(ctx *gin.Context) (*Response[T], error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err := fn(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}
		ctx.JSON(200, resp)
	}
}

func WrapWithUri[R any, T any](fn func(ctx *gin.Context, req R) (*Response[T], error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(400, FailWithMsg(err.Error()))
			return
		}
		resp, err := fn(ctx, req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		ctx.JSON(200, resp)
	}
}

func WrapWithQuery[R any, T any](fn func(ctx *gin.Context, req R) (*Response[T], error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(400, FailWithMsg(err.Error()))
			return
		}
		resp, err := fn(ctx, req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		ctx.JSON(200, resp)
	}
}

func WrapWithJson[R any, T any](fn func(ctx *gin.Context, req R) (*Response[T], error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R
		if ctx.Request.ContentLength == 0 {
			ctx.JSON(400, FailWithMsg("请求体为空"))
			return
		}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, FailWithMsg(err.Error()))
			return
		}
		resp, err := fn(ctx, req)
		if err != nil {
			handleError(ctx, err)
			return
		}
		ctx.JSON(200, resp)
	}
}

// handleError 错误处理
func handleError(ctx *gin.Context, err error) {
	switch e := err.(type) {
	case *BadRequest:
		ctx.JSON(400, &Response[any]{Code: e.Code, Msg: e.Msg})
	case *BizError:
		ctx.JSON(200, &Response[any]{Code: e.Code, Msg: e.Msg})
	case *InternalError:
		ctx.JSON(500, &Response[any]{Code: e.Code, Msg: e.Msg})
	default:
		ctx.JSON(500, FailWithMsg(e.Error()))
	}
}
