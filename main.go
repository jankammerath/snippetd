package main

import (
	"encoding/json"
	"net/http"

	"embed"

	"github.com/gin-gonic/gin"
	"github.com/jankammerath/snippetd/snippetd"
)

//go:embed config
var configDir embed.FS

func main() {
	router := gin.Default()

	// Load the configuration files from the embedded filesystem
	runtimesData, err := configDir.ReadFile("config/runtimes.json")
	if err != nil {
		panic(err)
	}

	var runtimes map[string]snippetd.RuntimeConfig
	err = json.Unmarshal(runtimesData, &runtimes)
	if err != nil {
		panic(err)
	}

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	router.GET("/runtimes", func(c *gin.Context) {
		// return the key of the map and the mime type array
		config := make(map[string][]string)
		for key, value := range runtimes {
			config[key] = value.MimeTypes
		}

		c.JSON(http.StatusOK, config)
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
