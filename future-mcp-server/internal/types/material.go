package types

import (
	"time"

	"github.com/google/uuid"
)

// MaterialType 素材类型枚举
type MaterialType string

const (
	MaterialTypeVideo      MaterialType = "video"
	MaterialTypePPT        MaterialType = "ppt"
	MaterialTypePDF        MaterialType = "pdf"
	MaterialTypeExercise   MaterialType = "exercise"
	MaterialTypeLessonPlan MaterialType = "lesson_plan"
	MaterialTypeAudio      MaterialType = "audio"
	MaterialTypeImage      MaterialType = "image"
)

// GradeLevel 年级枚举
type GradeLevel string

const (
	GradeLevel1  GradeLevel = "grade_1"
	GradeLevel2  GradeLevel = "grade_2"
	GradeLevel3  GradeLevel = "grade_3"
	GradeLevel4  GradeLevel = "grade_4"
	GradeLevel5  GradeLevel = "grade_5"
	GradeLevel6  GradeLevel = "grade_6"
	GradeLevel7  GradeLevel = "grade_7"
	GradeLevel8  GradeLevel = "grade_8"
	GradeLevel9  GradeLevel = "grade_9"
	GradeLevel10 GradeLevel = "grade_10"
	GradeLevel11 GradeLevel = "grade_11"
	GradeLevel12 GradeLevel = "grade_12"
)

// Subject 学科枚举
type Subject string

const (
	SubjectMath     Subject = "math"
	SubjectChinese  Subject = "chinese"
	SubjectEnglish  Subject = "english"
	SubjectPhysics  Subject = "physics"
	SubjectChemistry Subject = "chemistry"
	SubjectBiology   Subject = "biology"
	SubjectHistory   Subject = "history"
	SubjectGeography Subject = "geography"
	SubjectPolitics  Subject = "politics"
)

// Difficulty 难度枚举
type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

// UsageType 使用类型枚举
type UsageType string

const (
	UsageTypeView       UsageType = "view"
	UsageTypeDownload   UsageType = "download"
	UsageTypeEmbed      UsageType = "embed"
	UsageTypeModify     UsageType = "modify"
	UsageTypeDistribute UsageType = "distribute"
)

// TeachingMaterial 教学素材主模型
type TeachingMaterial struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string      `json:"title" gorm:"not null;index"`
	Description string      `json:"description" gorm:"type:text"`
	Type        MaterialType `json:"type" gorm:"not null"`
	GradeLevels []GradeLevel `json:"grade_levels" gorm:"type:text[]"`
	Subject     Subject     `json:"subject" gorm:"not null;index"`
	Tags        []string    `json:"tags" gorm:"type:text[]"`
	Difficulty  Difficulty  `json:"difficulty" gorm:"not null"`

	// 课标对齐信息
	CurriculumAlignment CurriculumAlignment `json:"curriculum_alignment" gorm:"embedded"`

	// 元数据
	Metadata MaterialMetadata `json:"metadata" gorm:"embedded"`

	// 权限控制
	Permissions MaterialPermissions `json:"permissions" gorm:"embedded"`

	// 统计信息
	Statistics MaterialStatistics `json:"statistics" gorm:"embedded"`

	// 向量嵌入（用于语义搜索）
	Embeddings []float32 `json:"embeddings,omitempty" gorm:"type:vector(768)"`

	// 时间戳
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// CurriculumAlignment 课标对齐信息
type CurriculumAlignment struct {
	Standard        string   `json:"standard" gorm:"column:curriculum_standard"` // 课标编号
	Objectives      []string `json:"objectives" gorm:"type:text[]"`             // 学习目标
	CompetencyLevel int      `json:"competency_level" gorm:"default:0"`         // 对齐度评分 (0-100)
}

// MaterialMetadata 素材元数据
type MaterialMetadata struct {
	Duration    *int    `json:"duration,omitempty"`     // 时长（秒）
	Pages       *int    `json:"pages,omitempty"`        // 页数
	FileSize    *int64  `json:"file_size,omitempty"`    // 文件大小（字节）
	Format      string  `json:"format"`                 // 文件格式
	Resolution  *string `json:"resolution,omitempty"`   // 视频分辨率
	Language    string  `json:"language" gorm:"default:'zh-CN'"` // 语言
	Quality     *string `json:"quality,omitempty"`      // 质量等级
	Source      string  `json:"source"`                 // 来源
	Version     string  `json:"version" gorm:"default:'1.0'"` // 版本
}

