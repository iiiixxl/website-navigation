package controller

import (
	"net/http"

	"sitenavigation/config"
	"sitenavigation/model"
	"sitenavigation/service"
	"sitenavigation/utils"

	"github.com/gin-gonic/gin"
)

var authService *service.UserService

func getAuthService() *service.UserService {
	if authService == nil {
		authService = service.NewUserService(config.DB)
	}
	return authService
}

// Login 用户登录
func Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	user, err := getAuthService().GetUserByAccount(req.Account)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "账号或密码错误"})
		return
	}

	if !getAuthService().VerifyPassword(user, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "账号或密码错误"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{
		Token: token,
		User:  user,
	})
}
