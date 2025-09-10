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
	// ctx, stop := signal.NotifyCo2ntext(context.Background(), os.Interrupt)
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

	iotagent := NewIoTAgent(&ServiceConfig{
		Nats: NATSConfig{
			URL: "localhost:4222",
			BasicAuth: &BasicAuth{
				Username: "test",
				Password: "test",
			},
		},
	})

	app := application{
		PodmanRuntime: &podmanRuntime,
		IoTAgent:      &iotagent,
		logger:        &logger,
		config: &config{
			Addr: ":8080",
		},
		// TODO add config number of messages in the queue
		MessageQueue: make(chan IoTMessage, 10000),
	}

	mux := app.mount()

	app.run(mux)
}
