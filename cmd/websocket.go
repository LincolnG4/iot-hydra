package main

import (
	"fmt"
	"time"

	"github.com/LincolnG4/iot-hydra/internal/brokers"
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
		msg := message. 
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
