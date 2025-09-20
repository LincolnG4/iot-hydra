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
	defer conn.Close()

	a.logger.Info().Msg("Client connected to Web Socket")
	for {
		msg := message.Message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// This is a real error (e.g., network failure).
				a.logger.Error().Err(err).Msg("Unexpected WebSocket close error")
			} else {
				// This is a normal disconnect from the client.
				a.logger.Info().Msg("WebSocket client disconnected normally.")
			}
			break
		}
		msg.ID = fmt.Sprintf("ws-%d", time.Now().UnixNano())
		msg.Timestamp = time.Now()

		select {
		case a.TelemetryAgent.Queue <- &msg:
			a.logger.Debug().Msg("Message received via WebSocket: " + msg.ID)
		default:
			a.logger.Warn().Msg("Telemetry queue is full. Dropping message from client.")
		}

		a.logger.Debug().Msg("Message received via WebSocket:" + msg.ID)
	}
}
