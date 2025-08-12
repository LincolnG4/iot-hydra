package main

import (
	"net/http"

	"github.com/LincolnG4/iot-hydra/cmd/api"
	"github.com/gin-gonic/gin"
)

type application struct {
	PodmanRuntime *api.PodmanHandler
}

func (a *application) mount() http.Handler {
	router := gin.Default()

	gin.DisableConsoleColor()

	{
		v1 := router.Group("/v1")
		{
			event := v1.Group("/containers")
			event.GET("/", a.PodmanRuntime.ListAll)
		}

	}

	return router
}
