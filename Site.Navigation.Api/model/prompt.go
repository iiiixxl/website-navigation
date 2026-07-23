package model

import "time"

// PromptCategory AI 模板分类
type PromptCategory struct {
	ID         int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Title      string    `json:"title" gorm:"column:title;size:100;not null"`
	SortOrder  int       `json:"sort_order" gorm:"column:sort_order;default:0;not null"`
	IsDeleted  bool      `json:"is_deleted" gorm:"column:is_deleted;default:0;not null"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"`
}

func (PromptCategory) TableName() string {
	return "t_prompt_category"
}

// PromptItem AI 提问模板
type PromptItem struct {
	ID         int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	CategoryID int       `json:"category_id" gorm:"column:category_id;not null;index"`
	Name       string    `json:"name" gorm:"column:name;size:100;not null"`
	Content    string    `json:"content" gorm:"column:content;type:longtext;not null"`
	SortOrder  int       `json:"sort_order" gorm:"column:sort_order;default:0;not null"`
	IsDeleted  bool      `json:"is_deleted" gorm:"column:is_deleted;default:0;not null"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"`
}

func (PromptItem) TableName() string {
	return "t_prompt_item"
}

// ---------- 请求 / 响应 DTO ----------

type CreatePromptCategoryRequest struct {
	Title     string `json:"title" binding:"required,max=100"`
	SortOrder int    `json:"sort_order"`
}

type UpdatePromptCategoryRequest struct {
	ID        int    `json:"id" binding:"required"`
	Title     string `json:"title" binding:"required,max=100"`
	SortOrder int    `json:"sort_order"`
}

type CreatePromptItemRequest struct {
	CategoryID int    `json:"category_id" binding:"required"`
	Name       string `json:"name" binding:"required,max=100"`
	Content    string `json:"content" binding:"required"`
	SortOrder  int    `json:"sort_order"`
}

type UpdatePromptItemRequest struct {
	ID         int    `json:"id" binding:"required"`
	CategoryID int    `json:"category_id" binding:"required"`
	Name       string `json:"name" binding:"required,max=100"`
	Content    string `json:"content" binding:"required"`
	SortOrder  int    `json:"sort_order"`
}

type DeleteByIDRequest struct {
	ID int `json:"id" binding:"required"`
}

// PromptCategoryTree 分类树（对齐前端 PROMPT_DATA 结构）
type PromptCategoryTree struct {
	ID        int               `json:"id"`
	Title     string            `json:"title"`
	SortOrder int               `json:"sort_order"`
	Items     []PromptItemBrief `json:"items"`
}

type PromptItemBrief struct {
	ID         int    `json:"id"`
	CategoryID int    `json:"category_id"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	SortOrder  int    `json:"sort_order"`
}

type PromptTreeResponse struct {
	Categories []PromptCategoryTree `json:"categories"`
}
