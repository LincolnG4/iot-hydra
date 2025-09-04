package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *application) healthChecker(c *gin.Context) {
	a.logger.Debug().Msg("Status OK")
	c.JSON(http.StatusOK, gin.H{"status": "Ok"})
}
