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

// MemoryCacheService 内存缓存服务实现（用于测试）
type MemoryCacheService struct {
	data map[string]interface{}
	ttl  map[string]time.Time
	mu   sync.RWMutex
}

// NewMemoryCacheService 创建内存缓存服务
func NewMemoryCacheService() CacheService {
	cache := &MemoryCacheService{
		data: make(map[string]interface{}),
		ttl:  make(map[string]time.Time),
	}

	// 启动清理过期缓存的goroutine
	go cache.cleanupExpired()

	return cache
}

// Get 获取缓存值
func (c *MemoryCacheService) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if expiry, exists := c.ttl[key]; exists && time.Now().After(expiry) {
		// 缓存已过期
		delete(c.data, key)
		delete(c.ttl, key)
		return "", fmt.Errorf("cache miss")
	}

	if value, exists := c.data[key]; exists {
		if str, ok := value.(string); ok {
			return str, nil
		}
		return fmt.Sprintf("%v", value), nil
	}

	return "", fmt.Errorf("cache miss")
}

// Set 设置缓存值
func (c *MemoryCacheService) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	c.ttl[key] = time.Now().Add(time.Duration(ttl) * time.Second)

	return nil
}

// Delete 删除缓存
func (c *MemoryCacheService) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
	delete(c.ttl, key)

	return nil
}

// Exists 检查键是否存在
func (c *MemoryCacheService) Exists(ctx context.Context, key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if expiry, exists := c.ttl[key]; exists && time.Now().After(expiry) {
		// 缓存已过期，清理
		delete(c.data, key)
		delete(c.ttl, key)
		return false
	}

	_, exists := c.data[key]
	return exists
}

// GetMaterialCache 获取素材缓存
func (c *MemoryCacheService) GetMaterialCache(materialID string) (*types.TeachingMaterial, error) {
	key := fmt.Sprintf("material:%s", materialID)
	value, err := c.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}

	var material types.TeachingMaterial
	if err := json.Unmarshal([]byte(value), &material); err != nil {
		return nil, fmt.Errorf("failed to unmarshal material: %w", err)
	}

	return &material, nil
}

// SetMaterialCache 设置素材缓存
func (c *MemoryCacheService) SetMaterialCache(material *types.TeachingMaterial, ttl int) error {
	key := fmt.Sprintf("material:%s", material.ID)
	return c.SetJSON(context.Background(), key, material, time.Duration(ttl)*time.Second)
}

// GetSearchCache 获取搜索缓存
func (c *MemoryCacheService) GetSearchCache(query string, filters map[string]interface{}) (*types.SearchResult, error) {
	key := fmt.Sprintf("search:%s", query)
	value, err := c.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}

	var result types.SearchResult
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search result: %w", err)
	}

	return &result, nil
}

// SetSearchCache 设置搜索缓存
func (c *MemoryCacheService) SetSearchCache(query string, filters map[string]interface{}, result *types.SearchResult, ttl int) error {
	key := fmt.Sprintf("search:%s", query)
	return c.SetJSON(context.Background(), key, result, time.Duration(ttl)*time.Second)
}

// SetJSON 设置JSON格式的缓存值
func (c *MemoryCacheService) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.Set(ctx, key, string(data), int(ttl.Seconds()))
}

// DeleteMaterialCache 删除素材缓存
func (c *MemoryCacheService) DeleteMaterialCache(materialID string) error {
	key := fmt.Sprintf("material:%s", materialID)
	return c.Delete(context.Background(), key)
}

// GetUserCache 获取用户缓存
func (c *MemoryCacheService) GetUserCache(userID uuid.UUID) (*types.User, error) {
	key := fmt.Sprintf("user:%s", userID)
	value, err := c.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}

	var user types.User
	if err := json.Unmarshal([]byte(value), &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// SetUserCache 设置用户缓存
func (c *MemoryCacheService) SetUserCache(user *types.User, ttl int) error {
	key := fmt.Sprintf("user:%s", user.ID)
	return c.SetJSON(context.Background(), key, user, time.Duration(ttl)*time.Second)
}

// cleanupExpired 清理过期缓存
func (c *MemoryCacheService) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()

		for key, expiry := range c.ttl {
			if now.After(expiry) {
				delete(c.data, key)
				delete(c.ttl, key)
				logger.Info("Cleaned up expired cache", logger.Any("key", key))
			}
		}

		c.mu.Unlock()
	}
}
