package handlers

import (
	"blog-backend/config"
	"blog-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	postIDStr := c.Param("id")

	postID64, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.UserID = userID
	comment.PostID = uint(postID64)

	config.DB.Create(&comment)
	c.JSON(http.StatusCreated, comment)
}

func GetComments(c *gin.Context) {
	postID := c.Param("id")
	var comments []models.Comment

	config.DB.Where("post_id = ?", postID).Preload("User").Find(&comments)
	c.JSON(200, comments)
}
