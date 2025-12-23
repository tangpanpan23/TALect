package service

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

// MCPService MCP服务
type MCPService struct {
	config         *MCPServiceConfig
	tools          map[string]*types.ToolDefinition
	resources      map[string]*types.ResourceDefinition
	toolRegistry   *ToolRegistry
	resourceRegistry *ResourceRegistry
	subscriptionManager *SubscriptionManager
	mu             sync.RWMutex
}

// MCPServiceConfig MCP服务配置
type MCPServiceConfig struct {
	MaterialService MaterialService
	ToolService     ToolService
	ResourceService ResourceService
	UserService     UserService
}

// NewMCPService 创建MCP服务
func NewMCPService(config *MCPServiceConfig) *MCPService {
	s := &MCPService{
		config:              config,
		tools:               make(map[string]*types.ToolDefinition),
		resources:           make(map[string]*types.ResourceDefinition),
		toolRegistry:        NewToolRegistry(),
		resourceRegistry:    NewResourceRegistry(),
		subscriptionManager: NewSubscriptionManager(),
	}

	s.registerDefaultTools()
	s.registerDefaultResources()

	return s
}

// HandleRequest 处理MCP请求
func (s *MCPService) HandleRequest(ctx context.Context, request *types.MCPRequest) (*types.MCPResponse, error) {
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
	case types.MCPMethodResourcesSubscribe:
		return s.handleResourcesSubscribe(ctx, request)
	case types.MCPMethodResourcesUnsubscribe:
		return s.handleResourcesUnsubscribe(ctx, request)
	case types.MCPMethodPing:
		return s.handlePing(request)
	default:
		return s.createErrorResponse(request.ID, types.MCPMethodNotFound, "Method not found")
	}
}

