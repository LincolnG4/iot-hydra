package runtimer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateContainer_Success(t *testing.T) {
	// Arrange
	mockManager := &MockPodmanManager{
		CreateContainerFunc: func(container Container) error {
			// Simulate a successful container creation
			return nil
		},
	}

	runtimer := NewRuntimer()
	runtimer.PodmanManager = mockManager

	container := Container{Name: "test-container", Image: "test-image"}

	// Act
	err := runtimer.PodmanManager.CreateContainer(container)

	// Assert
	assert.NoError(t, err)
}

func TestCreateContainer_Failure(t *testing.T) {
	// Arrange
	expectedErr := errors.New("failed to create container")
	mockManager := &MockPodmanManager{
		CreateContainerFunc: func(container Container) error {
			// Simulate a failure during container creation
			return expectedErr
		},
	}

	runtimer := NewRuntimer()
	runtimer.PodmanManager = mockManager

	container := Container{Name: "test-container", Image: "test-image"}

	// Act
	err := runtimer.PodmanManager.CreateContainer(container)

	// Assert
	assert.Equal(t, expectedErr, err)
}

func TestCheckContainer_Success(t *testing.T) {
	// Arrange
	expectedContainer := Container{
		Name:  "my-container",
		Image: "my-image",
		State: "running",
	}
	mockManager := &MockPodmanManager{
		CheckContainerFunc: func(name string) (Container, error) {
			assert.Equal(t, "my-container", name)
			return expectedContainer, nil
		},
	}

	runtimer := NewRuntimer()
	runtimer.PodmanManager = mockManager

	// Act
	result, err := runtimer.PodmanManager.CheckContainer("my-container")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedContainer, result)
}

func TestListContainers_Success(t *testing.T) {
	// Arrange
	expectedContainers := []Container{
		{Name: "container-1", Image: "image-1", State: "running"},
		{Name: "container-2", Image: "image-2", State: "exited"},
	}
	mockManager := &MockPodmanManager{
		ListContainersFunc: func() ([]Container, error) {
			return expectedContainers, nil
		},
	}

	runtimer := NewRuntimer()
	runtimer.PodmanManager = mockManager

	// Act
	result, err := runtimer.PodmanManager.ListContainers()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedContainers, result)
}

func TestStartContainer_Success(t *testing.T) {
	// Arrange
	mockManager := &MockPodmanManager{
		StartContainerFunc: func(name string) error {
			assert.Equal(t, "test-container", name)
			return nil
		},
	}

	runtimer := NewRuntimer()
	runtimer.PodmanManager = mockManager

	// Act
	err := runtimer.PodmanManager.StartContainer("test-container")

	// Assert
	assert.NoError(t, err)
}
