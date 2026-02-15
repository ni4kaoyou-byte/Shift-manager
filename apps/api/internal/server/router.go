package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/middleware"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(
		middleware.RequestID(),
		middleware.Logging(),
		middleware.Recover(),
	)
	routeRegistrars := newRouteRegistrars()

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	apiV1 := router.Group("/api/v1")
	apiV1.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "api v1"})
	})
	for _, routeRegistrar := range routeRegistrars {
		routeRegistrar.RegisterRoutes(apiV1)
	}

	return router
}
