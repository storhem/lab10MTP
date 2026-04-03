package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger возвращает middleware, которое логирует каждый HTTP-запрос.
// Формат: [METHOD] PATH | status=XXX | duration=Xms | ip=X.X.X.X
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.Printf("[%s] %s | status=%d | duration=%v | ip=%s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
			c.ClientIP(),
		)
	}
}
