package types

import (
	"time"

	"github.com/google/uuid"
)

// MCP协议相关类型定义
// 参考: https://modelcontextprotocol.io/specification

// MCPProtocolVersion 协议版本
const MCPProtocolVersion = "2024-11-05"

// MCPMessage MCP消息基础结构
type MCPMessage struct {
	JSONRPC string `json:"jsonrpc" binding:"required,eq=2.0"`
	ID      interface{} `json:"id,omitempty"` // string | number | null
}

// MCPRequest MCP请求
type MCPRequest struct {
	MCPMessage
	Method string `json:"method" binding:"required"`
	Params interface{} `json:"params,omitempty"`
}

// MCPResponse MCP响应
type MCPResponse struct {
	MCPMessage
	Result interface{} `json:"result,omitempty"`
	Error  *MCPError   `json:"error,omitempty"`
}

// MCPNotification MCP通知
type MCPNotification struct {
	MCPMessage
	Method string `json:"method" binding:"required"`
	Params interface{} `json:"params,omitempty"`
}

// MCPError MCP错误
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP错误码
const (
	MCPParseError     = -32700
	MCPInvalidRequest = -32600
	MCPMethodNotFound = -32601
	MCPInvalidParams  = -32602
	MCPInternalError  = -32603
)

// MCP标准方法
const (
	MCPMethodInitialize     = "initialize"
	MCPMethodInitialized    = "initialized"
	MCPMethodToolsList      = "tools/list"
	MCPMethodToolsCall      = "tools/call"
	MCPMethodResourcesList  = "resources/list"
	MCPMethodResourcesRead  = "resources/read"
	MCPMethodResourcesSubscribe = "resources/subscribe"
	MCPMethodResourcesUnsubscribe = "resources/unsubscribe"
	MCPMethodPing           = "ping"
	MCPMethodCancel         = "cancel"
	MCPMethodProgress       = "notifications/progress"
	MCPMethodResourcesUpdated = "notifications/resources/updated"
	MCPMethodToolsChanged   = "notifications/tools/list_changed"
)

// ==================== 初始化相关 ====================

// InitializeRequest 初始化请求
type InitializeRequest struct {
	ProtocolVersion string                 `json:"protocolVersion" binding:"required"`
	Capabilities    ClientCapabilities     `json:"capabilities" binding:"required"`
	ClientInfo      ImplementationInfo     `json:"clientInfo" binding:"required"`
}

// ClientCapabilities 客户端能力
type ClientCapabilities struct {
	Tools     *ToolsCapability     `json:"tools,omitempty"`
	Resources *ResourcesCapability `json:"resources,omitempty"`
	Sampling  *SamplingCapability  `json:"sampling,omitempty"`
	Logging   *LoggingCapability   `json:"logging,omitempty"`
}

// ToolsCapability 工具能力
type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ResourcesCapability 资源能力
type ResourcesCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
	Subscribe   bool `json:"subscribe,omitempty"`
}

// SamplingCapability 采样能力
type SamplingCapability struct {
}

// LoggingCapability 日志能力
type LoggingCapability struct {
}

// ImplementationInfo 实现信息
type ImplementationInfo struct {
	Name    string `json:"name" binding:"required"`
	Version string `json:"version" binding:"required"`
}

// InitializeResponse 初始化响应
type InitializeResponse struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    ServerCapabilities     `json:"capabilities"`
	ServerInfo      ImplementationInfo     `json:"serverInfo"`
	Instructions    string                 `json:"instructions,omitempty"`
}

// ServerCapabilities 服务器能力
type ServerCapabilities struct {
	Tools     *ServerToolsCapability     `json:"tools,omitempty"`
	Resources *ServerResourcesCapability `json:"resources,omitempty"`
	Logging   *ServerLoggingCapability   `json:"logging,omitempty"`
	Prompts   *ServerPromptsCapability   `json:"prompts,omitempty"`
}

// ServerToolsCapability 服务器工具能力
type ServerToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ServerResourcesCapability 服务器资源能力
type ServerResourcesCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
	Subscribe   bool `json:"subscribe,omitempty"`
}

// ServerLoggingCapability 服务器日志能力
type ServerLoggingCapability struct {
}

// ServerPromptsCapability 服务器提示能力
type ServerPromptsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ==================== 工具相关 ====================

// ToolsListRequest 工具列表请求
type ToolsListRequest struct{}

// ToolsListResponse 工具列表响应
type ToolsListResponse struct {
	Tools []Tool `json:"tools"`
}

// Tool 工具定义
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	InputSchema interface{} `json:"inputSchema"` // JSON Schema
}

// ToolsCallRequest 工具调用请求
type ToolsCallRequest struct {
	Name      string      `json:"name" binding:"required"`
	Arguments interface{} `json:"arguments,omitempty"`
}