// registerDefaultTools 注册默认工具
func (s *MCPService) registerDefaultTools() {
	// 检索类工具
	s.toolRegistry.RegisterTool(&types.ToolDefinition{
		Name:        "search_teaching_materials",
		Description: "按关键词搜索教学素材，支持学而思培优体系",
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
					"description": "年级列表 (学而思体系)",
				},
				"subject": map[string]interface{}{
					"type":        "string",
					"description": "学科",
					"enum":        []string{"math", "chinese", "english", "physics", "chemistry", "biology"},
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "返回数量限制",
					"default":     10,
					"maximum":     50,
				},
			},
			"required": []string{"query"},
		},
		Handler: s.handleSearchMaterials,
	})

	s.toolRegistry.RegisterTool(&types.ToolDefinition{
		Name:        "search_by_grade_subject",
		Description: "按年级学科筛选教学素材 (学而思培优体系)",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"grade": map[string]interface{}{
					"type":        "string",
					"description": "年级",
					"enum":        []string{"grade_1", "grade_2", "grade_3", "grade_4", "grade_5", "grade_6", "grade_7", "grade_8", "grade_9"},
				},
				"subject": map[string]interface{}{
					"type":        "string",
					"description": "学科",
					"enum":        []string{"math", "chinese", "english", "physics", "chemistry", "biology"},
				},
				"difficulty": map[string]interface{}{
					"type":        "string",
					"description": "难度级别",
					"enum":        []string{"easy", "medium", "hard"},
					"default":     "medium",
				},
				"teaching_stage": map[string]interface{}{
					"type":        "string",
					"description": "教学阶段",
					"enum":        []string{"basic", "advanced", "olympic"},
				},
			},
			"required": []string{"grade", "subject"},
		},
		Handler: s.handleSearchByGradeSubject,
	})

	s.toolRegistry.RegisterTool(&types.ToolDefinition{
		Name:        "get_recommended_materials",
		Description: "基于学习数据个性化推荐 (AI算法)",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user_id": map[string]interface{}{
					"type":        "string",
					"description": "用户ID",
				},
				"learning_goals": map[string]interface{}{
					"type":        "array",
					"items":       map[string]interface{}{"type": "string"},
					"description": "学习目标",
				},
				"history_records": map[string]interface{}{
					"type":        "array",
					"items":       map[string]interface{}{"type": "string"},
					"description": "历史学习记录",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "推荐数量",
					"default":     5,
					"maximum":     20,
				},
			},
			"required": []string{"user_id"},
		},
		Handler: s.handleGetRecommendedMaterials,
	})

	// 内容类工具
	s.toolRegistry.RegisterTool(&types.ToolDefinition{
		Name:        "get_material_detail",
		Description: "获取教学素材详细信息 (包含教学元数据)",
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

	s.toolRegistry.RegisterTool(&types.ToolDefinition{
		Name:        "get_related_materials",
		Description: "获取相关素材 (知识图谱关联)",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"material_id": map[string]interface{}{
					"type":        "string",
					"description": "素材ID",
				},
				"relation_type": map[string]interface{}{
					"type":        "string",
					"description": "关联类型",
					"enum":        []string{"prerequisite", "followup", "similar", "complementary"},
					"default":     "similar",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "返回数量",
					"default":     5,
				},
			},
			"required": []string{"material_id"},
		},
		Handler: s.handleGetRelatedMaterials,
	})

	// 生成类工具
	s.toolRegistry.RegisterTool(&types.ToolDefinition{
		Name:        "generate_lesson_plan",
		Description: "生成个性化教案 (学而思教研标准)",
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
				"student_level": map[string]interface{}{
					"type":        "string",
					"description": "学生水平",
					"enum":        []string{"beginner", "intermediate", "advanced"},
				},
				"duration": map[string]interface{}{
					"type":        "integer",
					"description": "教学时长（分钟）",
					"default":     45,
				},
			},
			"required": []string{"material_ids", "objectives", "grade"},
		},
		Handler: s.handleGenerateLessonPlan,
	})

	s.toolRegistry.RegisterTool(&types.ToolDefinition{
		Name:        "generate_exercises",
		Description: "生成智能练习题 (奥数/竞赛专项)",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"material_id": map[string]interface{}{
					"type":        "string",
					"description": "基于的素材ID",
				},
				"exercise_type": map[string]interface{}{
					"type":        "string",
					"description": "练习题类型",
					"enum":        []string{"practice", "homework", "quiz", "olympic", "competition"},
				},
				"difficulty": map[string]interface{}{
					"type":        "string",
					"description": "难度级别",
					"enum":        []string{"easy", "medium", "hard", "challenge"},
				},
				"knowledge_points": map[string]interface{}{
					"type":        "array",
					"items":       map[string]interface{}{"type": "string"},
					"description": "涉及的知识点",
				},
				"count": map[string]interface{}{
					"type":        "integer",
					"description": "生成题目数量",
					"default":     5,
					"maximum":     20,
				},
			},
			"required": []string{"material_id", "exercise_type"},
		},
		Handler: s.handleGenerateExercises,
	})
}

// registerDefaultResources 注册默认资源
func (s *MCPService) registerDefaultResources() {
	// 课程大纲资源
	s.resourceRegistry.RegisterResource(&types.ResourceDefinition{
		URI:         "curriculum://grade-1/math",
		Name:        "一年级数学课程大纲 (学而思标准)",
		Description: "一年级数学课程框架，包含知识体系、教学目标和评估标准",
		MimeType:    "application/json",
		Handler:     s.handleCurriculumResource,
	})

	s.resourceRegistry.RegisterResource(&types.ResourceDefinition{
		URI:         "curriculum://grade-2/math",
		Name:        "二年级数学课程大纲 (学而思标准)",
		Description: "二年级数学课程框架，包含知识体系、教学目标和评估标准",
		MimeType:    "application/json",
		Handler:     s.handleCurriculumResource,
	})

	// 知识图谱资源
	s.resourceRegistry.RegisterResource(&types.ResourceDefinition{
		URI:         "knowledge-graph://math/elementary",
		Name:        "小学数学知识图谱",
		Description: "小学数学知识点关联网络，包含概念关系和学习路径",
		MimeType:    "application/json",
		Handler:     s.handleKnowledgeGraphResource,
	})

	// 教学模板资源
	s.resourceRegistry.RegisterResource(&types.ResourceDefinition{
		URI:         "template://lesson-plan/5e-model",
		Name:        "5E教学模型模板",
		Description: "基于Engage-Explore-Explain-Elaborate-Evaluate的教学流程模板",
		MimeType:    "application/json",
		Handler:     s.handleTeachingTemplateResource,
	})
}

