package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GivePing(c *gin.Context) {
	c.JSON(http.StatusFound, gin.H{"message": "welcome to rest api"})
}
