package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type computeService interface {
	Square(ctx context.Context, value int64) (int64, error)
}

// NewRouter configures the HTTP routes for the API.
func NewRouter(computeSvc computeService) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "distributed-job-queue api is running"})
	})

	type squareRequest struct {
		Value *int64 `json:"value"`
	}

	r.POST("/compute/square", func(c *gin.Context) {
		if computeSvc == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "compute service not configured"})
			return
		}

		var req squareRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Value == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "value is required"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		square, err := computeSvc.Square(ctx, *req.Value)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"value":  *req.Value,
			"square": square,
			"engine": "rust-grpc",
		})
	})

	return r
}
