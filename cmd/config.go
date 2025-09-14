package main

import "github.com/LincolnG4/iot-hydra/internal/agent"

var AppConfig struct {
	Server         string               `yaml:"server"`
	TelemetryAgent agent.TelemetryAgent `yaml:"telemetryAgent"`
}
