package runtimer

import (
	"github.com/stretchr/testify/mock"
)

type MockPodmanManager struct {
	mock.Mock
	// CreateContainerFunc func(container Container) error
	// CheckContainerFunc  func(name string) (Container, error)
	// StartContainerFunc  func(name string) error
	// StopContainerFunc   func(name string) error
	// DeleteContainerFunc func(name string) error
	// ListContainersFunc  func() ([]Container, error)
}

func (_m *MockPodmanManager) CreateContainer(container Container) error {
	ret := _m.Called(container)
	return ret.Error(0)
}

// CheckContainer provides a mock function with given fields: name
func (_m *MockPodmanManager) CheckContainer(name string) (Container, error) {
	ret := _m.Called(name)
	return ret.Get(0).(Container), ret.Error(1)
}

func (_m *MockPodmanManager) StartContainer(name string) error {
	ret := _m.Called(name)
	return ret.Error(0)
}

func (_m *MockPodmanManager) StopContainer(name string) error {
	ret := _m.Called(name)
	return ret.Error(0)
}

func (_m *MockPodmanManager) DeleteContainer(name string) error {
	ret := _m.Called(name)
	return ret.Error(0)
}

func (_m *MockPodmanManager) ListContainers() ([]Container, error) {
	ret := _m.Called()
	return ret.Get(0).([]Container), ret.Error(1)
}
