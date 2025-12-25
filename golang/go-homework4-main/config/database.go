package config

import (
	"blog-backend/models"
	"log"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	//DB, err = gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	//if err != nil {
	//	log.Fatal("failed to connect database")
	//}
	dsn := "demo:cccc@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	DB, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err := DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	log.Println("Database connected and tables migrated successfully")

}
