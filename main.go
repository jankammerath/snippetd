package main

import (
	"context"
	"log"
	"net/http"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces" // Import the namespaces package
	"github.com/gin-gonic/gin"
)

const containerdSocket = "/run/containerd/containerd.sock"
const snippetNamespace = "snippetd" // Define a default namespace

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	router.GET("/containers", func(c *gin.Context) {
		client, err := containerd.New(containerdSocket)
		if err != nil {
			log.Printf("Error creating containerd client: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to containerd"})
			return
		}
		defer client.Close()

		// Use a context with a namespace
		ctx := namespaces.WithNamespace(context.Background(), snippetNamespace)

		containers, err := client.Containers(ctx)
		if err != nil {
			log.Printf("Error listing containers: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list containers"})
			return
		}

		containerNames := make([]string, len(containers))
		for i, container := range containers {
			containerNames[i] = container.ID()
		}

		c.JSON(http.StatusOK, gin.H{"containers": containerNames})
	})

	router.Run(":8080")
}
