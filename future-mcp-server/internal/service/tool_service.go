package service

import (
	"context"
	"fmt"

	"github.com/future-mcp/future-mcp-server/internal/types"
)

// ToolServiceImpl 工具服务实现
type ToolServiceImpl struct {
	mcpService *MCPService
}

// NewToolService 创建工具服务
func NewToolService(mcpService *MCPService) ToolService {
	return &ToolServiceImpl{
		mcpService: mcpService,
	}
}

// SetMCPService 设置MCP服务引用
func (s *ToolServiceImpl) SetMCPService(mcpService *MCPService) {
	s.mcpService = mcpService
}

// ExecuteTool 执行工具
func (s *ToolServiceImpl) ExecuteTool(ctx context.Context, toolName string, params map[string]interface{}) (interface{}, error) {
	// 构建MCP工具调用请求
	request := &types.MCPRequest{
		MCPMessage: types.MCPMessage{
			JSONRPC: "2.0",
			ID:      1, // 简化处理
		},
		Method: types.MCPMethodToolsCall,
		Params: map[string]interface{}{
			"name":      toolName,
			"arguments": params,
		},
	}

	// 调用MCP服务
	response, err := s.mcpService.HandleRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute tool: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("tool execution error: %s", response.Error.Message)
	}

	// 解析工具调用响应
	if toolResponse, ok := response.Result.(*types.ToolsCallResponse); ok {
		if len(toolResponse.Content) > 0 {
			return toolResponse.Content[0].Text, nil
		}
	}

	return response.Result, nil
}

// ListAvailableTools 列出可用工具
func (s *ToolServiceImpl) ListAvailableTools() ([]types.Tool, error) {
	tools := s.mcpService.toolRegistry.ListTools()
	result := make([]types.Tool, 0, len(tools))

	for _, tool := range tools {
		result = append(result, types.Tool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		})
	}

	return result, nil
}

// GetToolDefinition 获取工具定义
func (s *ToolServiceImpl) GetToolDefinition(toolName string) (*types.ToolDefinition, error) {
	tool := s.mcpService.toolRegistry.GetTool(toolName)
	if tool == nil {
		return nil, fmt.Errorf("tool not found: %s", toolName)
	}
	return tool, nil
}

// GetToolUsageStatistics 获取工具使用统计
func (s *ToolServiceImpl) GetToolUsageStatistics() (map[string]interface{}, error) {
	// 简化实现，返回基本统计信息
	return map[string]interface{}{
		"total_tools": len(s.mcpService.toolRegistry.ListTools()),
		"tools_active": true,
		"last_updated": "2024-01-01T00:00:00Z",
	}, nil
}
