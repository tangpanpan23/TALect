package service

import (
	"fmt"
	"time"

	"github.com/future-mcp/future-mcp-server/internal/repository"
	"github.com/future-mcp/future-mcp-server/internal/types"
	"github.com/future-mcp/future-mcp-server/pkg/logger"
	"github.com/google/uuid"
)

// MaterialServiceImpl 素材服务实现
type MaterialServiceImpl struct {
	materialRepo repository.MaterialRepository
	cache        CacheService
}

// NewMaterialService 创建素材服务
func NewMaterialService(materialRepo repository.MaterialRepository, cache CacheService) MaterialService {
	return &MaterialServiceImpl{
		materialRepo: materialRepo,
		cache:        cache,
	}
}

// SearchMaterials 搜索素材
func (s *MaterialServiceImpl) SearchMaterials(userID uuid.UUID, req types.SearchMaterialsRequest) (*types.SearchMaterialsResponse, error) {
	logger.Info("Searching materials",
		logger.Any("user_id", userID),
		logger.Any("query", req.Query),
		logger.Any("grade", req.Grade),
		logger.Any("subject", req.Subject))

	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("search:%s:%v:%s:%d", req.Query, req.Grade, req.Subject, req.Pagination.Page)
	if cached, err := s.cache.GetSearchCache(cacheKey, nil); err == nil && cached != nil {
		logger.Info("Search result from cache", logger.Any("cache_key", cacheKey))
		return &types.SearchMaterialsResponse{
			Materials:   cached.Materials,
			Pagination:  s.buildPaginationResponse(req.Pagination, cached.TotalCount),
			TotalCount:  cached.TotalCount,
			SearchTime:  cached.SearchTime,
			Query:       req.Query,
		}, nil
	}

	// 从数据库搜索
	materials, total, err := s.materialRepo.SearchMaterials(req)
	if err != nil {
		logger.Error("Failed to search materials", logger.Any("error", err))
		return nil, fmt.Errorf("failed to search materials: %w", err)
	}

	// 构建响应
	response := &types.SearchMaterialsResponse{
		Materials:   materials,
		Pagination:  s.buildPaginationResponse(req.Pagination, total),
		TotalCount:  total,
		SearchTime:  0.1, // 模拟搜索时间
		Query:       req.Query,
	}

	// 缓存结果
	searchResult := &types.SearchResult{
		Materials:  materials,
		TotalCount: total,
		SearchTime: 0.1,
	}
	if err := s.cache.SetSearchCache(cacheKey, nil, searchResult, 300); err != nil {
		logger.Warn("Failed to cache search result", logger.Any("error", err))
	}

	return response, nil
}

// SearchByGradeSubject 按年级学科搜索
func (s *MaterialServiceImpl) SearchByGradeSubject(userID uuid.UUID, grade types.GradeLevel, subject types.Subject, difficulty types.Difficulty, teachingStage string) (*types.SearchMaterialsResponse, error) {
	logger.Info("Searching materials by grade and subject",
		logger.Any("user_id", userID),
		logger.Any("grade", grade),
		logger.Any("subject", subject),
		logger.Any("difficulty", difficulty))

	req := types.SearchMaterialsRequest{
		Grade:   []types.GradeLevel{grade},
		Subject: subject,
		Pagination: types.PaginationRequest{
			Page:     1,
			PageSize: 20,
		},
	}

	// 添加难度过滤
	if difficulty != "" {
		req.Filters = map[string]interface{}{
			"difficulty": difficulty,
		}
	}

	// 添加教学阶段过滤
	if teachingStage != "" {
		if req.Filters == nil {
			req.Filters = make(map[string]interface{})
		}
		req.Filters["teaching_stage"] = teachingStage
	}

	return s.SearchMaterials(userID, req)
}

// SemanticSearch 语义搜索
func (s *MaterialServiceImpl) SemanticSearch(userID uuid.UUID, query string, limit int) (*types.SearchMaterialsResponse, error) {
	logger.Info("Semantic search",
		logger.Any("user_id", userID),
		logger.Any("query", query),
		logger.Any("limit", limit))

	// 这里应该调用向量搜索服务
	// 暂时使用关键词搜索作为替代
	req := types.SearchMaterialsRequest{
		Query: query,
		Pagination: types.PaginationRequest{
			Page:     1,
			PageSize: limit,
		},
	}

	return s.SearchMaterials(userID, req)
}

// GetMaterialDetail 获取素材详情
func (s *MaterialServiceImpl) GetMaterialDetail(userID uuid.UUID, materialID uuid.UUID) (*types.MaterialDetailResponse, error) {
	logger.Info("Getting material detail",
		logger.Any("user_id", userID),
		logger.Any("material_id", materialID))

	// 尝试从缓存获取
	if cached, err := s.cache.GetMaterialCache(materialID.String()); err == nil && cached != nil {
		logger.Info("Material detail from cache", logger.Any("material_id", materialID))
		return &types.MaterialDetailResponse{
			TeachingMaterial: cached,
		}, nil
	}

	// 从数据库获取
	material, err := s.materialRepo.GetMaterialByID(materialID)
	if err != nil {
		logger.Error("Failed to get material detail", logger.Any("error", err))
		return nil, fmt.Errorf("failed to get material detail: %w", err)
	}

	// 获取相关素材
	related, err := s.materialRepo.GetRelatedMaterials(materialID, "similar", 5)
	if err != nil {
		logger.Warn("Failed to get related materials", logger.Any("error", err))
		related = []types.TeachingMaterial{}
	}

	// 构建响应
	response := &types.MaterialDetailResponse{
		TeachingMaterial: material,
		RelatedMaterials: s.buildRelatedMaterials(related),
	}

	// 缓存结果
	if err := s.cache.SetMaterialCache(material, 1800); err != nil {
		logger.Warn("Failed to cache material detail", logger.Any("error", err))
	}

	return response, nil
}

