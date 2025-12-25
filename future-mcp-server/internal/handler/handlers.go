package handler

import (
	"net/http"

	"github.com/future-mcp/future-mcp-server/internal/svc"
	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, ctx *svc.ServiceContext) {
	// 健康检查路由
	server.AddRoutes([]rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/health",
			Handler: HealthCheck(ctx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/ready",
			Handler: ReadinessCheck(ctx),
		},
	})

	// MCP协议路由
	server.AddRoutes([]rest.Route{
		{
			Method:  http.MethodPost,
			Path:    "/mcp/jsonrpc",
			Handler: MCPHandler(ctx),
		},
	})
}
