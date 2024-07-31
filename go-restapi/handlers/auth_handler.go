// handlers/auth_handler.go
package handlers

import (
	"net/http"

	"examples.com/go-restapi/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", Login)
	}
}

func Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Replace with actual authentication logic (e.g., database check)
	if loginData.Username == "admin" && loginData.Password == "password" {
		token, err := middleware.GenerateToken(loginData.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}
