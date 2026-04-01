package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	customlogger "aubergine/logger"
)

// GinLogger injects our custom asynchronous logger into the Gin HTTP pipeline.
func GinLogger(l *customlogger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start).String()
		status := c.Writer.Status()

		fields := map[string]string{
			"status":  fmt.Sprintf("%d", status),
			"latency": latency,
			"method":  c.Request.Method,
			"path":    path,
			"ip":      c.ClientIP(),
		}

		if status >= 500 {
			l.Error("Server error on request", fields)
		} else if status >= 400 {
			l.Error("Client error on request", fields)
		} else {
			l.Info("Request handled successfully", fields)
		}
	}
}
