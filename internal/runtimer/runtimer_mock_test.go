package runtimer

// MockPodmanManager is a mock implementation of the PodmanManager interface for testing.
type MockPodmanManager struct {
	CreateContainerFunc func(container Container) error
	CheckContainerFunc  func(name string) (Container, error)
	StartContainerFunc  func(name string) error
	StopContainerFunc   func(name string) error
	DeleteContainerFunc func(name string) error
	ListContainersFunc  func() ([]Container, error)
}

func (m *MockPodmanManager) CreateContainer(container Container) error {
	if m.CreateContainerFunc != nil {
		return m.CreateContainerFunc(container)
	}
	return nil
}

func (m *MockPodmanManager) CheckContainer(name string) (Container, error) {
	if m.CheckContainerFunc != nil {
		return m.CheckContainerFunc(name)
	}
	return Container{}, nil
}

func (m *MockPodmanManager) StartContainer(name string) error {
	if m.StartContainerFunc != nil {
		return m.StartContainerFunc(name)
	}
	return nil
}

func (m *MockPodmanManager) StopContainer(name string) error {
	if m.StopContainerFunc != nil {
		return m.StopContainerFunc(name)
	}
	return nil
}

func (m *MockPodmanManager) DeleteContainer(name string) error {
	if m.DeleteContainerFunc != nil {
		return m.DeleteContainerFunc(name)
	}
	return nil
}

func (m *MockPodmanManager) ListContainers() ([]Container, error) {
	if m.ListContainersFunc != nil {
		return m.ListContainersFunc()
	}
	return nil, nil
}
