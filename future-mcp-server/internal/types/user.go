package types

import (
	"time"

	"github.com/google/uuid"
)

// UserRole 用户角色枚举
type UserRole string

const (
	UserRoleGuest      UserRole = "guest"
	UserRoleStudent    UserRole = "student"
	UserRoleTeacher    UserRole = "teacher"
	UserRoleDeveloper  UserRole = "developer"
	UserRolePartner    UserRole = "partner"
	UserRoleInternal   UserRole = "internal"
	UserRoleAdmin      UserRole = "admin"
)

// UserType 用户类型枚举
type UserType string

const (
	UserTypeIndividual UserType = "individual"
	UserTypeSchool     UserType = "school"
	UserTypeCompany    UserType = "company"
	UserTypeGovernment UserType = "government"
)

// UserStatus 用户状态枚举
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted  UserStatus = "deleted"
)

// User 用户主模型
type User struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email       string     `json:"email" gorm:"uniqueIndex;not null"`
	Username    string     `json:"username" gorm:"uniqueIndex;not null"`
	Type        UserType   `json:"type" gorm:"not null"`
	Role        UserRole   `json:"role" gorm:"not null"`
	Status      UserStatus `json:"status" gorm:"not null;default:'active'"`

	// 基本信息
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
	Avatar      string `json:"avatar"`

	// 联系信息
	Phone       string `json:"phone"`
	Company     string `json:"company"`
	Position    string `json:"position"`

	// 偏好设置
	Preferences UserPreferences `json:"preferences" gorm:"embedded"`

	// 配额和限制
	Quota UserQuota `json:"quota" gorm:"embedded"`

	// 统计信息
	Statistics UserStatistics `json:"statistics" gorm:"embedded"`

	// 认证信息
	PasswordHash string    `json:"-" gorm:"not null"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`

	// 元数据
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy uuid.UUID  `json:"created_by" gorm:"type:uuid"`
	UpdatedBy uuid.UUID  `json:"updated_by" gorm:"type:uuid"`
}

// UserPreferences 用户偏好设置
type UserPreferences struct {
	Language         string       `json:"language" gorm:"default:'zh-CN'"`
	Timezone         string       `json:"timezone" gorm:"default:'Asia/Shanghai'"`
	GradeLevels      []GradeLevel `json:"grade_levels" gorm:"type:text[]"`
	Subjects         []Subject    `json:"subjects" gorm:"type:text[]"`
	Difficulties     []Difficulty `json:"difficulties" gorm:"type:text[]"`
	MaterialTypes    []MaterialType `json:"material_types" gorm:"type:text[]"`

	// 个性化设置
	EnableNotifications bool `json:"enable_notifications" gorm:"default:true"`
	EnableRecommendations bool `json:"enable_recommendations" gorm:"default:true"`
	AutoSaveSearch       bool `json:"auto_save_search" gorm:"default:true"`

	// 显示设置
	PageSize            int  `json:"page_size" gorm:"default:20"`
	SortOrder           string `json:"sort_order" gorm:"default:'desc'"`
	DefaultSearchType   string `json:"default_search_type" gorm:"default:'keyword'"`
}

// UserQuota 用户配额
type UserQuota struct {
	DailyRequestLimit    int `json:"daily_request_limit" gorm:"default:1000"`
	MonthlyRequestLimit  int `json:"monthly_request_limit" gorm:"default:30000"`
	ConcurrentLimit      int `json:"concurrent_limit" gorm:"default:10"`

	// 当前使用量
	DailyRequests        int `json:"daily_requests" gorm:"default:0"`
	MonthlyRequests      int `json:"monthly_requests" gorm:"default:0"`
	CurrentConcurrency   int `json:"current_concurrency" gorm:"default:0"`

	// 重置时间
	DailyResetAt         time.Time `json:"daily_reset_at"`
	MonthlyResetAt       time.Time `json:"monthly_reset_at"`

	// 特殊权限
	UnlimitedAccess      bool `json:"unlimited_access" gorm:"default:false"`
	PrioritySupport      bool `json:"priority_support" gorm:"default:false"`
}

// UserStatistics 用户统计信息
type UserStatistics struct {
	TotalRequests        int64 `json:"total_requests" gorm:"default:0"`
	SuccessfulRequests   int64 `json:"successful_requests" gorm:"default:0"`
	FailedRequests       int64 `json:"failed_requests" gorm:"default:0"`

	// 搜索统计
	TotalSearches        int64 `json:"total_searches" gorm:"default:0"`
	AverageSearchTime    float64 `json:"average_search_time" gorm:"type:decimal(5,2)"`

	// 素材使用统计
	FavoriteMaterials    []uuid.UUID `json:"favorite_materials" gorm:"type:uuid[]"`
	ViewHistory          []MaterialViewRecord `json:"view_history" gorm:"type:jsonb"`

	// 偏好学习统计
	MostViewedGrades     map[string]int `json:"most_viewed_grades" gorm:"type:jsonb"`
	MostViewedSubjects   map[string]int `json:"most_viewed_subjects" gorm:"type:jsonb"`
	MostViewedTypes      map[string]int `json:"most_viewed_types" gorm:"type:jsonb"`

	// 时间统计
	TotalStudyTime       int64 `json:"total_study_time" gorm:"default:0"` // 分钟
	AverageSessionTime   float64 `json:"average_session_time" gorm:"type:decimal(5,2)"`
	LastActivityAt       *time.Time `json:"last_activity_at"`
}

