package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/LincolnG4/iot-hydra/internal/config"
	"github.com/LincolnG4/iot-hydra/internal/runtimer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	log.Info().Msg("starting IoT runtime")

	// Handle SIGINT (CTRL+C) and SIGTERM gracefully
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// load configuration file
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config.yaml" // default config file
		log.Info().Str("config_file", configFile).Msg("using default config file")
	}

	// create service configuration structure
	log.Debug().Str("config_file", configFile).Msg("loading configuration")
	cfg, err := config.NewConfigFromYAML(configFile)
	if err != nil {
		log.Error().Err(err).Str("config_file", configFile).Msg("failed to load configuration")
		os.Exit(1)
	}

	// OpenTelemetry setup
	log.Debug().Msg("starting OpenTelemetry")

	// Set up OpenTelemetry.
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
		log.Error().Err(err).Msg("")
	}()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample()))))

	log.Debug().Msg("starting PodmanRuntime")
	socketPath := os.Getenv("PODMAN_SOCKET_PATH")
	if socketPath == "" {
		socketPath = "unix:///run/podman/podman.sock" // default podman socket
		log.Info().Str("socket_path", socketPath).Msg("using default Podman socket")
	}

	podmanRuntime, err := runtimer.NewPodmanManager(
		&runtimer.ManagerOptions{
			SocketPath: socketPath,
		},
	)
	if err != nil {
		log.Error().Err(err).Str("socket_path", socketPath).Msg("failed to initialize Podman runtime")
		os.Exit(1)
	}

	app := application{
		PodmanRuntime: &podmanRuntime,
		logger:        &logger,
		config:        &cfg,
	}

	mux := app.mount()

	if err := app.run(ctx, mux); err != nil {
		log.Error().Err(err).Msg("application failed")
		os.Exit(1)
	}
}
