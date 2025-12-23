package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/future-mcp/future-mcp-server/internal/types"
	"github.com/future-mcp/future-mcp-server/pkg/logger"
	"github.com/google/uuid"
)

// Service MCP服务
type Service struct {
	config         ServiceConfig
	tools          map[string]*types.ToolDefinition
	resources      map[string]*types.ResourceDefinition
	toolRegistry   *ToolRegistry
	resourceRegistry *ResourceRegistry
	mu             sync.RWMutex
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	MaterialService interface{} // 素材服务接口
	ToolService     interface{} // 工具服务接口
	ResourceService interface{} // 资源服务接口
}

// NewService 创建MCP服务
func NewService(config ServiceConfig) *Service {
	s := &Service{
		config:           config,
		tools:            make(map[string]*types.ToolDefinition),
		resources:        make(map[string]*types.ResourceDefinition),
		toolRegistry:     NewToolRegistry(),
		resourceRegistry: NewResourceRegistry(),
	}

	s.registerDefaultTools()
	s.registerDefaultResources()

	return s
}

// HandleJSONRPC 处理JSON-RPC请求
func (s *Service) HandleJSONRPC(ctx context.Context, request *types.MCPRequest) (*types.MCPResponse, error) {
	mcpLogger := logger.NewMCPLogger(getRequestID(request.ID))

	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		mcpLogger.LogMCPResponse(nil, duration)
	}()

	mcpLogger.LogMCPRequest(request.Method, request.Params)

	switch request.Method {
	case types.MCPMethodInitialize:
		return s.handleInitialize(request)
	case types.MCPMethodToolsList:
		return s.handleToolsList(request)
	case types.MCPMethodToolsCall:
		return s.handleToolsCall(ctx, request)
	case types.MCPMethodResourcesList:
		return s.handleResourcesList(request)
	case types.MCPMethodResourcesRead:
		return s.handleResourcesRead(request)
	case types.MCPMethodPing:
		return s.handlePing(request)
	default:
		return s.createErrorResponse(request.ID, types.MCPMethodNotFound, "Method not found")
	}
}

// HandleSSE 处理SSE请求（用于流式响应）
func (s *Service) HandleSSE(ctx context.Context, request *types.MCPRequest) (<-chan *types.MCPResponse, error) {
	// TODO: 实现SSE流式响应
	return nil, fmt.Errorf("SSE not implemented yet")
}

// registerDefaultTools 注册默认工具
func (s *Service) registerDefaultTools() {
	// 搜索工具
	s.toolRegistry.RegisterTool(&types.ToolDefinition{
		Name:        "search_teaching_materials",
		Description: "搜索教学素材，支持关键词、年级、学科等条件筛选",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "搜索关键词",
				},
				"grade": map[string]interface{}{
					"type":        "array",
					"items":       map[string]interface{}{"type": "string"},
					"description": "年级列表",
				},
				"subject": map[string]interface{}{
					"type":        "string",
					"description": "学科",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "返回数量限制",
					"default":     10,
				},
			},
			"required": []string{"query"},
		},
		Handler: s.handleSearchMaterials,
	})

	// 获取素材详情工具
	s.toolRegistry.RegisterTool(&types.ToolDefinition{
		Name:        "get_material_detail",
		Description: "获取教学素材的详细信息",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"material_id": map[string]interface{}{
					"type":        "string",
					"description": "素材ID",
				},
			},
			"required": []string{"material_id"},
		},
		Handler: s.handleGetMaterialDetail,
	})

	// 生成教案工具
	s.toolRegistry.RegisterTool(&types.ToolDefinition{
		Name:        "generate_lesson_plan",
		Description: "基于教学素材生成结构化教案",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"material_ids": map[string]interface{}{
					"type":        "array",
					"items":       map[string]interface{}{"type": "string"},
					"description": "素材ID列表",
				},
				"objectives": map[string]interface{}{
					"type":        "array",
					"items":       map[string]interface{}{"type": "string"},
					"description": "教学目标",
				},
				"grade": map[string]interface{}{
					"type":        "string",
					"description": "年级",
				},
				"duration": map[string]interface{}{
					"type":        "integer",
					"description": "教学时长（分钟）",
				},
			},
			"required": []string{"material_ids", "objectives"},
		},
		Handler: s.handleGenerateLessonPlan,
	})
}

