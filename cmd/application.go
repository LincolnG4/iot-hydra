package main

import (
	"context"
	"fmt"
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

func (a *application) startTelemetryAgent(ctx context.Context) error {
	var err error
	a.TelemetryAgent, err = agent.NewTelemetryAgent(&a.config.TelemetryAgent)
	if err != nil {
		return err
	}

	// start workerpool
	a.TelemetryAgent.WorkerPool.Start()

	a.ctx = ctx
	go func() {
		for {
			select {
			case msg := <-a.TelemetryAgent.Queue:
				for _, brokerName := range msg.TargetBrokers {
					b, exist := a.TelemetryAgent.Brokers[brokerName]
					if !exist {
						a.logger.Error().Msg(fmt.Sprintf("broker '%s' is not configured", brokerName))
						continue
					}

					a.TelemetryAgent.WorkerPool.JobQueue <- func() error {
						if err := b.Publish(ctx, msg); err != nil {
							a.TelemetryAgent.WorkerPool.ResultQueue <- fmt.Errorf("failed to publish message to broker '%s': %w", brokerName, err)
						}
						return nil
					}
				}
			case <-ctx.Done():
				a.logger.Info().Msg("telemetry agent stopped")
				return
			}
		}
	}()
	return nil
}

func (a *application) run(ctx context.Context, r *gin.Engine) error {
	err := a.startTelemetryAgent(ctx)
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

	a.ctx = ctx
	a.logger.Info().Str("address", a.config.APIService.Address).Msg("starting server")

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for shutdown signal
	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		a.logger.Info().Msg("shutdown signal received, stopping server...")

		// Graceful shutdown with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			a.logger.Error().Err(err).Msg("server forced to shutdown")
			return err
		}

		a.logger.Info().Msg("server stopped gracefully")
		return nil
	}
}
