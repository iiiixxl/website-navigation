package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn := GetDatabaseDSN()
	if dsn == "" {
		return fmt.Errorf("数据库连接字符串未配置")
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("打开数据库连接失败: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	sqlDB.SetMaxOpenConns(GetDatabaseMaxOpenConns())
	sqlDB.SetMaxIdleConns(GetDatabaseMaxIdleConns())

	log.Println("数据库连接成功")
	return nil
}

func CloseDB() {
	if DB == nil {
		return
	}
	sqlDB, err := DB.DB()
	if err == nil {
		_ = sqlDB.Close()
		log.Println("数据库连接已关闭")
	}
}
