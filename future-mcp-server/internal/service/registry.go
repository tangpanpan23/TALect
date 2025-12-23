package service

import (
	"sync"

	"github.com/future-mcp/future-mcp-server/internal/types"
)

// ToolRegistry 工具注册器
type ToolRegistry struct {
	tools map[string]*types.ToolDefinition
	mu    sync.RWMutex
}

// NewToolRegistry 创建工具注册器
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]*types.ToolDefinition),
	}
}

// RegisterTool 注册工具
func (tr *ToolRegistry) RegisterTool(tool *types.ToolDefinition) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.tools[tool.Name] = tool
}

// GetTool 获取工具
func (tr *ToolRegistry) GetTool(name string) *types.ToolDefinition {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	return tr.tools[name]
}

// ListTools 列出所有工具
func (tr *ToolRegistry) ListTools() map[string]*types.ToolDefinition {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	tools := make(map[string]*types.ToolDefinition)
	for name, tool := range tr.tools {
		tools[name] = tool
	}
	return tools
}

// RemoveTool 移除工具
func (tr *ToolRegistry) RemoveTool(name string) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	delete(tr.tools, name)
}

// ResourceRegistry 资源注册器
type ResourceRegistry struct {
	resources map[string]*types.ResourceDefinition
	mu        sync.RWMutex
}

// NewResourceRegistry 创建资源注册器
func NewResourceRegistry() *ResourceRegistry {
	return &ResourceRegistry{
		resources: make(map[string]*types.ResourceDefinition),
	}
}

// RegisterResource 注册资源
func (rr *ResourceRegistry) RegisterResource(resource *types.ResourceDefinition) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	rr.resources[resource.URI] = resource
}

// GetResource 获取资源
func (rr *ResourceRegistry) GetResource(uri string) *types.ResourceDefinition {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	return rr.resources[uri]
}

// ListResources 列出所有资源
func (rr *ResourceRegistry) ListResources() map[string]*types.ResourceDefinition {
	rr.mu.RLock()
	defer rr.mu.RUnlock()

	resources := make(map[string]*types.ResourceDefinition)
	for uri, resource := range rr.resources {
		resources[uri] = resource
	}
	return resources
}

// RemoveResource 移除资源
func (rr *ResourceRegistry) RemoveResource(uri string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	delete(rr.resources, uri)
}

// SubscriptionManager 订阅管理器
type SubscriptionManager struct {
	subscriptions map[string]map[string]chan *types.MCPNotification // uri -> clientID -> channel
	mu           sync.RWMutex
}

// NewSubscriptionManager 创建订阅管理器
func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions: make(map[string]map[string]chan *types.MCPNotification),
	}
}

// Subscribe 订阅资源
func (sm *SubscriptionManager) Subscribe(uri, clientID string) chan *types.MCPNotification {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.subscriptions[uri] == nil {
		sm.subscriptions[uri] = make(map[string]chan *types.MCPNotification)
	}

	ch := make(chan *types.MCPNotification, 10) // 缓冲通道
	sm.subscriptions[uri][clientID] = ch

	return ch
}

// Unsubscribe 取消订阅资源
func (sm *SubscriptionManager) Unsubscribe(uri, clientID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if clients, exists := sm.subscriptions[uri]; exists {
		if ch, exists := clients[clientID]; exists {
			close(ch)
			delete(clients, clientID)
		}

		if len(clients) == 0 {
			delete(sm.subscriptions, uri)
		}
	}
}

// NotifyResourceUpdate 通知资源更新
func (sm *SubscriptionManager) NotifyResourceUpdate(uri string, notification *types.MCPNotification) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if clients, exists := sm.subscriptions[uri]; exists {
		for _, ch := range clients {
			select {
			case ch <- notification:
			default:
				// 通道已满，跳过
			}
		}
	}
}

// GetSubscribedClients 获取订阅的客户端
func (sm *SubscriptionManager) GetSubscribedClients(uri string) []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var clients []string
	if clientMap, exists := sm.subscriptions[uri]; exists {
		for clientID := range clientMap {
			clients = append(clients, clientID)
		}
	}
	return clients
}

// Cleanup 清理资源
func (sm *SubscriptionManager) Cleanup() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for uri, clients := range sm.subscriptions {
		for _, ch := range clients {
			close(ch)
		}
		delete(sm.subscriptions, uri)
	}
}
