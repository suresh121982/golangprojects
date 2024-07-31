package main

import (
	"fmt"

	"example.com/myapp/intializers"
	"github.com/gin-gonic/gin"
)

func init() {
	intializers.LoadEnvVariables()

}

func main() {
	fmt.Println("Hello Suresh")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong pong",
		})
	})
	r.Run() //
}
