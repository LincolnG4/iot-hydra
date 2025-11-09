package main

import (
	"context"
	"net/http"
	"time"

	"github.com/LincolnG4/iot-hydra/internal/agent"
	"github.com/LincolnG4/iot-hydra/internal/config"
	"github.com/LincolnG4/iot-hydra/internal/runtimer"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type application struct {
	PodmanRuntime  runtimer.PodmanRuntime // responsible to manage podman service
	config         *config.ConfigYAML     // configuration loaded from config.yaml file
	TelemetryAgent *agent.TelemetryAgent  // agent responsible to route telemetry messsages to the brokers

	logger *zerolog.Logger
	ctx    context.Context
	cancel *context.CancelFunc
}

func (a *application) run(ctx context.Context, r *gin.Engine) error {
	// Start agent to collect data and send to brokers
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

		return nil
	}
}
