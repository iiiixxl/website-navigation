package router

import (
	"sitenavigation/middleware"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	Method         string
	Path           string
	Handler        gin.HandlerFunc
	AllowAnonymous bool
}

// RegisterRoute 注册路由，支持 AllowAnonymous（类似 .NET [AllowAnonymous]）
func RegisterRoute(group *gin.RouterGroup, config RouteConfig) {
	fullPath := group.BasePath() + config.Path
	if fullPath == "" {
		fullPath = "/"
	}

	if config.AllowAnonymous {
		middleware.RegisterAnonymousRoute(config.Method, fullPath)
	}

	switch config.Method {
	case "GET":
		group.GET(config.Path, config.Handler)
	case "POST":
		group.POST(config.Path, config.Handler)
	case "PUT":
		group.PUT(config.Path, config.Handler)
	case "DELETE":
		group.DELETE(config.Path, config.Handler)
	case "PATCH":
		group.PATCH(config.Path, config.Handler)
	default:
		group.Any(config.Path, config.Handler)
	}
}
