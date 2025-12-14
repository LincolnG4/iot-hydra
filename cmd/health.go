package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (a *application) healthChecker(c *gin.Context) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "iot-hydra",
		"version":   "1.0.0", // TODO: Get from build info
	}

	// Add queue health if telemetry agent is available
	if a.TelemetryAgent != nil {
		health["telemetry"] = map[string]interface{}{
			"queue_length":      len(a.TelemetryAgent.Queue),
			"queue_capacity":    cap(a.TelemetryAgent.Queue),
			"brokers_connected": len(a.TelemetryAgent.Brokers),
		}
	}

	// Add podman runtime health
	if a.PodmanRuntime != nil {
		// Test podman connection by listing containers
		_, err := a.PodmanRuntime.ListContainers()
		if err != nil {
			health["podman"] = map[string]interface{}{
				"status": "unhealthy",
				"error":  err.Error(),
			}
			health["status"] = "degraded"
		} else {
			health["podman"] = map[string]interface{}{
				"status": "healthy",
			}
		}
	}

	statusCode := http.StatusOK
	if health["status"] == "degraded" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, health)
}
