package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", home)
	router.GET("/image/*path", image)
	router.GET("/health", healthready)
	router.GET("/ready", healthready)

	router.Run(":8010")
}

func home(c *gin.Context) {
	body := map[string]interface{}{
		"route":   "home",
		"message": "hello world",
	}

	c.IndentedJSON(http.StatusOK, body)
}

func image(c *gin.Context) {
	body := map[string]interface{}{
		"route": "image",
		"path":  c.Param("path"),
	}

	c.IndentedJSON(http.StatusOK, body)
}

func healthready(c *gin.Context) {
	c.Status(http.StatusOK)
}
