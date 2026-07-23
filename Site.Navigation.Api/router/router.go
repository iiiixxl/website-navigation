package router

import (
	"sitenavigation/controller"
	"sitenavigation/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/health", controller.Health)

	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware())
	{
		auth := v1.Group("/auth")
		{
			RegisterRoute(auth, RouteConfig{
				Method:         "POST",
				Path:           "/login",
				Handler:        controller.Login,
				AllowAnonymous: true,
			})
		}

		users := v1.Group("/users")
		{
			RegisterRoute(users, RouteConfig{
				Method:         "POST",
				Path:           "/CreateUser",
				Handler:        controller.CreateUser,
				AllowAnonymous: true,
			})
			users.GET("/GetUserList", controller.GetAllUsers)
			users.GET("/GetUserById", controller.GetUser)
			users.POST("/UpdateUser", controller.UpdateUser)
			users.POST("/DeleteUser", controller.DeleteUser)
		}

		// AI 提问模板
		prompts := v1.Group("/prompts")
		{
			RegisterRoute(prompts, RouteConfig{
				Method:         "GET",
				Path:           "/GetTree",
				Handler:        controller.GetPromptTree,
				AllowAnonymous: true,
			})
			RegisterRoute(prompts, RouteConfig{
				Method:         "GET",
				Path:           "/GetItemById",
				Handler:        controller.GetPromptItem,
				AllowAnonymous: true,
			})

			prompts.POST("/CreateCategory", controller.CreatePromptCategory)
			prompts.POST("/UpdateCategory", controller.UpdatePromptCategory)
			prompts.POST("/DeleteCategory", controller.DeletePromptCategory)
			prompts.POST("/CreateItem", controller.CreatePromptItem)
			prompts.POST("/UpdateItem", controller.UpdatePromptItem)
			prompts.POST("/DeleteItem", controller.DeletePromptItem)
		}
	}
}
