package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type AppConfig struct {
	Database DatabaseConfig `json:"database"`
	JWT      JWTConfig      `json:"jwt"`
	Redis    RedisConfig    `json:"redis"`
	Server   ServerConfig   `json:"server"`
}

type DatabaseConfig struct {
	DSN          string `json:"dsn"`
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
}

type JWTConfig struct {
	Secret      string `json:"secret"`
	ExpireHours int    `json:"expire_hours"`
}

type RedisConfig struct {
	Enabled      bool   `json:"enabled"`
	Address      string `json:"address"`
	Password     string `json:"password"`
	DB           int    `json:"db"`
	PoolSize     int    `json:"pool_size"`
	MinIdleConns int    `json:"min_idle_conns"`
}

type ServerConfig struct {
	Port string `json:"port"`
}

var cfg *AppConfig

// InitConfig 只从本地 config/config.json 读取配置
func InitConfig() error {
	candidates := []string{
		"config/config.json",
		filepath.Join("config", "config.json"),
	}

	var lastErr error
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			lastErr = err
			continue
		}

		var loaded AppConfig
		if err := json.Unmarshal(data, &loaded); err != nil {
			return fmt.Errorf("解析配置失败(%s): %w", path, err)
		}
		cfg = &loaded
		log.Printf("已加载配置: %s", path)
		return nil
	}

	return fmt.Errorf("未找到配置文件 config/config.json: %v", lastErr)
}

func getConfig() *AppConfig {
	if cfg == nil {
		log.Fatal("配置未初始化，请先调用 InitConfig")
	}
	return cfg
}

func GetDatabaseDSN() string {
	return getConfig().Database.DSN
}

func GetDatabaseMaxOpenConns() int {
	n := getConfig().Database.MaxOpenConns
	if n <= 0 {
		return 25
	}
	return n
}

func GetDatabaseMaxIdleConns() int {
	n := getConfig().Database.MaxIdleConns
	if n <= 0 {
		return 5
	}
	return n
}

func GetJWTSecret() string {
	return getConfig().JWT.Secret
}

func GetJWTExpireHours() int {
	n := getConfig().JWT.ExpireHours
	if n <= 0 {
		return 24
	}
	return n
}

func GetServerPort() string {
	port := getConfig().Server.Port
	if port == "" {
		port = "18080"
	}
	if !strings.HasPrefix(port, ":") {
		return ":" + port
	}
	return port
}

func IsRedisEnabled() bool {
	return getConfig().Redis.Enabled
}

func GetRedisAddress() string {
	addr := getConfig().Redis.Address
	if addr == "" {
		return "127.0.0.1:6379"
	}
	return addr
}

func GetRedisPassword() string {
	return getConfig().Redis.Password
}

func GetRedisDB() int {
	return getConfig().Redis.DB
}

func GetRedisPoolSize() int {
	n := getConfig().Redis.PoolSize
	if n <= 0 {
		return 10
	}
	return n
}

func GetRedisMinIdleConns() int {
	n := getConfig().Redis.MinIdleConns
	if n <= 0 {
		return 5
	}
	return n
}
