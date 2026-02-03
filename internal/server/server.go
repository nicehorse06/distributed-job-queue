package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewRouter configures the HTTP routes for the API.
func NewRouter() http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "distributed-job-queue api is running"})
	})

	return r
}
