package main

import (
	"blog-backend/config"
	"blog-backend/handlers"
	"blog-backend/middleware"
	"blog-backend/models"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	config.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	r := gin.Default()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/posts", handlers.CreatePost)
		auth.PUT("/posts/:id", handlers.UpdatePost)
		auth.DELETE("/posts/:id", handlers.DeletePost)
		auth.POST("/posts/:id/comments", handlers.CreateComment)
	}

	r.GET("/posts", handlers.GetPosts)
	r.GET("/posts/:id/comments", handlers.GetComments)

	r.Run(":8080")
}
