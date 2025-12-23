package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/future-mcp/future-mcp-server/internal/service"
	"github.com/future-mcp/future-mcp-server/internal/types"
	"github.com/future-mcp/future-mcp-server/pkg/logger"
	"github.com/gin-gonic/gin"
)

// MCPHandler MCP JSON-RPC处理器
func MCPHandler(mcpService *service.MCPService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析请求
		var request types.MCPRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			logger.Warn("Failed to parse MCP request", logger.Any("error", err))
			c.JSON(http.StatusBadRequest, types.MCPResponse{
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
			c.JSON(http.StatusBadRequest, types.MCPResponse{
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
		ctx := extractUserContext(c)

		// 处理请求
		response, err := mcpService.HandleRequest(ctx, &request)
		if err != nil {
			logger.Error("Failed to handle MCP request",
				logger.Any("method", request.Method),
				logger.Any("error", err))

			c.JSON(http.StatusInternalServerError, types.MCPResponse{
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

		c.JSON(http.StatusOK, response)
	}
}

// MCPSSEHandler MCP SSE处理器
func MCPSSEHandler(mcpService *service.MCPService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置SSE头
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")

		// 获取用户上下文
		ctx := extractUserContext(c)

		// 升级为SSE连接
		c.Writer.Flush()

		// 处理SSE请求（这里简化实现）
		// 实际应该解析SSE格式的请求并返回流式响应

		// 发送初始连接确认
		fmt.Fprintf(c.Writer, "data: %s\n\n", `{"type": "connected", "message": "TALink MCP Server connected"}`)
		c.Writer.Flush()

		// 保持连接（简化实现）
		<-ctx.Done()
	}
}

// extractUserContext 从Gin上下文中提取用户上下文
func extractUserContext(c *gin.Context) context.Context {
	ctx := c.Request.Context()

	// 从请求头中提取用户信息（如果有的话）
	if userID := c.GetHeader("X-User-ID"); userID != "" {
		ctx = context.WithValue(ctx, "user_id", userID)
	}

	if sessionID := c.GetHeader("X-Session-ID"); sessionID != "" {
		ctx = context.WithValue(ctx, "session_id", sessionID)
	}

	if clientID := c.GetHeader("X-Client-ID"); clientID != "" {
		ctx = context.WithValue(ctx, "client_id", clientID)
	}

	return ctx
}

// BatchMCPHandler 批量MCP请求处理器
func BatchMCPHandler(mcpService *service.MCPService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析批量请求
		var requests []types.MCPRequest
		if err := c.ShouldBindJSON(&requests); err != nil {
			logger.Warn("Failed to parse batch MCP requests", logger.Any("error", err))
			c.JSON(http.StatusBadRequest, types.MCPResponse{
				Error: &types.MCPError{
					Code:    types.MCPInvalidRequest,
					Message: "Invalid batch request",
				},
			})
			return
		}

		// 限制批量请求数量
		if len(requests) > 10 {
			c.JSON(http.StatusBadRequest, types.MCPResponse{
				Error: &types.MCPError{
					Code:    types.MCPInvalidRequest,
					Message: "Too many requests in batch (max 10)",
				},
			})
			return
		}

		// 获取用户上下文
		ctx := extractUserContext(c)

		// 处理批量请求
		responses := make([]types.MCPResponse, len(requests))
		for i, request := range requests {
			response, err := mcpService.HandleRequest(ctx, &request)
			if err != nil {
				logger.Error("Failed to handle batch MCP request",
					logger.Any("method", request.Method),
					logger.Any("error", err))

				responses[i] = types.MCPResponse{
					MCPMessage: types.MCPMessage{
						JSONRPC: "2.0",
						ID:      request.ID,
					},
					Error: &types.MCPError{
						Code:    types.MCPInternalError,
						Message: "Internal server error",
					},
				}
			} else {
				responses[i] = *response
			}
		}

		c.JSON(http.StatusOK, responses)
	}
}

// MCPWebSocketHandler WebSocket处理器（暂未实现）
func MCPWebSocketHandler(mcpService *service.MCPService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 暂未实现WebSocket支持
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "WebSocket support not implemented yet",
		})
	}
}
// upgradeToWebSocket 升级为WebSocket连接（这里需要导入gorilla/websocket）
// 注意：实际实现需要导入相应的WebSocket库
func upgradeToWebSocket(c *gin.Context) (interface{}, error) {
	// 这里是简化实现，实际需要使用WebSocket库
	return nil, fmt.Errorf("WebSocket not implemented")
}

// isNormalClose 检查是否是正常关闭 (暂时未使用)
func isNormalClose(err error) bool {
	return false // 简化实现
}

// MCPHealthHandler MCP健康检查处理器
func MCPHealthHandler(mcpService *service.MCPService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查MCP服务状态
		health := map[string]interface{}{
			"status":  "healthy",
			"service": "TALink MCP Server",
			"version": types.MCPProtocolVersion,
		}

		c.JSON(http.StatusOK, health)
	}
}
