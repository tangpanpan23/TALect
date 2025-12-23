package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/future-mcp/future-mcp-server/pkg/logger"
	"github.com/gin-gonic/gin"
)

// RequestID 请求ID中间件
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("HTTP Request",
			logger.Any("method", c.Request.Method),
			logger.Any("path", path),
			logger.Any("status", statusCode),
			logger.Any("latency", latency.Milliseconds()),
			logger.Any("ip", c.ClientIP()),
			logger.Any("user_agent", c.Request.UserAgent()),
		)
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 简化的实现，实际应该使用更复杂的ID生成逻辑
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
