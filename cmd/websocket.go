package main

import (
	"fmt"
	"time"

	"github.com/LincolnG4/iot-hydra/internal/message"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (a *application) websocketIoTHandler(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		a.logger.Error().Err(fmt.Errorf("failed to set websocket upgrade: %+v", err)).Msg("")
		return
	}

	a.logger.Info().Msg("Client connected to Web Socket")
	for {
		msg := message.Message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			a.logger.Error().Err(fmt.Errorf("WebSocket error: %v", err)).Msg("")
			break
		}

		msg.ID = fmt.Sprintf("ws-%d", time.Now().UnixNano())
		msg.Timestamp = time.Now()

		err = a.TelemetryAgent.RouteToBrokers(&msg)
		if err != nil {
			a.logger.Debug().Msg("Failed to send message" + msg.ID + err.Error())
		}

		a.logger.Debug().Msg("Message received via WebSocket:" + msg.ID)
	}
}
