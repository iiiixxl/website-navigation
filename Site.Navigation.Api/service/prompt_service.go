package service

import (
	"errors"
	"time"

	"sitenavigation/model"

	"gorm.io/gorm"
)

type PromptService struct {
	db *gorm.DB
}

func NewPromptService(db *gorm.DB) *PromptService {
	return &PromptService{db: db}
}

// GetTree 获取分类 + 模板树（供前端 AI 页使用）
func (s *PromptService) GetTree(keyword string) (*model.PromptTreeResponse, error) {
	var categories []model.PromptCategory
	if err := s.db.Where("is_deleted = ?", false).
		Order("sort_order ASC, id ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	resp := &model.PromptTreeResponse{Categories: make([]model.PromptCategoryTree, 0, len(categories))}
	for _, cat := range categories {
		query := s.db.Model(&model.PromptItem{}).
			Where("category_id = ? AND is_deleted = ?", cat.ID, false)
		if keyword != "" {
			like := "%" + keyword + "%"
			query = query.Where("name LIKE ? OR content LIKE ?", like, like)
		}

		var items []model.PromptItem
		if err := query.Order("sort_order ASC, id ASC").Find(&items).Error; err != nil {
			return nil, err
		}
		if keyword != "" && len(items) == 0 {
			continue
		}

		briefs := make([]model.PromptItemBrief, 0, len(items))
		for _, item := range items {
			briefs = append(briefs, model.PromptItemBrief{
				ID:         item.ID,
				CategoryID: item.CategoryID,
				Name:       item.Name,
				Content:    item.Content,
				SortOrder:  item.SortOrder,
			})
		}

		resp.Categories = append(resp.Categories, model.PromptCategoryTree{
			ID:        cat.ID,
			Title:     cat.Title,
			SortOrder: cat.SortOrder,
			Items:     briefs,
		})
	}
	return resp, nil
}

func (s *PromptService) GetItemByID(id int) (*model.PromptItem, error) {
	var item model.PromptItem
	if err := s.db.Where("id = ? AND is_deleted = ?", id, false).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("模板不存在")
		}
		return nil, err
	}
	return &item, nil
}

func (s *PromptService) CreateCategory(req *model.CreatePromptCategoryRequest) (*model.PromptCategory, error) {
	var count int64
	if err := s.db.Model(&model.PromptCategory{}).
		Where("title = ? AND is_deleted = ?", req.Title, false).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("分类名称已存在")
	}

	now := time.Now()
	cat := &model.PromptCategory{
		Title:      req.Title,
		SortOrder:  req.SortOrder,
		IsDeleted:  false,
		CreateTime: now,
		UpdateTime: now,
	}
	if err := s.db.Create(cat).Error; err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *PromptService) UpdateCategory(req *model.UpdatePromptCategoryRequest) (*model.PromptCategory, error) {
	var cat model.PromptCategory
	if err := s.db.Where("id = ? AND is_deleted = ?", req.ID, false).First(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, err
	}

	var count int64
	if err := s.db.Model(&model.PromptCategory{}).
		Where("title = ? AND is_deleted = ? AND id <> ?", req.Title, false, req.ID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("分类名称已存在")
	}

	if err := s.db.Model(&cat).Updates(map[string]interface{}{
		"title":       req.Title,
		"sort_order":  req.SortOrder,
		"update_time": time.Now(),
	}).Error; err != nil {
		return nil, err
	}
	return s.getCategoryByID(req.ID)
}

func (s *PromptService) DeleteCategory(id int) error {
	cat, err := s.getCategoryByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.PromptItem{}).
			Where("category_id = ? AND is_deleted = ?", id, false).
			Updates(map[string]interface{}{
				"is_deleted":  true,
				"update_time": now,
			}).Error; err != nil {
			return err
		}
		return tx.Model(cat).Updates(map[string]interface{}{
			"is_deleted":  true,
			"update_time": now,
		}).Error
	})
}

func (s *PromptService) CreateItem(req *model.CreatePromptItemRequest) (*model.PromptItem, error) {
	if _, err := s.getCategoryByID(req.CategoryID); err != nil {
		return nil, err
	}

	now := time.Now()
	item := &model.PromptItem{
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Content:    req.Content,
		SortOrder:  req.SortOrder,
		IsDeleted:  false,
		CreateTime: now,
		UpdateTime: now,
	}
	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (s *PromptService) UpdateItem(req *model.UpdatePromptItemRequest) (*model.PromptItem, error) {
	item, err := s.GetItemByID(req.ID)
	if err != nil {
		return nil, err
	}
	if _, err := s.getCategoryByID(req.CategoryID); err != nil {
		return nil, err
	}

	if err := s.db.Model(item).Updates(map[string]interface{}{
		"category_id": req.CategoryID,
		"name":        req.Name,
		"content":     req.Content,
		"sort_order":  req.SortOrder,
		"update_time": time.Now(),
	}).Error; err != nil {
		return nil, err
	}
	return s.GetItemByID(req.ID)
}

func (s *PromptService) DeleteItem(id int) error {
	item, err := s.GetItemByID(id)
	if err != nil {
		return err
	}
	return s.db.Model(item).Updates(map[string]interface{}{
		"is_deleted":  true,
		"update_time": time.Now(),
	}).Error
}

func (s *PromptService) getCategoryByID(id int) (*model.PromptCategory, error) {
	var cat model.PromptCategory
	if err := s.db.Where("id = ? AND is_deleted = ?", id, false).First(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, err
	}
	return &cat, nil
}
