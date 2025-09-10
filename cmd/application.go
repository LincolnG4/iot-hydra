package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/LincolnG4/iot-hydra/internal/runtimer"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type application struct {
	PodmanRuntime runtimer.PodmanRuntime
	IoTAgent      *IoTAgent
	logger        *zerolog.Logger
	config        *config
	MessageQueue  chan IoTMessage
	ctx           context.Context
	cancel        context.CancelFunc
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

func (a *application) startMessageAgent() {
	go func() {
		for {
			select {
			case msg := <-a.MessageQueue:
				if err := a.IoTAgent.Route(&msg); err != nil {
					a.logger.Error().Err(err).Msg("")
				}
			case <-a.ctx.Done():
				a.logger.Info().Msg("iot message agent stopped")
			}
		}
	}()
}

func (a *application) websocketIoTHandler(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		a.logger.Error().Err(fmt.Errorf("failed to set websocket upgrade: %+v", err)).Msg("")
		return
	}

	a.logger.Info().Msg("Client connected to Web Socket")
	for {
		msg := IoTMessage{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			a.logger.Error().Err(fmt.Errorf("WebSocket error: %v", err)).Msg("")
			break
		}

		msg.ID = fmt.Sprintf("ws-%d", time.Now().UnixNano())
		msg.Timestamp = time.Now()
		select {
		case a.MessageQueue <- msg:
			a.logger.Debug().Msg("Message received via WebSocket:" + msg.ID)
		default:
			a.logger.Debug().Msg("Message queue full, dropping message: %s" + msg.ID)
		}
	}
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
	a.startMessageAgent()

	a.logger.Info().Msg("server has started at localhost:8080")
	return srv.ListenAndServe()
}
