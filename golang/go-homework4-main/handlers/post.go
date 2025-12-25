package handlers

import (
	"blog-backend/config"
	"blog-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var post models.Post
	c.ShouldBindJSON(&post)

	post.UserID = userID
	config.DB.Create(&post)

	c.JSON(http.StatusCreated, post)
}
func GetPosts(c *gin.Context) {
	var posts []models.Post
	config.DB.Preload("User").Find(&posts)
	c.JSON(http.StatusOK, posts)
}
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(uint)

	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "post not found"})
		return
	}

	if post.UserID != userID {
		c.JSON(403, gin.H{"error": "no permission"})
		return
	}

	c.ShouldBindJSON(&post)
	config.DB.Save(&post)
	c.JSON(200, post)
}
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(uint)

	var post models.Post
	config.DB.First(&post, id)

	if post.UserID != userID {
		c.JSON(403, gin.H{"error": "no permission"})
		return
	}

	config.DB.Delete(&post)
	c.JSON(200, gin.H{"message": "deleted"})
}
