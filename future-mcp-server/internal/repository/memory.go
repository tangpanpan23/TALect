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
		{
			ID:          uuid.New(),
			Title:       "一元二次方程解法",
			Description: "详细讲解一元二次方程的解题方法和技巧",
			Type:        types.MaterialTypeVideo,
			GradeLevels: []types.GradeLevel{types.GradeLevel2},
			Subject:     types.SubjectMath,
			Tags:        []string{"代数", "方程", "二次方程"},
			Difficulty:  types.DifficultyMedium,
			CurriculumAlignment: types.CurriculumAlignment{
				Standard:        "人教版初中数学",
				Objectives:      []string{"掌握一元二次方程解法"},
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
			Title:       "英语语法基础：时态",
			Description: "系统学习英语16种时态的用法和区别",
			Type:        types.MaterialTypePPT,
			GradeLevels: []types.GradeLevel{types.GradeLevel1, types.GradeLevel2},
			Subject:     types.SubjectEnglish,
			Tags:        []string{"语法", "时态", "英语学习"},
			Difficulty:  types.DifficultyEasy,
			CurriculumAlignment: types.CurriculumAlignment{
				Standard:        "人教版初中英语",
				Objectives:      []string{"掌握英语时态用法"},
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
			Title:       "物理力学基础",
			Description: "牛顿力学核心概念和定律详解",
			Type:        types.MaterialTypeExercise,
			GradeLevels: []types.GradeLevel{types.GradeLevel2, types.GradeLevel3},
			Subject:     types.SubjectPhysics,
			Tags:        []string{"力学", "牛顿定律", "物理"},
			Difficulty:  types.DifficultyHard,
			CurriculumAlignment: types.CurriculumAlignment{
				Standard:        "人教版高中物理",
				Objectives:      []string{"理解牛顿运动定律"},
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
