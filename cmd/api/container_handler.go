package api

import (
	"github.com/LincolnG4/iot-hydra/internal/runtime"
	"github.com/gin-gonic/gin"
)

type PodmanHandler struct {
	runtime *runtime.Podman
}

func (a *PodmanHandler) ListAll(c *gin.Context) {}
