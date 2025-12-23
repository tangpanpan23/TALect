package repository

import (
	"github.com/future-mcp/future-mcp-server/internal/types"
	"github.com/google/uuid"
)

// MaterialRepository 素材仓库接口
type MaterialRepository interface {
	// 基础CRUD
	CreateMaterial(material *types.TeachingMaterial) error
	GetMaterialByID(id uuid.UUID) (*types.TeachingMaterial, error)
	UpdateMaterial(material *types.TeachingMaterial) error
	DeleteMaterial(id uuid.UUID) error

	// 搜索相关
	SearchMaterials(req types.SearchMaterialsRequest) ([]types.TeachingMaterial, int64, error)
	GetRelatedMaterials(materialID uuid.UUID, relationType string, limit int) ([]types.TeachingMaterial, error)

	// 统计相关
	GetPopularMaterials(limit int) ([]types.TeachingMaterial, error)
	GetMaterialsByGrade(grade types.GradeLevel, limit int) ([]types.TeachingMaterial, error)
	GetMaterialsBySubject(subject types.Subject, limit int) ([]types.TeachingMaterial, error)

	// 批量操作
	BatchCreateMaterials(materials []*types.TeachingMaterial) error
	BatchUpdateMaterials(materials []*types.TeachingMaterial) error
}

// UserRepository 用户仓库接口
type UserRepository interface {
	// 基础CRUD
	CreateUser(user *types.User) error
	GetUserByID(id uuid.UUID) (*types.User, error)
	GetUserByEmail(email string) (*types.User, error)
	GetUserByUsername(username string) (*types.User, error)
	UpdateUser(user *types.User) error
	DeleteUser(id uuid.UUID) error

	// 权限相关
	GetUserRoles(userID uuid.UUID) ([]string, error)
	UpdateUserRoles(userID uuid.UUID, roles []string) error

	// 配额相关
	GetUserQuota(userID uuid.UUID) (*types.UserQuota, error)
	UpdateUserQuota(userID uuid.UUID, quota *types.UserQuota) error

	// 统计相关
	GetUserStatistics(userID uuid.UUID) (*types.UserStatistics, error)
	UpdateUserStatistics(userID uuid.UUID, stats *types.UserStatistics) error

	// 活动记录
	LogUserActivity(activity *types.UserActivity) error
	GetUserActivities(userID uuid.UUID, limit int, offset int) ([]types.UserActivity, error)
}

// APIKeyRepository API密钥仓库接口
type APIKeyRepository interface {
	CreateAPIKey(apiKey *types.APIKey) error
	GetAPIKeyByKey(key string) (*types.APIKey, error)
	GetAPIKeysByUserID(userID uuid.UUID) ([]types.APIKey, error)
	UpdateAPIKey(apiKey *types.APIKey) error
	DeleteAPIKey(id uuid.UUID) error
	RevokeAPIKey(id uuid.UUID) error
}

// Repositories 仓库集合
type Repositories struct {
	Material MaterialRepository
	User     UserRepository
	APIKey   APIKeyRepository
}

// NewRepositories 创建仓库集合
func NewRepositories(materialRepo MaterialRepository, userRepo UserRepository, apiKeyRepo APIKeyRepository) *Repositories {
	return &Repositories{
		Material: materialRepo,
		User:     userRepo,
		APIKey:   apiKeyRepo,
	}
}
