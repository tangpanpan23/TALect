package repository

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/future-mcp/future-mcp-server/internal/types"
	"github.com/google/uuid"
)

// MemoryMaterialRepository 内存素材仓库实现（用于测试）
type MemoryMaterialRepository struct {
	materials map[uuid.UUID]*types.TeachingMaterial
	mu        sync.RWMutex
}

// NewMemoryMaterialRepository 创建内存素材仓库
func NewMemoryMaterialRepository() MaterialRepository {
	repo := &MemoryMaterialRepository{
		materials: make(map[uuid.UUID]*types.TeachingMaterial),
	}

	// 初始化一些测试数据
	repo.initSampleData()

	return repo
}

// initSampleData 初始化示例数据
func (r *MemoryMaterialRepository) initSampleData() {
	materials := []*types.TeachingMaterial{
		// 数学类素材
		{
			ID:          uuid.New(),
			Title:       "一元二次方程解法",
			Description: "详细讲解一元二次方程的解题方法和技巧，包含十字相乘法、公式法等多种解题思路",
			Type:        types.MaterialTypeVideo,
			GradeLevels: []types.GradeLevel{types.GradeLevel2},
			Subject:     types.SubjectMath,
			Tags:        []string{"代数", "方程", "二次方程", "解题技巧"},
			Difficulty:  types.DifficultyMedium,
			CurriculumAlignment: types.CurriculumAlignment{
				Standard:        "人教版初中数学",
				Objectives:      []string{"掌握一元二次方程解法", "能够灵活运用多种解题方法"},
				CompetencyLevel: 85,
			},
			Metadata: types.MaterialMetadata{
				Duration:    &[]int{1800}[0], // 30分钟
				Format:      "mp4",
				Language:    "zh-CN",
			},
			Permissions: types.MaterialPermissions{
				AllowedUsage: []types.UsageType{types.UsageTypeView, types.UsageTypeDownload},
				Licensing: types.LicenseInfo{
					Type: "creative-commons",
					Name: "知识共享",
				},
			},
			Statistics: types.MaterialStatistics{
				ViewCount:    1250,
				DownloadCount: 320,
				RatingCount:  45,
				AverageRating: 4.6,
			},
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Title:       "分数四则运算",
			Description: "小学分数基础运算，包括分数加减乘除的计算方法和通分技巧",
			Type:        types.MaterialTypeInteractive,
			GradeLevels: []types.GradeLevel{types.GradeLevel1},
			Subject:     types.SubjectMath,
			Tags:        []string{"分数", "四则运算", "通分", "小学数学"},
			Difficulty:  types.DifficultyEasy,
			CurriculumAlignment: types.CurriculumAlignment{
				Standard:        "人教版小学数学",
				Objectives:      []string{"掌握分数四则运算", "理解分数的基本概念"},
				CompetencyLevel: 90,
			},
			Metadata: types.MaterialMetadata{
				Duration:    &[]int{1200}[0], // 20分钟
				Format:      "html",
				Language:    "zh-CN",
			},
			Permissions: types.MaterialPermissions{
				AllowedUsage: []types.UsageType{types.UsageTypeView, types.UsageTypeInteractive},
				Licensing: types.LicenseInfo{
					Type: "creative-commons",
					Name: "知识共享",
				},
			},
			Statistics: types.MaterialStatistics{
				ViewCount:    2100,
				DownloadCount: 150,
				RatingCount:  67,
				AverageRating: 4.8,
			},
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Title:       "几何证明题解题思路",
			Description: "高中几何证明题的解题方法和思维训练，包含多种证明技巧",
			Type:        types.MaterialTypeExercise,
			GradeLevels: []types.GradeLevel{types.GradeLevel3},
			Subject:     types.SubjectMath,
			Tags:        []string{"几何", "证明题", "思维训练", "高中数学"},
			Difficulty:  types.DifficultyHard,
			CurriculumAlignment: types.CurriculumAlignment{
				Standard:        "人教版高中数学",
				Objectives:      []string{"掌握几何证明方法", "培养逻辑思维能力"},
				CompetencyLevel: 70,
			},
			Metadata: types.MaterialMetadata{
				Pages:       &[]int{30}[0],
				Format:      "pdf",
				Language:    "zh-CN",
			},
			Permissions: types.MaterialPermissions{
				AllowedUsage: []types.UsageType{types.UsageTypeView, types.UsageTypeDownload},
				Licensing: types.LicenseInfo{
					Type: "creative-commons",
					Name: "知识共享",
				},
			},
			Statistics: types.MaterialStatistics{
				ViewCount:    750,
				DownloadCount: 280,
				RatingCount:  28,
				AverageRating: 4.3,
			},
			CreatedAt: time.Now().Add(-18 * time.Hour),
			UpdatedAt: time.Now(),
		},
		// 英语类素材
		{
			ID:          uuid.New(),
			Title:       "英语语法基础：时态",
			Description: "系统学习英语16种时态的用法和区别，包含练习和常见错误分析",
			Type:        types.MaterialTypePPT,
			GradeLevels: []types.GradeLevel{types.GradeLevel1, types.GradeLevel2},
			Subject:     types.SubjectEnglish,
			Tags:        []string{"语法", "时态", "英语学习", "语言学习"},
			Difficulty:  types.DifficultyEasy,
			CurriculumAlignment: types.CurriculumAlignment{
				Standard:        "人教版初中英语",
				Objectives:      []string{"掌握英语时态用法", "能够正确使用各种时态"},
				CompetencyLevel: 90,
			},
			Metadata: types.MaterialMetadata{
				Pages:       &[]int{45}[0],
				Format:      "pptx",
				Language:    "zh-CN",
			},
			Permissions: types.MaterialPermissions{
				AllowedUsage: []types.UsageType{types.UsageTypeView},
				Licensing: types.LicenseInfo{
					Type: "copyright",
					Name: "版权保护",
				},
			},
			Statistics: types.MaterialStatistics{
				ViewCount:    2100,
				DownloadCount: 180,
				RatingCount:  78,
				AverageRating: 4.8,
			},
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Title:       "英语单词记忆方法",
			Description: "高效的英语单词记忆技巧和方法，适合不同年级的学生",
			Type:        types.MaterialTypeVideo,
			GradeLevels: []types.GradeLevel{types.GradeLevel1, types.GradeLevel2, types.GradeLevel3},
			Subject:     types.SubjectEnglish,
			Tags:        []string{"单词记忆", "学习方法", "英语学习技巧"},
			Difficulty:  types.DifficultyEasy,
			CurriculumAlignment: types.CurriculumAlignment{
				Standard:        "通用英语教学",
				Objectives:      []string{"掌握单词记忆方法", "提高词汇量"},
				CompetencyLevel: 95,
			},
			Metadata: types.MaterialMetadata{
				Duration:    &[]int{1500}[0], // 25分钟
				Format:      "mp4",
				Language:    "zh-CN",
			},
			Permissions: types.MaterialPermissions{
				AllowedUsage: []types.UsageType{types.UsageTypeView, types.UsageTypeDownload},
				Licensing: types.LicenseInfo{
					Type: "creative-commons",
					Name: "知识共享",
				},
			},
			Statistics: types.MaterialStatistics{
				ViewCount:    3200,
				DownloadCount: 450,
				RatingCount:  120,
				AverageRating: 4.9,
			},
			CreatedAt: time.Now().Add(-72 * time.Hour),
			UpdatedAt: time.Now(),
		},
		// 物理类素材
		{
			ID:          uuid.New(),
			Title:       "物理力学基础",
			Description: "牛顿力学核心概念和定律详解，包含实验演示和习题解析",
			Type:        types.MaterialTypeExercise,
			GradeLevels: []types.GradeLevel{types.GradeLevel2, types.GradeLevel3},
			Subject:     types.SubjectPhysics,
			Tags:        []string{"力学", "牛顿定律", "物理实验", "高中物理"},
			Difficulty:  types.DifficultyHard,
			CurriculumAlignment: types.CurriculumAlignment{
				Standard:        "人教版高中物理",
				Objectives:      []string{"理解牛顿运动定律", "掌握力学基本概念"},
				CompetencyLevel: 75,
			},
			Metadata: types.MaterialMetadata{
				Format:      "pdf",
				Language:    "zh-CN",
			},
			Permissions: types.MaterialPermissions{
				AllowedUsage: []types.UsageType{types.UsageTypeView, types.UsageTypeDownload},
				Licensing: types.LicenseInfo{
					Type: "creative-commons",
					Name: "知识共享",
				},
			},
			Statistics: types.MaterialStatistics{
				ViewCount:    890,
				DownloadCount: 450,
				RatingCount:  32,
				AverageRating: 4.4,
			},
			CreatedAt: time.Now().Add(-12 * time.Hour),
			UpdatedAt: time.Now(),
		},
		// 语文类素材
		{
			ID:          uuid.New(),
			Title:       "古诗文赏析方法",
			Description: "初中古诗文阅读和赏析技巧，包含名家赏析和学生范文",
			Type:        types.MaterialTypeDocument,
			GradeLevels: []types.GradeLevel{types.GradeLevel1, types.GradeLevel2},
			Subject:     types.SubjectChinese,
			Tags:        []string{"古诗文", "赏析", "阅读技巧", "语文学习"},
			Difficulty:  types.DifficultyMedium,
			CurriculumAlignment: types.CurriculumAlignment{
				Standard:        "人教版初中语文",
				Objectives:      []string{"掌握古诗文赏析方法", "提高文学素养"},
				CompetencyLevel: 80,
			},
			Metadata: types.MaterialMetadata{
				Pages:       &[]int{25}[0],
				Format:      "docx",
				Language:    "zh-CN",
			},
			Permissions: types.MaterialPermissions{
				AllowedUsage: []types.UsageType{types.UsageTypeView, types.UsageTypeDownload},
				Licensing: types.LicenseInfo{
					Type: "creative-commons",
					Name: "知识共享",
				},
			},
			Statistics: types.MaterialStatistics{
				ViewCount:    1650,
				DownloadCount: 380,
				RatingCount:  55,
				AverageRating: 4.7,
			},
			CreatedAt: time.Now().Add(-36 * time.Hour),
			UpdatedAt: time.Now(),
		},
		// 化学类素材
		{
			ID:          uuid.New(),
			Title:       "化学反应基本类型",
			Description: "高中化学反应分类和基本规律，包含实验演示视频",
			Type:        types.MaterialTypeVideo,
			GradeLevels: []types.GradeLevel{types.GradeLevel2, types.GradeLevel3},
			Subject:     types.SubjectChemistry,
			Tags:        []string{"化学反应", "反应类型", "实验", "高中化学"},
			Difficulty:  types.DifficultyMedium,
			CurriculumAlignment: types.CurriculumAlignment{
				Standard:        "人教版高中化学",
				Objectives:      []string{"掌握化学反应基本类型", "理解反应规律"},
				CompetencyLevel: 85,
			},
			Metadata: types.MaterialMetadata{
				Duration:    &[]int{2100}[0], // 35分钟
				Format:      "mp4",
				Language:    "zh-CN",
			},
			Permissions: types.MaterialPermissions{
				AllowedUsage: []types.UsageType{types.UsageTypeView},
				Licensing: types.LicenseInfo{
					Type: "copyright",
					Name: "版权保护",
				},
			},
			Statistics: types.MaterialStatistics{
				ViewCount:    980,
				DownloadCount: 120,
				RatingCount:  38,
				AverageRating: 4.5,
			},
			CreatedAt: time.Now().Add(-30 * time.Hour),
			UpdatedAt: time.Now(),
		},
		// 生物类素材
		{
			ID:          uuid.New(),
			Title:       "细胞结构与功能",
			Description: "初中生物细胞基础知识，包含微观结构图解和互动练习",
			Type:        types.MaterialTypeInteractive,
			GradeLevels: []types.GradeLevel{types.GradeLevel1, types.GradeLevel2},
			Subject:     types.SubjectBiology,
			Tags:        []string{"细胞", "生物学", "微观结构", "初中生物"},
			Difficulty:  types.DifficultyEasy,
			CurriculumAlignment: types.CurriculumAlignment{
				Standard:        "人教版初中生物",
				Objectives:      []string{"了解细胞基本结构", "理解细胞功能"},
				CompetencyLevel: 88,
			},
			Metadata: types.MaterialMetadata{
				Duration:    &[]int{900}[0], // 15分钟
				Format:      "html",
				Language:    "zh-CN",
			},
			Permissions: types.MaterialPermissions{
				AllowedUsage: []types.UsageType{types.UsageTypeView, types.UsageTypeInteractive},
				Licensing: types.LicenseInfo{
					Type: "creative-commons",
					Name: "知识共享",
				},
			},
			Statistics: types.MaterialStatistics{
				ViewCount:    1400,
				DownloadCount: 200,
				RatingCount:  42,
				AverageRating: 4.6,
			},
			CreatedAt: time.Now().Add(-42 * time.Hour),
			UpdatedAt: time.Now(),
		},
	}

	for _, material := range materials {
		r.materials[material.ID] = material
	}
}

