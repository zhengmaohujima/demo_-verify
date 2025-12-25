package handlers

import (
	"net/http"

	"blog-backend/config"
	"blog-backend/models"
	"blog-backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "register success"})
}
func Login(c *gin.Context) {
	var req models.User
	var user models.User

	c.ShouldBindJSON(&req)

	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.Username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