// MaterialViewRecord 素材查看记录
type MaterialViewRecord struct {
	MaterialID    uuid.UUID `json:"material_id"`
	ViewCount     int       `json:"view_count"`
	TotalTime     int       `json:"total_time"` // 秒
	LastViewedAt  time.Time `json:"last_viewed_at"`
	Rating        *int      `json:"rating,omitempty"` // 1-5星
}

// UserProfileResponse 用户资料响应
type UserProfileResponse struct {
	ID          uuid.UUID       `json:"id"`
	Email       string          `json:"email"`
	Username    string          `json:"username"`
	Type        UserType        `json:"type"`
	Role        UserRole        `json:"role"`
	Status      UserStatus      `json:"status"`
	DisplayName string          `json:"display_name"`
	Avatar      string          `json:"avatar"`
	Phone       string          `json:"phone"`
	Company     string          `json:"company"`
	Position    string          `json:"position"`
	Preferences UserPreferences `json:"preferences"`
	Statistics  UserStatistics  `json:"statistics"`
	CreatedAt   time.Time       `json:"created_at"`
}

// UpdateUserProfileRequest 更新用户资料请求
type UpdateUserProfileRequest struct {
	DisplayName *string          `json:"display_name,omitempty"`
	Avatar      *string          `json:"avatar,omitempty"`
	Phone       *string          `json:"phone,omitempty"`
	Company     *string          `json:"company,omitempty"`
	Position    *string          `json:"position,omitempty"`
	Preferences *UserPreferences `json:"preferences,omitempty"`
}

// UserQuotaResponse 用户配额响应
type UserQuotaResponse struct {
	DailyLimit         int       `json:"daily_limit"`
	DailyUsed          int       `json:"daily_used"`
	DailyRemaining     int       `json:"daily_remaining"`
	MonthlyLimit       int       `json:"monthly_limit"`
	MonthlyUsed        int       `json:"monthly_used"`
	MonthlyRemaining   int       `json:"monthly_remaining"`
	ConcurrentLimit    int       `json:"concurrent_limit"`
	CurrentConcurrency int       `json:"current_concurrency"`
	DailyResetAt       time.Time `json:"daily_reset_at"`
	MonthlyResetAt     time.Time `json:"monthly_reset_at"`
	UnlimitedAccess    bool      `json:"unlimited_access"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Email       string     `json:"email" binding:"required,email"`
	Username    string     `json:"username" binding:"required,min=3,max=50"`
	Password    string     `json:"password" binding:"required,min=8"`
	Type        UserType   `json:"type" binding:"required"`
	Role        UserRole   `json:"role" binding:"required"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	DisplayName string     `json:"display_name"`
	Phone       string     `json:"phone"`
	Company     string     `json:"company"`
	Position    string     `json:"position"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 用户登录响应
type LoginResponse struct {
	User         UserProfileResponse `json:"user"`
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	TokenType    string              `json:"token_type"`
	ExpiresIn    int                 `json:"expires_in"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ConfirmResetPasswordRequest 确认重置密码请求
type ConfirmResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserActivity 用户活动记录
type UserActivity struct {
	ID         uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID    `json:"user_id" gorm:"type:uuid;index;not null"`
	Action     string       `json:"action" gorm:"not null"` // view/search/download/share/favorite/rate
	ResourceType string     `json:"resource_type" gorm:"not null"` // material/user/tool/resource
	ResourceID uuid.UUID    `json:"resource_id" gorm:"type:uuid"`
	Details    map[string]interface{} `json:"details" gorm:"type:jsonb"`
	IPAddress  string       `json:"ip_address"`
	UserAgent  string       `json:"user_agent"`
	CreatedAt  time.Time    `json:"created_at" gorm:"autoCreateTime"`
}

// APIKey API密钥模型
type APIKey struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;index;not null"`
	Name        string     `json:"name" gorm:"not null"`
	Key         string     `json:"key" gorm:"uniqueIndex;not null"`
	Permissions []string   `json:"permissions" gorm:"type:text[]"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// CreateAPIKeyRequest 创建API密钥请求
type CreateAPIKeyRequest struct {
	Name        string   `json:"name" binding:"required,min=1,max=100"`
	Permissions []string `json:"permissions" binding:"required,min=1"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// APIKeyResponse API密钥响应
type APIKeyResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Key         string     `json:"key"`
	Permissions []string   `json:"permissions"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
}