// MaterialPermissions 素材权限控制
type MaterialPermissions struct {
	AllowedUsage   []UsageType `json:"allowed_usage" gorm:"type:text[]"`
	Licensing      LicenseInfo `json:"licensing" gorm:"embedded"`
	Restrictions   []string    `json:"restrictions" gorm:"type:text[]"`
	AccessLevel    string      `json:"access_level" gorm:"default:'public'"` // public/protected/private
	AllowedRoles   []string    `json:"allowed_roles" gorm:"type:text[]"`    // 允许的角色
	AllowedUsers   []uuid.UUID `json:"allowed_users" gorm:"type:uuid[]"`    // 允许的用户ID
}

// LicenseInfo 授权信息
type LicenseInfo struct {
	Type        string `json:"type"`         // creative-commons/copyright/etc
	Name        string `json:"name"`         // 许可证名称
	URL         string `json:"url"`          // 许可证URL
	Description string `json:"description"`  // 许可证描述
}

// MaterialStatistics 素材统计信息
type MaterialStatistics struct {
	ViewCount         int64 `json:"view_count" gorm:"default:0"`
	DownloadCount     int64 `json:"download_count" gorm:"default:0"`
	FavoriteCount     int64 `json:"favorite_count" gorm:"default:0"`
	ShareCount        int64 `json:"share_count" gorm:"default:0"`
	RatingCount       int64 `json:"rating_count" gorm:"default:0"`
	AverageRating     float64 `json:"average_rating" gorm:"type:decimal(3,2);default:0"`
	LastAccessedAt    *time.Time `json:"last_accessed_at,omitempty"`
	TotalAccessTime   int64 `json:"total_access_time" gorm:"default:0"` // 总访问时长（秒）

	// 按维度统计
	UsageByGrade     map[string]int64 `json:"usage_by_grade,omitempty" gorm:"type:jsonb"`
	UsageBySubject   map[string]int64 `json:"usage_by_subject,omitempty" gorm:"type:jsonb"`
	UsageByTimeRange map[string]int64 `json:"usage_by_time_range,omitempty" gorm:"type:jsonb"`
}

// SearchMaterialsRequest 搜索素材请求
type SearchMaterialsRequest struct {
	Query     string       `json:"query" form:"query"`
	Grade     []GradeLevel `json:"grade,omitempty" form:"grade"`
	Subject   Subject      `json:"subject,omitempty" form:"subject"`
	Type      MaterialType `json:"type,omitempty" form:"type"`
	Difficulty Difficulty   `json:"difficulty,omitempty" form:"difficulty"`
	Tags      []string     `json:"tags,omitempty" form:"tags"`

	Pagination PaginationRequest `json:"pagination"`
	Sort       SortRequest       `json:"sort"`

	// 高级搜索选项
	SemanticSearch bool     `json:"semantic_search,omitempty"` // 是否启用语义搜索
	MinScore       *float64 `json:"min_score,omitempty"`       // 最小匹配分数
	Filters        map[string]interface{} `json:"filters,omitempty"` // 自定义过滤器
}

// PaginationRequest 分页请求
type PaginationRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

// SortRequest 排序请求
type SortRequest struct {
	Field string `json:"field" form:"field"`
	Order string `json:"order" form:"order" binding:"oneof=asc desc"`
}