// ToolsCallResponse 工具调用响应
type ToolsCallResponse struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

// Content 内容
type Content struct {
	Type string `json:"type" binding:"required,oneof=text image"` // text | image
	Text string `json:"text,omitempty"`
	Data string `json:"data,omitempty"` // base64 for images
	MimeType string `json:"mimeType,omitempty"`
}

// ==================== 资源相关 ====================

// ResourcesListRequest 资源列表请求
type ResourcesListRequest struct{}

// ResourcesListResponse 资源列表响应
type ResourcesListResponse struct {
	Resources []Resource `json:"resources"`
}

// Resource 资源定义
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// ResourcesReadRequest 资源读取请求
type ResourcesReadRequest struct {
	URI string `json:"uri" binding:"required"`
}

// ResourcesReadResponse 资源读取响应
type ResourcesReadResponse struct {
	Contents []ResourceContent `json:"contents"`
}

// ResourceContent 资源内容
type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType,omitempty"`
	Text     string `json:"text,omitempty"`
	Blob     string `json:"blob,omitempty"` // base64 encoded
}

// ResourcesSubscribeRequest 资源订阅请求
type ResourcesSubscribeRequest struct {
	URI string `json:"uri" binding:"required"`
}

// ResourcesUnsubscribeRequest 资源取消订阅请求
type ResourcesUnsubscribeRequest struct {
	URI string `json:"uri" binding:"required"`
}

// ==================== 通知相关 ====================

// ProgressNotification 进度通知
type ProgressNotification struct {
	Method string           `json:"method" binding:"eq=notifications/progress"`
	Params ProgressParams   `json:"params"`
}

// ProgressParams 进度参数
type ProgressParams struct {
	ProgressToken string  `json:"progressToken"`
	Progress      float64 `json:"progress"`      // 0-1
	Total         *float64 `json:"total,omitempty"`
	Message       string   `json:"message,omitempty"`
}

// ResourcesUpdatedNotification 资源更新通知
type ResourcesUpdatedNotification struct {
	Method string                    `json:"method" binding:"eq=notifications/resources/updated"`
	Params ResourcesUpdatedParams    `json:"params"`
}

// ResourcesUpdatedParams 资源更新参数
type ResourcesUpdatedParams struct {
	URI string `json:"uri"`
}

// ToolsListChangedNotification 工具列表变更通知
type ToolsListChangedNotification struct {
	Method string `json:"method" binding:"eq=notifications/tools/list_changed"`
}

// ==================== 工具实现相关 ====================

// ToolDefinition 工具定义（内部使用）
type ToolDefinition struct {
	Name        string
	Description string
	Handler     ToolHandler
	InputSchema interface{}
}

// ToolHandler 工具处理器
type ToolHandler func(ctx *ToolContext, args interface{}) (*ToolsCallResponse, error)

// ToolContext 工具上下文
type ToolContext struct {
	UserID      uuid.UUID
	SessionID   string
	RequestID   string
	StartTime   time.Time
	Parameters  map[string]interface{}
}

// ResourceDefinition 资源定义（内部使用）
type ResourceDefinition struct {
	URI         string
	Name        string
	Description string
	MimeType    string
	Handler     ResourceHandler
}

// ResourceHandler 资源处理器
type ResourceHandler func(uri string) (*ResourcesReadResponse, error)

// ==================== 扩展类型 ====================

// PaginatedRequest 分页请求（扩展）
type PaginatedRequest struct {
	Cursor *string `json:"cursor,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
}

// CursorResponse 游标响应（扩展）
type CursorResponse struct {
	Data       interface{} `json:"data"`
	NextCursor *string     `json:"nextCursor,omitempty"`
	HasMore    bool        `json:"hasMore"`
}

// SearchFilters 搜索过滤器
type SearchFilters struct {
	Query         string       `json:"query,omitempty"`
	GradeLevels   []GradeLevel `json:"grade_levels,omitempty"`
	Subjects      []Subject    `json:"subjects,omitempty"`
	MaterialTypes []MaterialType `json:"material_types,omitempty"`
	Difficulties  []Difficulty `json:"difficulties,omitempty"`
	Tags          []string     `json:"tags,omitempty"`
	DateRange     *DateRange   `json:"date_range,omitempty"`
}

// DateRange 日期范围
type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Materials  []TeachingMaterial `json:"materials"`
	TotalCount int64              `json:"total_count"`
	Cursor     *string            `json:"cursor,omitempty"`
	HasMore    bool               `json:"has_more"`
	SearchTime float64            `json:"search_time"`
}

// RecommendationResult 推荐结果
type RecommendationResult struct {
	Materials     []TeachingMaterial `json:"materials"`
	Algorithm     string             `json:"algorithm"`
	Confidence    float64            `json:"confidence"`
	Reason        string             `json:"reason"`
}
