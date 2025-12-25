package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/future-mcp/future-mcp-server/internal/svc"
	"github.com/future-mcp/future-mcp-server/internal/types"
	"github.com/future-mcp/future-mcp-server/pkg/logger"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// MCPHandler MCP JSON-RPC处理器
func MCPHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求
		var request types.MCPRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			logger.Warn("Failed to parse MCP request", logger.Any("error", err))
			httpx.OkJson(w, types.MCPResponse{
				MCPMessage: types.MCPMessage{
					JSONRPC: "2.0",
					ID:      request.ID,
				},
				Error: &types.MCPError{
					Code:    types.MCPInvalidRequest,
					Message: "Invalid JSON-RPC request",
				},
			})
			return
		}

		// 验证协议版本
		if request.JSONRPC != "2.0" {
			httpx.OkJson(w, types.MCPResponse{
				MCPMessage: types.MCPMessage{
					JSONRPC: "2.0",
					ID:      request.ID,
				},
				Error: &types.MCPError{
					Code:    types.MCPInvalidRequest,
					Message: "Unsupported JSON-RPC version",
				},
			})
			return
		}

		// 获取用户上下文
		ctx := extractUserContext(r)

		// 处理请求
		response, err := svcCtx.MCPService.HandleRequest(ctx, &request)
		if err != nil {
			logger.Error("Failed to handle MCP request",
				logger.Any("method", request.Method),
				logger.Any("error", err))

			httpx.OkJson(w, types.MCPResponse{
				MCPMessage: types.MCPMessage{
					JSONRPC: "2.0",
					ID:      request.ID,
				},
				Error: &types.MCPError{
					Code:    types.MCPInternalError,
					Message: "Internal server error",
				},
			})
			return
		}

		httpx.OkJson(w, response)
	}
}


// extractUserContext 从HTTP请求中提取用户上下文
func extractUserContext(r *http.Request) context.Context {
	ctx := r.Context()

	// 从请求头中提取用户信息（如果有的话）
	if userID := r.Header.Get("X-User-ID"); userID != "" {
		ctx = context.WithValue(ctx, "user_id", userID)
	}

	if sessionID := r.Header.Get("X-Session-ID"); sessionID != "" {
		ctx = context.WithValue(ctx, "session_id", sessionID)
	}

	if clientID := r.Header.Get("X-Client-ID"); clientID != "" {
		ctx = context.WithValue(ctx, "client_id", clientID)
	}

	return ctx
}


