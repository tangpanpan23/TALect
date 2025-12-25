package handler

import (
	"net/http"

	"github.com/future-mcp/future-mcp-server/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// HealthCheck 健康检查
func HealthCheck(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpx.OkJson(w, map[string]interface{}{
			"status":  "healthy",
			"service": "TALink MCP Server",
		})
	}
}

// ReadinessCheck 就绪检查
func ReadinessCheck(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 这里可以添加更复杂的就绪检查逻辑
		// 比如检查数据库连接、缓存连接等
		httpx.OkJson(w, map[string]interface{}{
			"status":  "ready",
			"service": "TALink MCP Server",
		})
	}
}