// CreateMaterial 创建素材
func (r *MemoryMaterialRepository) CreateMaterial(material *types.TeachingMaterial) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if material.ID == uuid.Nil {
		material.ID = uuid.New()
	}

	material.CreatedAt = time.Now()
	material.UpdatedAt = time.Now()

	r.materials[material.ID] = material
	return nil
}

// GetMaterialByID 根据ID获取素材
func (r *MemoryMaterialRepository) GetMaterialByID(id uuid.UUID) (*types.TeachingMaterial, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	material, exists := r.materials[id]
	if !exists {
		return nil, fmt.Errorf("material not found: %s", id)
	}

	return material, nil
}

// UpdateMaterial 更新素材
func (r *MemoryMaterialRepository) UpdateMaterial(material *types.TeachingMaterial) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.materials[material.ID]; !exists {
		return fmt.Errorf("material not found: %s", material.ID)
	}

	material.UpdatedAt = time.Now()
	r.materials[material.ID] = material
	return nil
}

// DeleteMaterial 删除素材
func (r *MemoryMaterialRepository) DeleteMaterial(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.materials[id]; !exists {
		return fmt.Errorf("material not found: %s", id)
	}

	delete(r.materials, id)
	return nil
}

// SearchMaterials 搜索素材
func (r *MemoryMaterialRepository) SearchMaterials(req types.SearchMaterialsRequest) ([]types.TeachingMaterial, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []types.TeachingMaterial

	// 简单筛选逻辑
	for _, material := range r.materials {
		// 关键词匹配
		if req.Query != "" {
			if !containsIgnoreCase(material.Title, req.Query) &&
			   !containsIgnoreCase(material.Description, req.Query) {
				continue
			}
		}

		// 年级筛选
		if len(req.Grade) > 0 {
			found := false
			for _, grade := range req.Grade {
				for _, materialGrade := range material.GradeLevels {
					if materialGrade == grade {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				continue
			}
		}

		// 学科筛选
		if req.Subject != "" && material.Subject != req.Subject {
			continue
		}

		results = append(results, *material)
	}

	// 分页
	total := int64(len(results))
	start := (req.Pagination.Page - 1) * req.Pagination.PageSize
	end := start + req.Pagination.PageSize

	if start >= len(results) {
		return []types.TeachingMaterial{}, total, nil
	}

	if end > len(results) {
		end = len(results)
	}

	return results[start:end], total, nil
}

// GetRelatedMaterials 获取相关素材
func (r *MemoryMaterialRepository) GetRelatedMaterials(materialID uuid.UUID, relationType string, limit int) ([]types.TeachingMaterial, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	material, exists := r.materials[materialID]
	if !exists {
		return nil, fmt.Errorf("material not found: %s", materialID)
	}

	var related []types.TeachingMaterial

	// 简单相关性逻辑：相同学科和年级的素材
	for _, m := range r.materials {
		if m.ID == materialID {
			continue
		}

		// 检查是否相同学科
		if m.Subject != material.Subject {
			continue
		}

		// 检查是否有相同年级
		hasCommonGrade := false
		for _, mg := range material.GradeLevels {
			for _, og := range m.GradeLevels {
				if mg == og {
					hasCommonGrade = true
					break
				}
			}
			if hasCommonGrade {
				break
			}
		}

		if hasCommonGrade {
			related = append(related, *m)
			if len(related) >= limit {
				break
			}
		}
	}

	return related, nil
}

// GetPopularMaterials 获取热门素材
func (r *MemoryMaterialRepository) GetPopularMaterials(limit int) ([]types.TeachingMaterial, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var materials []types.TeachingMaterial
	for _, m := range r.materials {
		materials = append(materials, *m)
	}

	// 按查看次数排序
	for i := 0; i < len(materials)-1; i++ {
		for j := i + 1; j < len(materials); j++ {
			if materials[j].Statistics.ViewCount > materials[i].Statistics.ViewCount {
				materials[i], materials[j] = materials[j], materials[i]
			}
		}
	}

	if len(materials) > limit {
		materials = materials[:limit]
	}

	return materials, nil
}

// GetMaterialsByGrade 按年级获取素材
func (r *MemoryMaterialRepository) GetMaterialsByGrade(grade types.GradeLevel, limit int) ([]types.TeachingMaterial, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var materials []types.TeachingMaterial
	for _, m := range r.materials {
		for _, g := range m.GradeLevels {
			if g == grade {
				materials = append(materials, *m)
				break
			}
		}
		if len(materials) >= limit {
			break
		}
	}

	return materials, nil
}

// GetMaterialsBySubject 按学科获取素材
func (r *MemoryMaterialRepository) GetMaterialsBySubject(subject types.Subject, limit int) ([]types.TeachingMaterial, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var materials []types.TeachingMaterial
	for _, m := range r.materials {
		if m.Subject == subject {
			materials = append(materials, *m)
			if len(materials) >= limit {
				break
			}
		}
	}

	return materials, nil
}

// BatchCreateMaterials 批量创建素材
func (r *MemoryMaterialRepository) BatchCreateMaterials(materials []*types.TeachingMaterial) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, material := range materials {
		if material.ID == uuid.Nil {
			material.ID = uuid.New()
		}
		material.CreatedAt = time.Now()
		material.UpdatedAt = time.Now()
		r.materials[material.ID] = material
	}

	return nil
}

// BatchUpdateMaterials 批量更新素材
func (r *MemoryMaterialRepository) BatchUpdateMaterials(materials []*types.TeachingMaterial) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, material := range materials {
		if _, exists := r.materials[material.ID]; !exists {
			return fmt.Errorf("material not found: %s", material.ID)
		}
		material.UpdatedAt = time.Now()
		r.materials[material.ID] = material
	}

	return nil
}

// 辅助函数
func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}
