package main

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/LincolnG4/iot-hydra/internal/agent"
	"github.com/LincolnG4/iot-hydra/internal/runtimer"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type application struct {
	PodmanRuntime runtimer.PodmanRuntime
	logger        *zerolog.Logger
	config        *config

	// Telemetry Agent
	ctx            context.Context
	cancel         *context.CancelFunc
	TelemetryAgent *agent.TelemetryAgent
}

type config struct {
	Addr string
}

func (a *application) mount() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	gin.DisableConsoleColor()
	// Add the otelgin middleware
	router.Use(otelgin.Middleware("iot-hydra-runtime"))
	{
		v1 := router.Group("/v1")
		{
			containers := v1.Group("/containers")
			containers.POST("/", a.createContainer)
			containers.GET("/", a.listContainer)

			containers.GET("/:name", a.checkContainer)
			containers.POST("/:name/start", a.startContainer)
			containers.POST("/:name/stop", a.stopContainer)
			containers.DELETE("/:name", a.deleteContainer)

			// health endpoint
			health := v1.Group("/health")
			health.GET("/", a.healthChecker)

			// websocket message driven
			iotAgent := v1.Group("/ws")
			iotAgent.GET("/", a.websocketIoTHandler)
		}

	}

	return router
}

func (a *application) startTelemetryAgent() {
	a.TelemetryAgent.Start()
	go func() {
		for {
			select {
			case msg := <-a.TelemetryAgent.Queue:
				if err := a.TelemetryAgent.RouteToBrokers(msg); err != nil {
					a.logger.Error().Err(err).Msg("")
				}
			case <-a.ctx.Done():
				a.logger.Info().Msg("telemtry agent stopped")
			}
		}
	}()
}

func (a *application) run(r *gin.Engine) error {
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	a.ctx = context.Background()
	a.logger.Info().Msg("start message reader")

	a.logger.Info().Msg("server has started at localhost:8080")
	return srv.ListenAndServe()
}
