package controller

import (
	"net/http"
	"strconv"

	"sitenavigation/config"
	"sitenavigation/model"
	"sitenavigation/service"

	"github.com/gin-gonic/gin"
)

var userService *service.UserService

func getUserService() *service.UserService {
	if userService == nil {
		userService = service.NewUserService(config.DB)
	}
	return userService
}

func CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var createUserID *int
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(int); ok {
			createUserID = &uid
		}
	}

	user, err := getUserService().CreateUser(&req, createUserID)
	if err != nil {
		if err.Error() == "账号已存在" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	user, err := getUserService().GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetAllUsers(c *gin.Context) {
	users, err := getUserService().GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func UpdateUser(c *gin.Context) {
	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := getUserService().UpdateUser(req.ID, &req)
	if err != nil {
		if err.Error() == "用户不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := getUserService().DeleteUser(req.ID); err != nil {
		switch err.Error() {
		case "用户不存在":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "不允许删除 admin 用户":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}
