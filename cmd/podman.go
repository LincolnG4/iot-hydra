package main

import (
	"net/http"

	"github.com/LincolnG4/iot-hydra/internal/runtimer"
	"github.com/gin-gonic/gin"
)

type newContainerPayload struct {
	Name  string `json:"name" uri:"name" inding:"required,uuid"`
	Image string `json:"image"`
}

func (a *application) createContainer(c *gin.Context) {
	var newContainer newContainerPayload

	if err := c.BindJSON(&newContainer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	container := runtimer.Container{
		Name:  newContainer.Name,
		Image: newContainer.Image,
	}

	if err := a.PodmanRuntime.CreateContainer(container); err != nil {
		a.logger.Error().Err(err).Msg("")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "request created"})
}

func (a *application) checkContainer(c *gin.Context) {
	var newContainer newContainerPayload
	if err := c.ShouldBindUri(&newContainer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	container := runtimer.Container{
		Name: newContainer.Name,
	}

	conInfo, err := a.PodmanRuntime.CheckContainer(container.Name)
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"response": conInfo})
}

func (a *application) startContainer(c *gin.Context) {
	var newContainer newContainerPayload
	if err := c.ShouldBindUri(&newContainer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	container := runtimer.Container{
		Name: newContainer.Name,
	}

	err := a.PodmanRuntime.StartContainer(container.Name)
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"response": "container started"})
}

func (a *application) stopContainer(c *gin.Context) {
	var newContainer newContainerPayload
	if err := c.ShouldBindUri(&newContainer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	container := runtimer.Container{
		Name: newContainer.Name,
	}

	err := a.PodmanRuntime.StopContainer(container.Name)
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"response": "container stopped"})
}

func (a *application) deleteContainer(c *gin.Context) {
	var newContainer newContainerPayload
	if err := c.ShouldBindUri(&newContainer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	container := runtimer.Container{
		Name: newContainer.Name,
	}

	err := a.PodmanRuntime.DeleteContainer(container.Name)
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"response": "container deleted"})
}

func (a *application) listContainer(c *gin.Context) {
	containers, err := a.PodmanRuntime.ListContainers()
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"response": containers})
}
