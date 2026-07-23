package middleware

import (
	"net/http"
	"strings"
	"sitenavigation/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
// 白名单路由或 AllowAnonymous 标记的接口会跳过鉴权
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if allowAnonymous, exists := c.Get(AllowAnonymousKey); exists {
			if allow, ok := allowAnonymous.(bool); ok && allow {
				c.Next()
				return
			}
		}

		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		if IsAnonymousRoute(method, path) {
			c.Next()
			return
		}
		if path != c.Request.URL.Path && IsAnonymousRoute(method, c.Request.URL.Path) {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少Authorization头"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization格式错误，应为: Bearer <token>"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token: " + err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("account", claims.Account)
		c.Next()
	}
}
