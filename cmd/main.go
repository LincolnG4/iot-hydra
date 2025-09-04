package main

import (
	"os"

	"github.com/LincolnG4/iot-hydra/internal/runtimer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)

	log.Info().Msg("starting IoT runtime")
	// OpenTelemetry setup
	log.Debug().Msg("starting OpenTelemetry")
	// Handle SIGINT (CTRL+C) gracefully.
	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer stop()

	// // Set up OpenTelemetry.
	// otelShutdown, err := setupOTelSDK(ctx)
	// if err != nil {
	// 	log.Error().Err(err).Msg("")
	// 	return
	// }
	// // Handle shutdown properly so nothing leaks.
	// defer func() {
	// 	err = errors.Join(err, otelShutdown(context.Background()))
	// 	log.Error().Err(err).Msg("")
	// }()
	// otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample()))))

	sock := os.Getenv("podman_socket_path")
	log.Info().Msg(sock)
	log.Debug().Msg("starting PodmanRuntime")
	podmanRuntime, err := runtimer.NewPodmanManager(
		&runtimer.ManagerOptions{
			SocketPath: os.Getenv("podman_socket_path"),
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	app := application{
		PodmanRuntime: &podmanRuntime,
		logger:        &logger,
		config: &config{
			Addr: ":8080",
		},
	}

	mux := app.mount()

	app.run(mux)
}