// registerDefaultResources 注册默认资源
func (s *Service) registerDefaultResources() {
	// 课程大纲资源
	s.resourceRegistry.RegisterResource(&types.ResourceDefinition{
		URI:         "curriculum://grade-1/math",
		Name:        "一年级数学课程大纲",
		Description: "一年级数学课程标准和知识体系",
		MimeType:    "application/json",
		Handler:     s.handleCurriculumResource,
	})

	// 知识图谱资源
	s.resourceRegistry.RegisterResource(&types.ResourceDefinition{
		URI:         "knowledge-graph://math/elementary",
		Name:        "小学数学知识图谱",
		Description: "小学数学知识点关系网络",
		MimeType:    "application/json",
		Handler:     s.handleKnowledgeGraphResource,
	})
}

// handleInitialize 处理初始化请求
func (s *Service) handleInitialize(request *types.MCPRequest) (*types.MCPResponse, error) {
	initReq := &types.InitializeRequest{}
	if err := s.parseParams(request.Params, initReq); err != nil {
		return s.createErrorResponse(request.ID, types.MCPInvalidParams, err.Error())
	}

	response := &types.InitializeResponse{
		ProtocolVersion: types.MCPProtocolVersion,
		Capabilities: types.ServerCapabilities{
			Tools: &types.ServerToolsCapability{
				ListChanged: true,
			},
			Resources: &types.ServerResourcesCapability{
				ListChanged: true,
				Subscribe:   true,
			},
		},
		ServerInfo: types.ImplementationInfo{
			Name:    "Future Education MCP Server",
			Version: "1.0.0",
		},
		Instructions: "这是一个好未来教育MCP服务器，提供教学素材搜索、教案生成等功能。",
	}

	return s.createSuccessResponse(request.ID, response)
}

// handleToolsList 处理工具列表请求
func (s *Service) handleToolsList(request *types.MCPRequest) (*types.MCPResponse, error) {
	tools := s.toolRegistry.ListTools()
	toolDefs := make([]types.Tool, 0, len(tools))

	for _, tool := range tools {
		toolDefs = append(toolDefs, types.Tool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		})
	}

	response := &types.ToolsListResponse{
		Tools: toolDefs,
	}

	return s.createSuccessResponse(request.ID, response)
}

// handleToolsCall 处理工具调用请求
func (s *Service) handleToolsCall(ctx context.Context, request *types.MCPRequest) (*types.MCPResponse, error) {
	callReq := &types.ToolsCallRequest{}
	if err := s.parseParams(request.Params, callReq); err != nil {
		return s.createErrorResponse(request.ID, types.MCPInvalidParams, err.Error())
	}

	tool := s.toolRegistry.GetTool(callReq.Name)
	if tool == nil {
		return s.createErrorResponse(request.ID, types.MCPMethodNotFound, "Tool not found")
	}

	toolContext := &types.ToolContext{
		UserID:     getUserIDFromContext(ctx),
		SessionID:  getSessionIDFromContext(ctx),
		RequestID:  getRequestID(request.ID),
		StartTime:  time.Now(),
		Parameters: map[string]interface{}{},
	}

	result, err := tool.Handler(toolContext, callReq.Arguments)
	if err != nil {
		return s.createErrorResponse(request.ID, types.MCPInternalError, err.Error())
	}

	return s.createSuccessResponse(request.ID, result)
}

// handleResourcesList 处理资源列表请求
func (s *Service) handleResourcesList(request *types.MCPRequest) (*types.MCPResponse, error) {
	resources := s.resourceRegistry.ListResources()
	resourceDefs := make([]types.Resource, 0, len(resources))

	for _, resource := range resources {
		resourceDefs = append(resourceDefs, types.Resource{
			URI:         resource.URI,
			Name:        resource.Name,
			Description: resource.Description,
			MimeType:    resource.MimeType,
		})
	}

	response := &types.ResourcesListResponse{
		Resources: resourceDefs,
	}

	return s.createSuccessResponse(request.ID, response)
}