// handleInitialize 处理初始化请求
func (s *MCPService) handleInitialize(request *types.MCPRequest) (*types.MCPResponse, error) {
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
			Name:    "TALink MCP Server",
			Version: "1.0.0",
		},
		Instructions: "TALink是好未来AI教育基础设施的核心组件，提供基于MCP协议的教育内容智能访问服务。支持教学素材搜索、个性化推荐、教案生成等AI教育工具。",
	}

	return s.createSuccessResponse(request.ID, response)
}

// handleToolsList 处理工具列表请求
func (s *MCPService) handleToolsList(request *types.MCPRequest) (*types.MCPResponse, error) {
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
func (s *MCPService) handleToolsCall(ctx context.Context, request *types.MCPRequest) (*types.MCPResponse, error) {
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
func (s *MCPService) handleResourcesList(request *types.MCPRequest) (*types.MCPResponse, error) {
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
func (s *MCPService) handleResourcesRead(request *types.MCPRequest) (*types.MCPResponse, error) {
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

// handleResourcesSubscribe 处理资源订阅请求
func (s *MCPService) handleResourcesSubscribe(ctx context.Context, request *types.MCPRequest) (*types.MCPResponse, error) {
	subscribeReq := &types.ResourcesSubscribeRequest{}
	if err := s.parseParams(request.Params, subscribeReq); err != nil {
		return s.createErrorResponse(request.ID, types.MCPInvalidParams, err.Error())
	}

	// 检查资源是否存在
	if s.resourceRegistry.GetResource(subscribeReq.URI) == nil {
		return s.createErrorResponse(request.ID, types.MCPInvalidParams, "Resource not found")
	}

	// 创建订阅（这里简化处理，实际应该基于用户ID）
	clientID := getClientIDFromContext(ctx)
	ch := s.subscriptionManager.Subscribe(subscribeReq.URI, clientID)

	// 在后台监听资源更新（简化实现）
	go s.monitorResourceUpdates(subscribeReq.URI, ch)

	return s.createSuccessResponse(request.ID, map[string]string{"status": "subscribed"})
}

// handleResourcesUnsubscribe 处理资源取消订阅请求
func (s *MCPService) handleResourcesUnsubscribe(ctx context.Context, request *types.MCPRequest) (*types.MCPResponse, error) {
	unsubscribeReq := &types.ResourcesUnsubscribeRequest{}
	if err := s.parseParams(request.Params, unsubscribeReq); err != nil {
		return s.createErrorResponse(request.ID, types.MCPInvalidParams, err.Error())
	}

	clientID := getClientIDFromContext(ctx)
	s.subscriptionManager.Unsubscribe(unsubscribeReq.URI, clientID)

	return s.createSuccessResponse(request.ID, map[string]string{"status": "unsubscribed"})
}

// handlePing 处理ping请求
func (s *MCPService) handlePing(request *types.MCPRequest) (*types.MCPResponse, error) {
	return s.createSuccessResponse(request.ID, map[string]string{"status": "pong"})
}

// 工具处理器实现（这里提供基础实现，实际需要调用具体服务）
func (s *MCPService) handleSearchMaterials(ctx *types.ToolContext, args interface{}) (*types.ToolsCallResponse, error) {
	// 解析参数
	var params struct {
		Query  string   `json:"query"`
		Grade  []string `json:"grade,omitempty"`
		Subject string  `json:"subject,omitempty"`
		Limit  int      `json:"limit,omitempty"`
	}

	// 直接将args转换为JSON并解析
	argsData, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	if err := json.Unmarshal(argsData, &params); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	// 调用素材服务
	result, err := s.config.MaterialService.SearchMaterials(ctx.UserID, types.SearchMaterialsRequest{
		Query:   params.Query,
		Grade:   convertToGradeLevels(params.Grade),
		Subject: types.Subject(params.Subject),
		Pagination: types.PaginationRequest{
			Page:     1,
			PageSize: params.Limit,
		},
	})
	if err != nil {
		return nil, err
	}

	// 格式化响应
	materials := make([]string, len(result.Materials))
	for i, material := range result.Materials {
		materials[i] = fmt.Sprintf("%s (ID: %s)", material.Title, material.ID)
	}

	return &types.ToolsCallResponse{
		Content: []types.Content{
			{
				Type: "text",
				Text: fmt.Sprintf("找到 %d 个相关教学素材：\n%s", len(materials), fmt.Sprintf("%v", materials)),
			},
		},
		IsError: false,
	}, nil
}

// 其他工具处理器的基础实现
func (s *MCPService) handleSearchByGradeSubject(ctx *types.ToolContext, args interface{}) (*types.ToolsCallResponse, error) {
	return &types.ToolsCallResponse{
		Content: []types.Content{
			{
				Type: "text",
				Text: "年级学科筛选功能正在开发中...",
			},
		},
		IsError: false,
	}, nil
}

func (s *MCPService) handleGetRecommendedMaterials(ctx *types.ToolContext, args interface{}) (*types.ToolsCallResponse, error) {
	return &types.ToolsCallResponse{
		Content: []types.Content{
			{
				Type: "text",
				Text: "个性化推荐功能正在开发中...",
			},
		},
		IsError: false,
	}, nil
}

func (s *MCPService) handleGetMaterialDetail(ctx *types.ToolContext, args interface{}) (*types.ToolsCallResponse, error) {
	// 解析参数
	var params struct {
		MaterialID string `json:"material_id"`
	}

	argsData, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	if err := json.Unmarshal(argsData, &params); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	if params.MaterialID == "" {
		return nil, fmt.Errorf("material_id is required")
	}

	// 解析素材ID
	materialID, err := uuid.Parse(params.MaterialID)
	if err != nil {
		return nil, fmt.Errorf("invalid material_id format: %w", err)
	}

	// 调用素材服务获取详情
	detail, err := s.config.MaterialService.GetMaterialDetail(ctx.UserID, materialID)
	if err != nil {
		return nil, err
	}

	// 格式化响应
	response := fmt.Sprintf(`素材详情：
标题: %s
描述: %s
类型: %s
年级: %v
学科: %s
难度: %s
时长: %s
查看次数: %d
评分: %.1f (%d人评价)
标签: %v
教学目标: %v`,
		detail.TeachingMaterial.Title,
		detail.TeachingMaterial.Description,
		detail.TeachingMaterial.Type,
		detail.TeachingMaterial.GradeLevels,
		detail.TeachingMaterial.Subject,
		detail.TeachingMaterial.Difficulty,
		formatDuration(detail.TeachingMaterial.Metadata.Duration),
		detail.TeachingMaterial.Statistics.ViewCount,
		detail.TeachingMaterial.Statistics.AverageRating,
		detail.TeachingMaterial.Statistics.RatingCount,
		detail.TeachingMaterial.Tags,
		detail.TeachingMaterial.CurriculumAlignment.Objectives,
	)

	return &types.ToolsCallResponse{
		Content: []types.Content{
			{
				Type: "text",
				Text: response,
			},
		},
		IsError: false,
	}, nil
}

func (s *MCPService) handleGetRelatedMaterials(ctx *types.ToolContext, args interface{}) (*types.ToolsCallResponse, error) {
	return &types.ToolsCallResponse{
		Content: []types.Content{
			{
				Type: "text",
				Text: "关联素材推荐功能正在开发中...",
			},
		},
		IsError: false,
	}, nil
}

func (s *MCPService) handleGenerateLessonPlan(ctx *types.ToolContext, args interface{}) (*types.ToolsCallResponse, error) {
	// 解析参数
	var params struct {
		MaterialIDs []string `json:"material_ids"`
		Objectives  []string `json:"objectives"`
		Grade       string   `json:"grade"`
		StudentLevel string  `json:"student_level"`
		Duration    int      `json:"duration"`
	}

	argsData, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	if err := json.Unmarshal(argsData, &params); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	if len(params.MaterialIDs) == 0 || len(params.Objectives) == 0 {
		return nil, fmt.Errorf("material_ids and objectives are required")
	}

	// 解析素材ID列表
	var materialIDs []uuid.UUID
	for _, idStr := range params.MaterialIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid material_id format: %s", idStr)
		}
		materialIDs = append(materialIDs, id)
	}

	// 基于学而思教研标准生成教案
	lessonPlan := s.generateLessonPlan(materialIDs, params.Objectives, params.Grade, params.StudentLevel, params.Duration)

	return &types.ToolsCallResponse{
		Content: []types.Content{
			{
				Type: "text",
				Text: lessonPlan,
			},
		},
		IsError: false,
	}, nil
}

func (s *MCPService) handleGenerateExercises(ctx *types.ToolContext, args interface{}) (*types.ToolsCallResponse, error) {
	return &types.ToolsCallResponse{
		Content: []types.Content{
			{
				Type: "text",
				Text: "练习题生成功能正在开发中...",
			},
		},
		IsError: false,
	}, nil
}

// 资源处理器实现
func (s *MCPService) handleCurriculumResource(uri string) (*types.ResourcesReadResponse, error) {
	// 模拟课程大纲数据
	curriculumData := map[string]interface{}{
		"grade": "grade_1",
		"subject": "math",
		"objectives": []string{
			"认识数字1-20",
			"掌握基本加减法",
			"理解几何图形",
		},
		"units": []string{
			"数字与运算",
			"图形与测量",
			"统计与概率初步",
		},
	}

	data, _ := json.Marshal(curriculumData)

	return &types.ResourcesReadResponse{
		Contents: []types.ResourceContent{
			{
				URI:      uri,
				MimeType: "application/json",
				Text:     string(data),
			},
		},
	}, nil
}

func (s *MCPService) handleKnowledgeGraphResource(uri string) (*types.ResourcesReadResponse, error) {
	// 模拟知识图谱数据
	graphData := map[string]interface{}{
		"nodes": []map[string]interface{}{
			{"id": "addition", "label": "加法", "level": "basic"},
			{"id": "subtraction", "label": "减法", "level": "basic"},
			{"id": "multiplication", "label": "乘法", "level": "intermediate"},
		},
		"edges": []map[string]interface{}{
			{"source": "addition", "target": "multiplication", "relation": "prerequisite"},
		},
	}

	data, _ := json.Marshal(graphData)

	return &types.ResourcesReadResponse{
		Contents: []types.ResourceContent{
			{
				URI:      uri,
				MimeType: "application/json",
				Text:     string(data),
			},
		},
	}, nil
}

func (s *MCPService) handleTeachingTemplateResource(uri string) (*types.ResourcesReadResponse, error) {
	// 模拟教学模板数据
	templateData := map[string]interface{}{
		"model": "5E",
		"phases": []map[string]interface{}{
			{"phase": "Engage", "description": "激发兴趣", "duration": 5},
			{"phase": "Explore", "description": "自主探索", "duration": 15},
			{"phase": "Explain", "description": "概念讲解", "duration": 10},
			{"phase": "Elaborate", "description": "深化应用", "duration": 10},
			{"phase": "Evaluate", "description": "效果评估", "duration": 5},
		},
	}

	data, _ := json.Marshal(templateData)

	return &types.ResourcesReadResponse{
		Contents: []types.ResourceContent{
			{
				URI:      uri,
				MimeType: "application/json",
				Text:     string(data),
			},
		},
	}, nil
}

// 辅助方法
func (s *MCPService) parseParams(params interface{}, target interface{}) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func (s *MCPService) createSuccessResponse(id interface{}, result interface{}) (*types.MCPResponse, error) {
	return &types.MCPResponse{
		MCPMessage: types.MCPMessage{
			JSONRPC: "2.0",
			ID:      id,
		},
		Result: result,
	}, nil
}

func (s *MCPService) createErrorResponse(id interface{}, code int, message string) (*types.MCPResponse, error) {
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

func (s *MCPService) monitorResourceUpdates(uri string, ch chan *types.MCPNotification) {
	// 简化实现：定时发送更新通知
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	defer close(ch)

	for {
		select {
		case <-ticker.C:
			// 发送资源更新通知
			notification := &types.MCPNotification{
				MCPMessage: types.MCPMessage{
					JSONRPC: "2.0",
				},
				Method: types.MCPMethodResourcesUpdated,
				Params: types.ResourcesUpdatedParams{
					URI: uri,
				},
			}
			select {
			case ch <- notification:
			default:
				return // 通道已满，退出
			}
		}
	}
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

func getClientIDFromContext(ctx context.Context) string {
	// TODO: 从上下文获取客户端ID
	return "default-client"
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

func convertToGradeLevels(grades []string) []types.GradeLevel {
	result := make([]types.GradeLevel, len(grades))
	for i, grade := range grades {
		result[i] = types.GradeLevel(grade)
	}
	return result
}

// formatDuration 格式化时长
func formatDuration(duration *int) string {
	if duration == nil {
		return "未知"
	}
	minutes := *duration / 60
	if minutes < 60 {
		return fmt.Sprintf("%d分钟", minutes)
	}
	hours := minutes / 60
	remainingMinutes := minutes % 60
	return fmt.Sprintf("%d小时%d分钟", hours, remainingMinutes)
}

// generateLessonPlan 基于学而思教研标准生成教案
func (s *MCPService) generateLessonPlan(materialIDs []uuid.UUID, objectives []string, grade, studentLevel string, duration int) string {
	// 获取素材信息（简化实现，这里只用第一个素材）
	var mainMaterial *types.TeachingMaterial
	if len(materialIDs) > 0 {
		if material, err := s.config.MaterialService.GetMaterialDetail(uuid.Nil, materialIDs[0]); err == nil {
			mainMaterial = material.TeachingMaterial
		}
	}

	// 基于学而思5E教学模型生成教案
	lessonPlan := fmt.Sprintf(`# 学而思教研标准教案

## 基本信息
- **年级**: %s
- **学生水平**: %s
- **总时长**: %d分钟
- **教学模式**: 线上直播 + 互动练习

## 教学目标
`, grade, studentLevel, duration)

	for i, obj := range objectives {
		lessonPlan += fmt.Sprintf("%d. %s\n", i+1, obj)
	}

	lessonPlan += `
## 教学流程 (5E模型)

### 1. Engage (激发兴趣) - 5分钟
- **教学活动**: 通过实际问题引入概念
- **互动方式**: 学生分享相关生活经验
- **预期效果**: 激发学习兴趣，建立情感连接

### 2. Explore (自主探索) - 15分钟
- **教学活动**: 学生自主探索和发现规律
- **材料使用**: `

	if mainMaterial != nil {
		lessonPlan += fmt.Sprintf("%s (%s)", mainMaterial.Title, mainMaterial.Type)
	} else {
		lessonPlan += "教学演示视频"
	}

	lessonPlan += `
- **分组活动**: 小组讨论和初步尝试
- **教师指导**: 巡视指导，及时干预

### 3. Explain (概念讲解) - 10分钟
- **教学重点**: 清晰阐述核心概念和原理
- **板书设计**: 重点概念和公式突出显示
- **举例说明**: 结合实际生活案例
- **常见问题**: 解答学生疑惑

### 4. Elaborate (深化应用) - 10分钟
- **练习巩固**: 分层练习，针对不同水平
- **应用拓展**: 实际问题解决
- **思维训练**: 培养解题思维和方法

### 5. Evaluate (效果评估) - 5分钟
- ** formative 评价**: 过程性评价
- ** summative 评价**: 学习效果检查
- **反馈收集**: 学生学习感受和建议

## 教学资源
- **主教材**: 学而思同步教材
- **辅助材料**: 互动练习册
- **技术工具**: 智能白板 + 在线练习平台

## 教学评价标准
- **知识掌握**: 能够正确理解和应用概念
- **技能培养**: 具备基本的解题能力和思维方法
- **情感态度**: 对学习产生兴趣，形成良好学习习惯

## 注意事项
- 根据学生实际反馈调整教学节奏
- 注重个别化指导，关注学困生
- 做好教学反思，为下次教学改进提供依据

---
*本教案基于学而思教研标准自动生成，可根据实际教学情况进行调整*`

	return lessonPlan
}
