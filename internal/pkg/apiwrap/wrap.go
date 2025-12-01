package apiwrap

import (
	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Code  int    `json:"code"`
	Data  T      `json:"data,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Error string `json:"error,omitempty"`
}

func JsonBack[T any](code int, msg string, data T) *Response[T] {
	return &Response[T]{Code: code, Msg: msg, Data: data}
}

func Wrap[T any](fn func(ctx *gin.Context) (int, string, T)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code, msg, data := fn(ctx)
		if code != 200 {
			ctx.JSON(200, &Response[T]{Code: code, Error: msg})
			return
		}
		ctx.JSON(200, &Response[T]{Code: code, Msg: msg, Data: data})
	}
}

func WrapWithUri[R any, T any](fn func(ctx *gin.Context, req R) (int, string, T)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(200, &Response[T]{Code: 400, Error: err.Error()})
			return
		}
		code, msg, data := fn(ctx, req)
		if code != 200 {
			ctx.JSON(200, &Response[T]{Code: code, Error: msg})
			return
		}
		ctx.JSON(200, &Response[T]{Code: code, Msg: msg, Data: data})
	}
}

func WrapWithQuery[R any, T any](fn func(ctx *gin.Context, req R) (int, string, T)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(200, &Response[T]{Code: 400, Error: err.Error()})
			return
		}
		code, msg, data := fn(ctx, req)
		if code != 200 {
			ctx.JSON(200, &Response[T]{Code: code, Error: msg})
			return
		}
		ctx.JSON(200, &Response[T]{Code: code, Msg: msg, Data: data})
	}
}

func WrapWithJson[R any, T any](fn func(ctx *gin.Context, req R) (int, string, T)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R
		if ctx.Request.ContentLength == 0 {
			ctx.JSON(200, &Response[T]{Code: 400, Error: "request body is empty"})
			return
		}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(200, &Response[T]{Code: 400, Error: err.Error()})
			return
		}
		code, msg, data := fn(ctx, req)
		if code != 200 {
			ctx.JSON(200, &Response[T]{Code: code, Error: msg})
			return
		}
		ctx.JSON(200, &Response[T]{Code: code, Msg: msg, Data: data})
	}
}
