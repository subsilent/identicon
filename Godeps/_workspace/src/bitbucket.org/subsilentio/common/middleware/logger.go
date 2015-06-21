package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    log "github.com/mgutz/logxi/v1"
)

// Logger is middleware which logs every request. It does not record the client IP on purpose.
func Logger(l log.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start timer
        start := time.Now()
        path := c.Request.URL.Path

        // Process request
        c.Next()

        // Stop timer
        end := time.Now()
        latency := end.Sub(start)

        // clientIP := c.ClientIP()
        method := c.Request.Method
        statusCode := c.Writer.Status()
        comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

        l.Info("Request", "method", method, "status", statusCode, "latency", latency, "path", path, "comment", comment)
    }
}
