package controller

import (
	"net/http"
	"strconv"

	"sitenavigation/config"
	"sitenavigation/model"
	"sitenavigation/service"

	"github.com/gin-gonic/gin"
)

var promptService *service.PromptService

func getPromptService() *service.PromptService {
	if promptService == nil {
		promptService = service.NewPromptService(config.DB)
	}
	return promptService
}

// GetPromptTree 获取模板树（对齐前端 PROMPT_DATA）
func GetPromptTree(c *gin.Context) {
	keyword := c.Query("keyword")
	tree, err := getPromptService().GetTree(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tree)
}

// GetPromptItem 获取单个模板
func GetPromptItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"})
		return
	}

	item, err := getPromptService().GetItemByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreatePromptCategory 创建分类
func CreatePromptCategory(c *gin.Context) {
	var req model.CreatePromptCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := getPromptService().CreateCategory(&req)
	if err != nil {
		if err.Error() == "分类名称已存在" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

// UpdatePromptCategory 更新分类
func UpdatePromptCategory(c *gin.Context) {
	var req model.UpdatePromptCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := getPromptService().UpdateCategory(&req)
	if err != nil {
		switch err.Error() {
		case "分类不存在":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "分类名称已存在":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, cat)
}

// DeletePromptCategory 删除分类（连同模板软删除）
func DeletePromptCategory(c *gin.Context) {
	var req model.DeleteByIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的分类ID"})
		return
	}

	if err := getPromptService().DeleteCategory(req.ID); err != nil {
		if err.Error() == "分类不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "分类删除成功"})
}

// CreatePromptItem 创建模板
func CreatePromptItem(c *gin.Context) {
	var req model.CreatePromptItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := getPromptService().CreateItem(&req)
	if err != nil {
		if err.Error() == "分类不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// UpdatePromptItem 更新模板（正文原样保存，含换行）
func UpdatePromptItem(c *gin.Context) {
	var req model.UpdatePromptItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := getPromptService().UpdateItem(&req)
	if err != nil {
		switch err.Error() {
		case "模板不存在", "分类不存在":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeletePromptItem 删除模板
func DeletePromptItem(c *gin.Context) {
	var req model.DeleteByIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"})
		return
	}

	if err := getPromptService().DeleteItem(req.ID); err != nil {
		if err.Error() == "模板不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "模板删除成功"})
}
