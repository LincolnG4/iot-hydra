package main

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// mount setups routes for mutex
func (a *application) routes() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	gin.DisableConsoleColor()

	// Add the otelgin middleware
	router.Use(otelgin.Middleware("iot-hydra-runtime-app"))
	{
		v1 := router.Group("/v1")
		{
			// Podman Route
			containers := v1.Group("/containers")
			containers.POST("/", a.createContainer)           // Create Container
			containers.GET("/", a.listContainer)              // List all containers
			containers.GET("/:name", a.checkContainer)        // Check status container
			containers.POST("/:name/start", a.startContainer) // Start container
			containers.POST("/:name/stop", a.stopContainer)   // Stop container
			containers.DELETE("/:name", a.deleteContainer)    // Delete container

			// health endpoint
			health := v1.Group("/health")
			health.GET("/", a.healthChecker) // Check Health

			// websocket message driven
			iotAgent := v1.Group("/ws")
			iotAgent.GET("", a.websocketIoTHandler) // Websocket message driven
		}

	}
	return router
}