// SearchMaterialsResponse 搜索素材响应
type SearchMaterialsResponse struct {
	Materials   []TeachingMaterial `json:"materials"`
	Pagination  PaginationResponse `json:"pagination"`
	TotalCount  int64              `json:"total_count"`
	SearchTime  float64            `json:"search_time"`
	Query       string             `json:"query"`
	Suggestions []string           `json:"suggestions,omitempty"`
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	TotalPages int  `json:"total_pages"`
	TotalCount int64 `json:"total_count"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// MaterialDetailResponse 素材详情响应
type MaterialDetailResponse struct {
	*TeachingMaterial
	RelatedMaterials []RelatedMaterial `json:"related_materials,omitempty"`
	Recommendations  []TeachingMaterial `json:"recommendations,omitempty"`
	AccessURL        string            `json:"access_url,omitempty"`
	PreviewURL       string            `json:"preview_url,omitempty"`
}

// RelatedMaterial 相关素材
type RelatedMaterial struct {
	Material    TeachingMaterial `json:"material"`
	RelationType string          `json:"relation_type"` // prerequisite/followup/similar
	Strength     float64         `json:"strength"`      // 关联强度 (0-1)
}

// MaterialAnalysisRequest 素材分析请求
type MaterialAnalysisRequest struct {
	MaterialID uuid.UUID `json:"material_id" binding:"required"`
	AnalysisType string  `json:"analysis_type" binding:"required,oneof=difficulty alignment effectiveness"`
}

// MaterialAnalysisResponse 素材分析响应
type MaterialAnalysisResponse struct {
	MaterialID   uuid.UUID         `json:"material_id"`
	AnalysisType string           `json:"analysis_type"`
	Result       interface{}      `json:"result"`
	Confidence   float64          `json:"confidence"`
	AnalysisTime time.Time        `json:"analysis_time"`
}

// DifficultyAnalysis 难度分析结果
type DifficultyAnalysis struct {
	EstimatedDifficulty Difficulty `json:"estimated_difficulty"`
	Confidence          float64    `json:"confidence"`
	Reasons             []string   `json:"reasons"`
	ReadabilityScore    float64    `json:"readability_score"`
	ComplexityScore     float64    `json:"complexity_score"`
}

// AlignmentAnalysis 对齐分析结果
type AlignmentAnalysis struct {
	CurriculumStandard string  `json:"curriculum_standard"`
	AlignmentScore     float64 `json:"alignment_score"`
	Objectives         []ObjectiveAlignment `json:"objectives"`
	Gaps               []string `json:"gaps"`
	Recommendations    []string `json:"recommendations"`
}

// ObjectiveAlignment 目标对齐
type ObjectiveAlignment struct {
	Objective     string  `json:"objective"`
	AlignmentScore float64 `json:"alignment_score"`
	Evidence      string  `json:"evidence"`
}

// EffectivenessAnalysis 教学效果分析结果
type EffectivenessAnalysis struct {
	OverallScore     float64             `json:"overall_score"`
	EngagementScore  float64             `json:"engagement_score"`
	ComprehensionScore float64           `json:"comprehension_score"`
	RetentionScore   float64             `json:"retention_score"`
	Metrics          map[string]float64  `json:"metrics"`
	Insights         []string            `json:"insights"`
	Recommendations  []string            `json:"recommendations"`
}

// GenerateLessonPlanRequest 生成教案请求
type GenerateLessonPlanRequest struct {
	MaterialIDs []uuid.UUID `json:"material_ids" binding:"required,min=1"`
	Objectives  []string    `json:"objectives" binding:"required,min=1"`
	Grade       GradeLevel  `json:"grade" binding:"required"`
	Duration    int         `json:"duration" binding:"required,min=15,max=120"` // 分钟
	Style       string      `json:"style,omitempty" binding:"oneof=traditional activity inquiry project"`
}

// LessonPlanResponse 教案响应
type LessonPlanResponse struct {
	Title         string                `json:"title"`
	Grade         GradeLevel            `json:"grade"`
	Subject       Subject               `json:"subject"`
	Duration      int                   `json:"duration"`
	Objectives    []string              `json:"objectives"`
	Materials     []LessonPlanMaterial  `json:"materials"`
	Procedure     []LessonPlanStep      `json:"procedure"`
	Assessment    LessonPlanAssessment  `json:"assessment"`
	Differentiation LessonPlanDifferentiation `json:"differentiation"`
	Standards     []string              `json:"standards"`
	GeneratedAt   time.Time             `json:"generated_at"`
}

// LessonPlanMaterial 教案材料
type LessonPlanMaterial struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Purpose     string    `json:"purpose"`
	Timing      string    `json:"timing"`
}

// LessonPlanStep 教案步骤
type LessonPlanStep struct {
	Phase       string `json:"phase"`       // introduction/development/closure
	Title       string `json:"title"`
	Description string `json:"description"`
	Timing      int    `json:"timing"`      // 分钟
	Materials   []string `json:"materials"` // 材料ID列表
	Activities  []string `json:"activities"`
}

// LessonPlanAssessment 教案评估
type LessonPlanAssessment struct {
	Formative   []string `json:"formative"`   // 形成性评估
	Summative   []string `json:"summative"`   // 总结性评估
	Rubric      []RubricItem `json:"rubric"` // 评分标准
}

// RubricItem 评分标准项
type RubricItem struct {
	Criteria    string   `json:"criteria"`
	Levels      []string `json:"levels"`
	Points      int      `json:"points"`
}

// LessonPlanDifferentiation 教案差异化
type LessonPlanDifferentiation struct {
	Support     []string `json:"support"`     // 支持策略
	Extension   []string `json:"extension"`   // 延伸活动
	Grouping    string   `json:"grouping"`    // 分组策略
}
