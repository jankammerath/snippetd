package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"embed"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jankammerath/snippetd/snippetd"
)

//go:embed config
var configDir embed.FS

// runtimes is a map of runtime configurations
var runtimes map[string]snippetd.RuntimeConfig

func GetRuntimeConfig(mimeType string) (snippetd.RuntimeConfig, error) {
	// iterate over the runtimes map and find the runtime config for the given mime type
	for language, runtime := range runtimes {
		for _, mt := range runtime.MimeTypes {
			if mt == mimeType {
				result := runtime

				runScript, err := configDir.ReadFile("config/" + language + ".sh")
				if err != nil {
					panic(err)
				}

				result.RunScript = string(runScript)

				return result, nil
			}
		}
	}

	return snippetd.RuntimeConfig{}, errors.New("Mime type not supported: " + mimeType)
}

func main() {
	router := gin.Default()

	// Load the configuration files from the embedded filesystem
	runtimesData, err := configDir.ReadFile("config/runtimes.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(runtimesData, &runtimes)
	if err != nil {
		panic(err)
	}

	// instanciate the code runtime
	codeRuntime, err := snippetd.NewCodeRuntime()
	if err != nil {
		panic(err)
	}
	defer codeRuntime.Close()

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
		// determine the mime type of the posted file
		mimeType := c.GetHeader("Content-Type")
		if mimeType == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Mime type not specified",
			})
			return
		}

		// get the runtime config for the mime type
		runtimeConfig, err := GetRuntimeConfig(mimeType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// read the source code from the request body
		sourceCode, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read source code",
			})
			return
		}

		executionUuid := uuid.New().String()
		result := codeRuntime.Execute(executionUuid, string(sourceCode), runtimeConfig)

		// return the result object
		c.JSON(http.StatusOK, gin.H{
			"uuid":   executionUuid,
			"result": result,
		})
	})

	router.Run(":8080")
}
