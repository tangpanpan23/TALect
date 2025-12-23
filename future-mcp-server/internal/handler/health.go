package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "TALink MCP Server",
	})
}

// ReadinessCheck 就绪检查
func ReadinessCheck(c *gin.Context) {
	// 这里可以添加更复杂的就绪检查逻辑
	// 比如检查数据库连接、缓存连接等
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"service": "TALink MCP Server",
	})
}
