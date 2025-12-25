package test

import (
	"blog-backend/handlers"
	"blog-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
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

	return r
}