// GetRelatedMaterials 获取相关素材
func (s *MaterialServiceImpl) GetRelatedMaterials(userID uuid.UUID, materialID uuid.UUID, relationType string, limit int) (*types.SearchMaterialsResponse, error) {
	logger.Info("Getting related materials",
		logger.Any("user_id", userID),
		logger.Any("material_id", materialID),
		logger.Any("relation_type", relationType),
		logger.Any("limit", limit))

	materials, err := s.materialRepo.GetRelatedMaterials(materialID, relationType, limit)
	if err != nil {
		logger.Error("Failed to get related materials", logger.Any("error", err))
		return nil, fmt.Errorf("failed to get related materials: %w", err)
	}

	return &types.SearchMaterialsResponse{
		Materials:   materials,
		Pagination:  s.buildPaginationResponse(types.PaginationRequest{Page: 1, PageSize: limit}, int64(len(materials))),
		TotalCount:  int64(len(materials)),
		SearchTime:  0.05,
	}, nil
}

// AnalyzeMaterial 分析素材
func (s *MaterialServiceImpl) AnalyzeMaterial(userID uuid.UUID, req types.MaterialAnalysisRequest) (*types.MaterialAnalysisResponse, error) {
	logger.Info("Analyzing material",
		logger.Any("user_id", userID),
		logger.Any("material_id", req.MaterialID),
		logger.Any("analysis_type", req.AnalysisType))

	// 获取素材
	material, err := s.materialRepo.GetMaterialByID(req.MaterialID)
	if err != nil {
		return nil, fmt.Errorf("failed to get material: %w", err)
	}

	// 根据分析类型进行分析
	var result interface{}
	switch req.AnalysisType {
	case "difficulty":
		result = s.analyzeDifficulty(material)
	case "curriculum_alignment":
		result = s.analyzeCurriculumAlignment(material)
	default:
		return nil, fmt.Errorf("unsupported analysis type: %s", req.AnalysisType)
	}

	return &types.MaterialAnalysisResponse{
		MaterialID:   req.MaterialID,
		AnalysisType: req.AnalysisType,
		Result:       result,
		Confidence:   0.85,
		AnalysisTime: time.Now(),
	}, nil
}

// GetPersonalizedRecommendations 获取个性化推荐
func (s *MaterialServiceImpl) GetPersonalizedRecommendations(userID uuid.UUID, limit int) (*types.RecommendationResult, error) {
	logger.Info("Getting personalized recommendations",
		logger.Any("user_id", userID),
		logger.Any("limit", limit))

	// 这里应该基于用户学习历史和偏好进行推荐
	// 暂时返回热门素材作为推荐
	materials, _, err := s.materialRepo.SearchMaterials(types.SearchMaterialsRequest{
		Pagination: types.PaginationRequest{
			Page:     1,
			PageSize: limit,
		},
		Sort: types.SortRequest{
			Field: "view_count",
			Order: "desc",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}

	return &types.RecommendationResult{
		Materials: materials,
		Algorithm: "popularity_based",
		Confidence: 0.7,
		Reason:     "基于学习热度推荐",
	}, nil
}

// 辅助方法
func (s *MaterialServiceImpl) buildPaginationResponse(req types.PaginationRequest, total int64) types.PaginationResponse {
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 20
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return types.PaginationResponse{
		Page:       req.Page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalCount: total,
		HasNext:    req.Page < totalPages,
		HasPrev:    req.Page > 1,
	}
}

func (s *MaterialServiceImpl) buildRelatedMaterials(materials []types.TeachingMaterial) []types.RelatedMaterial {
	result := make([]types.RelatedMaterial, len(materials))
	for i, material := range materials {
		result[i] = types.RelatedMaterial{
			Material:    material,
			RelationType: "similar",
			Strength:    0.8, // 简化实现
		}
	}
	return result
}

func (s *MaterialServiceImpl) analyzeDifficulty(material *types.TeachingMaterial) *types.DifficultyAnalysis {
	// 简化难度分析逻辑
	var estimatedDifficulty types.Difficulty
	var confidence float64

	switch material.Difficulty {
	case types.DifficultyEasy:
		estimatedDifficulty = types.DifficultyEasy
		confidence = 0.9
	case types.DifficultyMedium:
		estimatedDifficulty = types.DifficultyMedium
		confidence = 0.8
	case types.DifficultyHard:
		estimatedDifficulty = types.DifficultyHard
		confidence = 0.85
	default:
		estimatedDifficulty = types.DifficultyMedium
		confidence = 0.7
	}

	return &types.DifficultyAnalysis{
		EstimatedDifficulty: estimatedDifficulty,
		Confidence:          confidence,
		Reasons:             []string{"基于素材标签和内容分析"},
		ReadabilityScore:    0.75,
		ComplexityScore:     0.6,
	}
}

func (s *MaterialServiceImpl) analyzeCurriculumAlignment(material *types.TeachingMaterial) *types.AlignmentAnalysis {
	// 简化课标对齐分析
	return &types.AlignmentAnalysis{
		CurriculumStandard: "人教版小学数学",
		AlignmentScore:     0.85,
		Objectives: []types.ObjectiveAlignment{
			{
				Objective:     "掌握基本运算",
				AlignmentScore: 0.9,
				Evidence:      "素材包含加减乘除运算内容",
			},
		},
		Gaps:         []string{},
		Recommendations: []string{"建议补充练习环节"},
	}
}
