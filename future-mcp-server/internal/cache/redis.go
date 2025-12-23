package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/future-mcp/future-mcp-server/pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// RedisClient Redis客户端
var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() (*redis.Client, error) {
	host := viper.GetString("redis.host")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")
	poolSize := viper.GetInt("redis.pool_size")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         host,
		Password:     password,
		DB:           db,
		PoolSize:     poolSize,
		MinIdleConns: viper.GetInt("redis.min_idle_conns"),
		ConnMaxLifetime: time.Duration(viper.GetInt("redis.conn_max_lifetime")) * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis connected successfully",
		logger.Any("host", host),
		logger.Any("db", db),
		logger.Any("pool_size", poolSize),
	)

	return RedisClient, nil
}

// GetRedis 获取Redis客户端
func GetRedis() *redis.Client {
	return RedisClient
}

// Close 关闭Redis连接
func Close() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// HealthCheck Redis健康检查
func HealthCheck() error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return RedisClient.Ping(ctx).Err()
}

// Cache 缓存接口
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建Redis缓存
func NewRedisCache(client *redis.Client) Cache {
	return &RedisCache{client: client}
}

// Get 获取缓存值
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set 设置缓存值
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var val string
	switch v := value.(type) {
	case string:
		val = v
	case int, int64, float64, bool:
		val = fmt.Sprintf("%v", v)
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		val = string(data)
	}

	return c.client.Set(ctx, key, val, expiration).Err()
}

// Delete 删除缓存
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func (c *RedisCache) Exists(ctx context.Context, key string) bool {
	count, err := c.client.Exists(ctx, key).Result()
	return err == nil && count > 0
}

// Expire 设置过期时间
func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// GetJSON 获取JSON格式的缓存值
func (c *RedisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// SetJSON 设置JSON格式的缓存值
func (c *RedisCache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.Set(ctx, key, string(data), ttl)
}

// CacheKey 缓存键生成器
type CacheKey struct {
	prefix string
}

// NewCacheKey 创建缓存键生成器
func NewCacheKey(prefix string) *CacheKey {
	return &CacheKey{prefix: prefix}
}

// Key 生成缓存键
func (ck *CacheKey) Key(parts ...interface{}) string {
	key := ck.prefix
	for _, part := range parts {
		key += fmt.Sprintf(":%v", part)
	}
	return key
}

// MaterialCacheKeys 素材缓存键
var MaterialCacheKeys = NewCacheKey("material")

// UserCacheKeys 用户缓存键
var UserCacheKeys = NewCacheKey("user")

// SearchCacheKeys 搜索缓存键
var SearchCacheKeys = NewCacheKey("search")

// CacheConfig 缓存配置
type CacheConfig struct {
	DefaultTTL time.Duration
	MaterialTTL time.Duration
	UserTTL     time.Duration
	SearchTTL   time.Duration
}

// GetCacheConfig 获取缓存配置
func GetCacheConfig() CacheConfig {
	return CacheConfig{
		DefaultTTL:  time.Duration(viper.GetInt("cache.default_ttl")) * time.Second,
		MaterialTTL: time.Duration(viper.GetInt("cache.material_ttl")) * time.Second,
		UserTTL:     time.Duration(viper.GetInt("cache.user_ttl")) * time.Second,
		SearchTTL:   time.Duration(viper.GetInt("cache.search_ttl")) * time.Second,
	}
}

// MaterialCache 素材缓存管理器
type MaterialCache struct {
	cache Cache
	config CacheConfig
}

// NewMaterialCache 创建素材缓存管理器
func NewMaterialCache(cache Cache) *MaterialCache {
	return &MaterialCache{
		cache:  cache,
		config: GetCacheConfig(),
	}
}

// GetMaterial 获取素材缓存
func (mc *MaterialCache) GetMaterial(ctx context.Context, materialID string) (string, error) {
	key := MaterialCacheKeys.Key("detail", materialID)
	return mc.cache.Get(ctx, key)
}

// SetMaterial 设置素材缓存
func (mc *MaterialCache) SetMaterial(ctx context.Context, materialID string, data interface{}) error {
	key := MaterialCacheKeys.Key("detail", materialID)
	return mc.cache.SetJSON(ctx, key, data, mc.config.MaterialTTL)
}

// DeleteMaterial 删除素材缓存
func (mc *MaterialCache) DeleteMaterial(ctx context.Context, materialID string) error {
	key := MaterialCacheKeys.Key("detail", materialID)
	return mc.cache.Delete(ctx, key)
}

// SearchCache 搜索缓存管理器
type SearchCache struct {
	cache Cache
	config CacheConfig
}

// NewSearchCache 创建搜索缓存管理器
func NewSearchCache(cache Cache) *SearchCache {
	return &SearchCache{
		cache:  cache,
		config: GetCacheConfig(),
	}
}

// GetSearchResult 获取搜索结果缓存
func (sc *SearchCache) GetSearchResult(ctx context.Context, query string, filters map[string]interface{}) (string, error) {
	// 生成基于查询和过滤器的缓存键
	key := SearchCacheKeys.Key("query", query)
	for k, v := range filters {
		key = SearchCacheKeys.Key(key, k, v)
	}
	return sc.cache.Get(ctx, key)
}

// SetSearchResult 设置搜索结果缓存
func (sc *SearchCache) SetSearchResult(ctx context.Context, query string, filters map[string]interface{}, result interface{}) error {
	key := SearchCacheKeys.Key("query", query)
	for k, v := range filters {
		key = SearchCacheKeys.Key(key, k, v)
	}
	return sc.cache.SetJSON(ctx, key, result, sc.config.SearchTTL)
}
