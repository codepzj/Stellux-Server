package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CacheControlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/images/") {
			// 设置 1 小时强缓存
			c.Writer.Header().Set("Cache-Control", "public, max-age=3600")
			c.Writer.Header().Set("Expires", time.Now().Add(time.Hour).UTC().Format(http.TimeFormat))
			// 可选：模拟 Last-Modified，视实际资源修改时间而定
			c.Writer.Header().Set("Last-Modified", time.Now().Add(-1*time.Hour).UTC().Format(http.TimeFormat))
		}
		c.Next()
	}
}
