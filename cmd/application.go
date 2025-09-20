package main

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/LincolnG4/iot-hydra/internal/agent"
	"github.com/LincolnG4/iot-hydra/internal/config"
	"github.com/LincolnG4/iot-hydra/internal/runtimer"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type application struct {
	// responsible to manage podman service
	PodmanRuntime runtimer.PodmanRuntime
	logger        *zerolog.Logger

	// configuration loaded from config.yaml file
	config *config.ConfigYAML

	ctx    context.Context
	cancel *context.CancelFunc

	// agent responsible to route telemetry messsages to the brokers
	TelemetryAgent *agent.TelemetryAgent
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
			iotAgent.GET("", a.websocketIoTHandler)
		}

	}
	return router
}

func (a *application) startTelemetryAgent() error {
	var err error
	a.TelemetryAgent, err = agent.NewTelemetryAgent(&a.config.TelemetryAgent)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case msg := <-a.TelemetryAgent.Queue:
				if err := a.TelemetryAgent.RouteToBrokers(msg); err != nil {
					a.logger.Error().Err(err).Msg("")
				}
			case <-a.ctx.Done():
				a.logger.Info().Msg("telemetry agent stopped")
			}
		}
	}()
	return nil
}

func (a *application) run(r *gin.Engine) error {
	err := a.startTelemetryAgent()
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:         a.config.APIService.Address,
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
