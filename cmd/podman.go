package main

import (
	"net/http"

	"github.com/LincolnG4/iot-hydra/internal/runtimer"
	"github.com/gin-gonic/gin"
)

type newContainerPayload struct {
	Name  string `json:"name" uri:"name"`
	Image string `json:"image"`
}

// createContaiener deploys new podman container from image.
func (a *application) createContainer(c *gin.Context) {
	var newContainer newContainerPayload

	if err := c.BindJSON(&newContainer); err != nil {
		a.logger.Error().Err(err).Msg("failed to bind JSON payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload", "details": err.Error()})
		return
	}

	container := runtimer.Container{
		Name:  newContainer.Name,
		Image: newContainer.Image,
	}

	if err := a.PodmanRuntime.CreateContainer(container); err != nil {
		a.logger.Error().Err(err).Str("container_name", container.Name).Str("image", container.Image).Msg("failed to create container")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create container", "details": err.Error()})
		return
	}

	a.logger.Info().Str("container_name", container.Name).Str("image", container.Image).Msg("container created successfully")
	c.JSON(http.StatusCreated, gin.H{"status": "container created successfully", "container": container})
}

// checkContainer returns the current status of the container requested.
func (a *application) checkContainer(c *gin.Context) {
	var newContainer newContainerPayload
	if err := c.ShouldBindUri(&newContainer); err != nil {
		a.logger.Error().Err(err).Msg("failed to bind URI parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid container name", "details": err.Error()})
		return
	}

	container := runtimer.Container{
		Name: newContainer.Name,
	}

	conInfo, err := a.PodmanRuntime.CheckContainer(container.Name)
	if err != nil {
		a.logger.Error().Err(err).Str("container_name", container.Name).Msg("failed to check container")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check container", "details": err.Error()})
		return
	}

	a.logger.Debug().Str("container_name", container.Name).Msg("container checked successfully")
	c.JSON(http.StatusOK, gin.H{"status": "container checked successfully", "container": conInfo})
}

// startContainer starts the container based on it name.
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

// stopContainer stops the container based on it name.
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

// deleteContainer remove the container.
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

// listContainer lists all containers on the device.
func (a *application) listContainer(c *gin.Context) {
	containers, err := a.PodmanRuntime.ListContainers()
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to list containers")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list containers", "details": err.Error()})
		return
	}

	a.logger.Debug().Int("container_count", len(containers)).Msg("containers listed successfully")
	c.JSON(http.StatusOK, gin.H{"status": "containers listed successfully", "containers": containers, "count": len(containers)})
}