// handleResourcesRead 处理资源读取请求
func (s *Service) handleResourcesRead(request *types.MCPRequest) (*types.MCPResponse, error) {
	readReq := &types.ResourcesReadRequest{}
	if err := s.parseParams(request.Params, readReq); err != nil {
		return s.createErrorResponse(request.ID, types.MCPInvalidParams, err.Error())
	}

	resource := s.resourceRegistry.GetResource(readReq.URI)
	if resource == nil {
		return s.createErrorResponse(request.ID, types.MCPInvalidParams, "Resource not found")
	}

	result, err := resource.Handler(readReq.URI)
	if err != nil {
		return s.createErrorResponse(request.ID, types.MCPInternalError, err.Error())
	}

	return s.createSuccessResponse(request.ID, result)
}

// handlePing 处理ping请求
func (s *Service) handlePing(request *types.MCPRequest) (*types.MCPResponse, error) {
	return s.createSuccessResponse(request.ID, map[string]string{"status": "pong"})
}

// 工具处理器实现
func (s *Service) handleSearchMaterials(ctx *types.ToolContext, args interface{}) (*types.ToolsCallResponse, error) {
	// TODO: 实现素材搜索逻辑
	return &types.ToolsCallResponse{
		Content: []types.Content{
			{
				Type: "text",
				Text: "素材搜索功能开发中...",
			},
		},
		IsError: false,
	}, nil
}

func (s *Service) handleGetMaterialDetail(ctx *types.ToolContext, args interface{}) (*types.ToolsCallResponse, error) {
	// TODO: 实现获取素材详情逻辑
	return &types.ToolsCallResponse{
		Content: []types.Content{
			{
				Type: "text",
				Text: "获取素材详情功能开发中...",
			},
		},
		IsError: false,
	}, nil
}

func (s *Service) handleGenerateLessonPlan(ctx *types.ToolContext, args interface{}) (*types.ToolsCallResponse, error) {
	// TODO: 实现生成教案逻辑
	return &types.ToolsCallResponse{
		Content: []types.Content{
			{
				Type: "text",
				Text: "生成教案功能开发中...",
			},
		},
		IsError: false,
	}, nil
}

// 资源处理器实现
func (s *Service) handleCurriculumResource(uri string) (*types.ResourcesReadResponse, error) {
	// TODO: 实现课程大纲资源逻辑
	return &types.ResourcesReadResponse{
		Contents: []types.ResourceContent{
			{
				URI:      uri,
				MimeType: "application/json",
				Text:     `{"message": "课程大纲资源开发中"}`,
			},
		},
	}, nil
}

func (s *Service) handleKnowledgeGraphResource(uri string) (*types.ResourcesReadResponse, error) {
	// TODO: 实现知识图谱资源逻辑
	return &types.ResourcesReadResponse{
		Contents: []types.ResourceContent{
			{
				URI:      uri,
				MimeType: "application/json",
				Text:     `{"message": "知识图谱资源开发中"}`,
			},
		},
	}, nil
}

// 辅助方法
func (s *Service) parseParams(params interface{}, target interface{}) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func (s *Service) createSuccessResponse(id interface{}, result interface{}) (*types.MCPResponse, error) {
	return &types.MCPResponse{
		MCPMessage: types.MCPMessage{
			JSONRPC: "2.0",
			ID:      id,
		},
		Result: result,
	}, nil
}

func (s *Service) createErrorResponse(id interface{}, code int, message string) (*types.MCPResponse, error) {
	return &types.MCPResponse{
		MCPMessage: types.MCPMessage{
			JSONRPC: "2.0",
			ID:      id,
		},
		Error: &types.MCPError{
			Code:    code,
			Message: message,
		},
	}, nil
}

// 上下文获取辅助函数
func getUserIDFromContext(ctx context.Context) uuid.UUID {
	// TODO: 从上下文获取用户ID
	return uuid.Nil
}

func getSessionIDFromContext(ctx context.Context) string {
	// TODO: 从上下文获取会话ID
	return ""
}

func getRequestID(id interface{}) string {
	if id == nil {
		return ""
	}
	switch v := id.(type) {
	case string:
		return v
	case int, int64:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
