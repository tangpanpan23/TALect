package service

import (
	"context"
	"time"

	"github.com/future-mcp/future-mcp-server/internal/types"
	"github.com/google/uuid"
)

// MaterialService 素材服务接口
type MaterialService interface {
	// 搜索相关
	SearchMaterials(userID uuid.UUID, req types.SearchMaterialsRequest) (*types.SearchMaterialsResponse, error)
	SearchByGradeSubject(userID uuid.UUID, grade types.GradeLevel, subject types.Subject, difficulty types.Difficulty, teachingStage string) (*types.SearchMaterialsResponse, error)
	SemanticSearch(userID uuid.UUID, query string, limit int) (*types.SearchMaterialsResponse, error)

	// 详情相关
	GetMaterialDetail(userID uuid.UUID, materialID uuid.UUID) (*types.MaterialDetailResponse, error)
	GetRelatedMaterials(userID uuid.UUID, materialID uuid.UUID, relationType string, limit int) (*types.SearchMaterialsResponse, error)

	// 分析相关
	AnalyzeMaterial(userID uuid.UUID, req types.MaterialAnalysisRequest) (*types.MaterialAnalysisResponse, error)

	// 推荐相关
	GetPersonalizedRecommendations(userID uuid.UUID, limit int) (*types.RecommendationResult, error)
}

// ToolService 工具服务接口
type ToolService interface {
	// 工具执行相关
	ExecuteTool(ctx context.Context, toolName string, params map[string]interface{}) (interface{}, error)

	// 工具管理相关
	ListAvailableTools() ([]types.Tool, error)
	GetToolDefinition(toolName string) (*types.ToolDefinition, error)

	// 工具统计相关
	GetToolUsageStatistics() (map[string]interface{}, error)
}

// ResourceService 资源服务接口
type ResourceService interface {
	// 资源读取相关
	GetResource(uri string) (*types.ResourcesReadResponse, error)
	ListResources() (*types.ResourcesListResponse, error)

	// 资源订阅相关
	SubscribeResource(ctx context.Context, uri string) (<-chan *types.MCPNotification, error)
	UnsubscribeResource(ctx context.Context, uri string) error

	// 资源管理相关
	CreateResource(resource *types.ResourceDefinition) error
	UpdateResource(uri string, resource *types.ResourceDefinition) error
	DeleteResource(uri string) error
}

// UserService 用户服务接口
type UserService interface {
	// 用户认证相关
	AuthenticateUser(ctx context.Context, token string) (*types.User, error)
	ValidateUserPermission(userID uuid.UUID, resource string, action string) (bool, error)

	// 用户信息相关
	GetUserProfile(userID uuid.UUID) (*types.UserProfileResponse, error)
	UpdateUserProfile(userID uuid.UUID, updates *types.UpdateUserProfileRequest) error

	// 用户配额相关
	GetUserQuota(userID uuid.UUID) (*types.UserQuotaResponse, error)
	CheckUserQuota(userID uuid.UUID, action string) (bool, error)
	UpdateUserQuota(userID uuid.UUID, action string) error

	// 用户活动相关
	LogUserActivity(userID uuid.UUID, activity *types.UserActivity) error
	GetUserActivityHistory(userID uuid.UUID, limit int) ([]types.UserActivity, error)
}

// AuthService 认证服务接口
type AuthService interface {
	// JWT相关
	GenerateToken(userID uuid.UUID, claims map[string]interface{}) (string, error)
	ValidateToken(token string) (*types.User, error)
	RefreshToken(refreshToken string) (string, error)

	// API Key相关
	ValidateAPIKey(apiKey string) (*types.User, error)
	GenerateAPIKey(userID uuid.UUID, name string) (*types.APIKey, error)
	ListUserAPIKeys(userID uuid.UUID) ([]types.APIKey, error)
	RevokeAPIKey(userID uuid.UUID, keyID uuid.UUID) error

	// 权限检查相关
	CheckPermission(userID uuid.UUID, resource string, action string) (bool, error)
	GetUserRoles(userID uuid.UUID) ([]string, error)
}

// CacheService 缓存服务接口
type CacheService interface {
	// 基础缓存操作
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl int) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool

	// 素材缓存
	GetMaterialCache(materialID string) (*types.TeachingMaterial, error)
	SetMaterialCache(material *types.TeachingMaterial, ttl int) error
	DeleteMaterialCache(materialID string) error

	// 搜索缓存
	GetSearchCache(query string, filters map[string]interface{}) (*types.SearchResult, error)
	SetSearchCache(query string, filters map[string]interface{}, result *types.SearchResult, ttl int) error

	// 用户缓存
	GetUserCache(userID uuid.UUID) (*types.User, error)
	SetUserCache(user *types.User, ttl int) error

	// JSON操作
	SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

// VectorService 向量服务接口
type VectorService interface {
	// 向量操作
	EmbedText(text string) ([]float32, error)
	SearchSimilar(queryVector []float32, limit int, threshold float64) ([]types.VectorSearchResult, error)

	// 素材向量管理
	AddMaterialVector(materialID uuid.UUID, vector []float32, metadata map[string]interface{}) error
	UpdateMaterialVector(materialID uuid.UUID, vector []float32, metadata map[string]interface{}) error
	DeleteMaterialVector(materialID uuid.UUID) error

	// 批量操作
	BatchAddVectors(vectors []types.VectorRecord) error
	BatchSearchSimilar(queryVectors [][]float32, limit int) ([][]types.VectorSearchResult, error)
}

// VectorSearchResult 向量搜索结果
type VectorSearchResult struct {
	ID       string                 `json:"id"`
	Score    float64                `json:"score"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// VectorRecord 向量记录
type VectorRecord struct {
	ID       string                 `json:"id"`
	Vector   []float32              `json:"vector"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
