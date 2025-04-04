package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	router.POST("/execute", func(c *gin.Context) {
		// Here you would handle the code execution request

		// For now, just return a placeholder response
		c.JSON(http.StatusOK, gin.H{
			"message": "Code execution request received",
		})
	})

	router.Run(":8080")
}
