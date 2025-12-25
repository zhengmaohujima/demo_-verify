package test

import (
	"log"

	"blog-backend/config"
	"blog-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupTestDB() {
	// MySQL 测试数据库 DSN
	dsn := "demo:cccc@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to MySQL test database:", err)
	}

	// 自动迁移表结构
	if err := config.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	log.Println("MySQL test database connected and tables migrated successfully")
}
