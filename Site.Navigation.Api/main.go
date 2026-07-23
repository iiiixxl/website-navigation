package main

import (
	"log"
	"sitenavigation/config"
	"sitenavigation/router"
	"sitenavigation/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	jwtSecret := config.GetJWTSecret()
	if jwtSecret == "" {
		log.Fatal("JWT secret 未配置，请在 config/config.json 中设置 jwt.secret")
	}
	utils.InitJWT(jwtSecret, config.GetJWTExpireHours())

	if err := config.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer config.CloseDB()

	if config.IsRedisEnabled() {
		if err := config.InitRedis(); err != nil {
			log.Fatalf("Redis 初始化失败: %v", err)
		}
		defer config.CloseRedis()
	} else {
		log.Println("Redis 未启用（redis.enabled=false）")
	}

	r := gin.Default()
	router.RegisterRoutes(r)

	port := config.GetServerPort()
	log.Printf("服务启动: http://localhost%s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
